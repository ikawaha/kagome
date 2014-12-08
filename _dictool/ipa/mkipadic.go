//  Copyright (c) 2014 ikawaha.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"code.google.com/p/go.text/encoding/japanese"
	"code.google.com/p/go.text/transform"

	"github.com/ikawaha/kagome"
)

const (
	ipaMatrixDefFileName = "matrix.def"
	ipaCharDefFileName   = "char.def"
	ipaUnkDefFileName    = "unk.def"

	ipaDicArchiveFileName    = "ipa.dic"
	ipaDicMorphFileName      = "morph.dic"
	ipaDicIndexFileName      = "index.dic"
	ipaDicConnectionFileName = "connection.dic"
	ipaDicCharDefFileName    = "chardef.dic"
	ipaDicUnkFileName        = "unk.dic"

	ipaMorphCsvColSize                    = 13
	ipaMrophRecordSurfaceIndex            = 0
	ipaMorphRecordLeftIdIndex             = 1
	ipaMorphRecordRightIdIndex            = 2
	ipaMorphRecordWeightIndex             = 3
	ipaMorphRecordOtherContentsStartIndex = 4

	ipaUnkRecordSize                    = 11
	ipaUnkRecordCategoryIndex           = 0
	ipaUnkRecordLeftIdIndex             = 1
	ipaUnkRecordRightIndex              = 2
	ipaUnkRecordWeigthIndex             = 3
	ipaUnkRecordOtherContentsStartIndex = 4
)

type IpaDic struct {
	Morphs       []kagome.Morph
	Contents     [][]string
	Index        kagome.Trie
	IndexDup     map[int]int
	Connection   kagome.ConnectionTable
	CharClass    []string
	CharCategory []byte
	InvokeList   []bool
	GroupList    []bool

	UnkMorphs   []kagome.Morph
	UnkIndex    map[int]int
	UnkIndexDup map[int]int
	UnkContents [][]string
}

type ipaDicPath struct {
	Morph      string
	Index      string
	Connection string
	CharDef    string
	Unk        string
}

func loadIpaDic(path ipaDicPath) (dic *IpaDic, err error) {
	dic = new(IpaDic)
	if err = func() error {
		f, e := os.Open(path.Morph)
		if e != nil {
			return e
		}
		dec := gob.NewDecoder(f)
		if e = dec.Decode(&dic.Morphs); e != nil {
			return e
		}
		if e = dec.Decode(&dic.Contents); e != nil {
			return e
		}
		return nil
	}(); err != nil {
		return
	}
	if err = func() error {
		f, e := os.Open(path.Index)
		if e != nil {
			return e
		}
		var da kagome.DoubleArray
		dec := gob.NewDecoder(f)
		if e = dec.Decode(&da); e != nil {
			return e
		}
		dic.Index = &da
		if e = dec.Decode(&dic.IndexDup); e != nil {
			return e
		}
		return nil
	}(); err != nil {
		return
	}
	if err = func() error {
		f, e := os.Open(path.Connection)
		if e != nil {
			return e
		}
		dec := gob.NewDecoder(f)
		if e = dec.Decode(&dic.Connection); e != nil {
			return e
		}
		return nil
	}(); err != nil {
		return
	}

	if err = func() error {
		f, e := os.Open(path.CharDef)
		if e != nil {
			return e
		}
		dec := gob.NewDecoder(f)
		if e = dec.Decode(&dic.CharClass); e != nil {
			return e
		}
		if e = dec.Decode(&dic.CharCategory); e != nil {
			return e
		}
		if e = dec.Decode(&dic.InvokeList); e != nil {
			return e
		}
		if e = dec.Decode(&dic.GroupList); e != nil {
			return e
		}
		return nil
	}(); err != nil {
		return
	}

	if err = func() error {
		f, e := os.Open(path.Unk)
		if e != nil {
			return e
		}
		dec := gob.NewDecoder(f)
		if e = dec.Decode(&dic.UnkMorphs); e != nil {
			return e
		}
		if e = dec.Decode(&dic.UnkIndex); e != nil {
			return e
		}
		if e = dec.Decode(&dic.UnkIndexDup); e != nil {
			return e
		}
		if e = dec.Decode(&dic.UnkContents); e != nil {
			return e
		}
		return nil
	}(); err != nil {
		return
	}
	return
}

var usageMessage = "usage: ipadic [-o output_path] [-z] mecab_path"

func usage() {
	fmt.Fprintln(os.Stderr, usageMessage)
	flag.PrintDefaults()
	os.Exit(0)
}

var (
	fOutputPath = flag.String("o", ".", "the path to output files")
	fArchive    = flag.Bool("a", false, "archive files")
)

func main() {
	flag.Parse()
	mecabPath := flag.Arg(0)
	if mecabPath == "" {
		usage()
	}
	d, err := buildIpaDic(mecabPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "build error: %v\n", err)
	}
	err = saveIpaDic(d, *fOutputPath, *fArchive)
	if err != nil {
		fmt.Fprintf(os.Stderr, "build error: %v\n", err)
	}
}

type ipaMorphRecordSlice [][]string

func (p ipaMorphRecordSlice) Len() int           { return len(p) }
func (p ipaMorphRecordSlice) Less(i, j int) bool { return p[i][0] < p[j][0] }
func (p ipaMorphRecordSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func saveIpaDic(dic *IpaDic, base string, archive bool) (err error) {
	var zw *zip.Writer
	if archive {
		p := path.Join(base, ipaDicArchiveFileName)
		f, e := os.Create(p)
		if e != nil {
			return e
		}
		defer f.Close()
		zw = zip.NewWriter(f)
	}

	if err = func() (e error) {
		p := path.Join(base, ipaDicMorphFileName)
		var out io.Writer
		if archive {
			out, e = zw.Create(p)
			if e != nil {
				return
			}
		} else {
			var f *os.File
			if f, e = os.OpenFile(p, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666); e != nil {
				return
			}
			defer f.Close()
			out = f
		}
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if e = enc.Encode(dic.Morphs); e != nil {
			return
		}
		if _, e = buf.WriteTo(out); e != nil {
			return
		}
		if e = enc.Encode(dic.Contents); e != nil {
			return
		}
		if _, e = buf.WriteTo(out); e != nil {
			return
		}
		return
	}(); err != nil {
		return
	}

	if err = func() (e error) {
		p := path.Join(base, ipaDicIndexFileName)
		var out io.Writer
		if archive {
			if out, e = zw.Create(p); e != nil {
				return
			}

		} else {
			var f *os.File
			if f, e = os.OpenFile(p, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666); e != nil {
				return
			}
			defer f.Close()
			out = f
		}
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if e != nil {
			return e
		}
		if e = enc.Encode(dic.Index); e != nil {
			return e
		}
		if _, e = buf.WriteTo(out); e != nil {
			return e
		}
		if e = enc.Encode(dic.IndexDup); e != nil {
			return e
		}
		if _, e = buf.WriteTo(out); e != nil {
			return e
		}
		return nil
	}(); err != nil {
		return
	}

	if err = func() (e error) {
		p := path.Join(base, ipaDicConnectionFileName)
		var out io.Writer
		if archive {
			if out, e = zw.Create(p); e != nil {
				return
			}

		} else {
			var f *os.File
			if f, e = os.OpenFile(p, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666); e != nil {
				return
			}
			defer f.Close()
			out = f
		}
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if e = enc.Encode(dic.Connection); e != nil {
			return e
		}
		if _, e = buf.WriteTo(out); e != nil {
			return e
		}
		return e
	}(); err != nil {
		return
	}

	if err = func() (e error) {
		p := path.Join(base, ipaDicCharDefFileName)
		var out io.Writer
		if archive {
			if out, e = zw.Create(p); e != nil {
				return
			}

		} else {
			var f *os.File
			if f, e = os.OpenFile(p, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666); e != nil {
				return
			}
			defer f.Close()
			out = f
		}
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if e = enc.Encode(dic.CharClass); e != nil {
			return e
		}
		if _, e = buf.WriteTo(out); e != nil {
			return e
		}
		if e = enc.Encode(dic.CharCategory); e != nil {
			return e
		}
		if _, e = buf.WriteTo(out); e != nil {
			return e
		}
		if e = enc.Encode(dic.InvokeList); e != nil {
			return e
		}
		if _, e = buf.WriteTo(out); e != nil {
			return e
		}
		if e = enc.Encode(dic.GroupList); e != nil {
			return e
		}
		if _, e = buf.WriteTo(out); e != nil {
			return e
		}
		return nil
	}(); err != nil {
		return
	}

	if err = func() (e error) {
		p := path.Join(base, ipaDicUnkFileName)
		var out io.Writer
		if archive {
			if out, e = zw.Create(p); e != nil {
				return
			}

		} else {
			var f *os.File
			if f, e = os.OpenFile(p, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666); e != nil {
				return
			}
			defer f.Close()
			out = f
		}
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if e = enc.Encode(dic.UnkMorphs); e != nil {
			return e
		}
		if _, e = buf.WriteTo(out); e != nil {
			return e
		}
		if e = enc.Encode(dic.UnkIndex); e != nil {
			return e
		}
		if _, e = buf.WriteTo(out); e != nil {
			return e
		}
		if err = enc.Encode(dic.UnkIndexDup); err != nil {
			return e
		}
		if _, e = buf.WriteTo(out); e != nil {
			return e
		}
		if e = enc.Encode(dic.UnkContents); e != nil {
			return e
		}
		if _, e = buf.WriteTo(out); e != nil {
			return e
		}
		return nil
	}(); err != nil {
		return
	}

	if archive {
		err = zw.Close()
	}
	return
}

func buildIpaDic(path string) (dic *IpaDic, err error) {
	// Morphs, Contents, Index
	var files []string
	files, err = filepath.Glob(path + "/*.csv")
	if err != nil {
		return
	}
	var records ipaMorphRecordSlice
	for _, file := range files {
		if err = func() error {
			f, e := os.Open(file)
			if e != nil {
				return e
			}
			defer f.Close()
			r := csv.NewReader(transform.NewReader(f, japanese.EUCJP.NewDecoder()))
			r.Comma = ','
			for {
				rec, e := r.Read()
				if e == io.EOF {
					break
				} else if e != nil {
					return e
				} else if len(rec) != ipaMorphCsvColSize {
					return fmt.Errorf("invalid format csv: %v, %v", file, rec)
				}
				records = append(records, rec)
			}
			return nil
		}(); err != nil {
			return
		}
	}
	sort.Sort(records)
	dic = new(IpaDic)
	dic.Morphs = make([]kagome.Morph, 0, len(records))
	dic.Contents = make([][]string, 0, len(records))
	dic.IndexDup = make(map[int]int)
	keywords := make([]string, 0, len(records))
	ids := make([]int, 0, len(records))
	prev := ""
	for i, rec := range records {
		if prev == rec[ipaMrophRecordSurfaceIndex] {
			prevId := ids[len(ids)-1]
			if _, ok := dic.IndexDup[prevId]; !ok {
				dic.IndexDup[prevId]++
			}
			dic.IndexDup[prevId]++
		} else {
			keywords = append(keywords, rec[ipaMrophRecordSurfaceIndex])
			ids = append(ids, i)
			prev = rec[ipaMrophRecordSurfaceIndex]
		}
		var l, r, w int
		if l, err = strconv.Atoi(rec[ipaMorphRecordLeftIdIndex]); err != nil {
			return
		}
		if r, err = strconv.Atoi(rec[ipaMorphRecordRightIdIndex]); err != nil {
			return
		}
		if w, err = strconv.Atoi(rec[ipaMorphRecordWeightIndex]); err != nil {
			return
		}
		m := kagome.Morph{LeftId: int16(l), RightId: int16(r), Weight: int16(w)}
		dic.Morphs = append(dic.Morphs, m)
		dic.Contents = append(dic.Contents, rec[ipaMorphRecordOtherContentsStartIndex:])
	}
	da := &kagome.DoubleArray{}
	err = da.BuildWithIds(keywords, ids)
	if err != nil {
		return
	}
	dic.Index = da

	// ConnectionTable
	if r, c, v, e := loadIpaMatrixDefFile(path + "/" + ipaMatrixDefFileName); e != nil {
		err = e
		return
	} else {
		dic.Connection.Row = r
		dic.Connection.Col = c
		dic.Connection.Vec = v
	}

	// CharDef
	if cc, cm, inv, grp, e := loadIpaCharClassDefFile(path + "/" + ipaCharDefFileName); e != nil {
		err = e
		return
	} else {
		dic.CharClass = cc
		dic.CharCategory = cm
		dic.InvokeList = inv
		dic.GroupList = grp
	}

	// Unk
	if records, e := loadIpaUnkFile(path + "/" + ipaUnkDefFileName); e != nil {
		err = e
		return
	} else {
		dic.UnkIndex = make(map[int]int)
		dic.UnkIndexDup = make(map[int]int)
		sort.Sort(ipaMorphRecordSlice(records))
		for _, rec := range records {
			catid := -1
			for id, cat := range dic.CharClass {
				if cat == rec[ipaUnkRecordCategoryIndex] {
					catid = id
					break
				}
			}
			if catid < 0 {
				err = fmt.Errorf("unknown unk category: %v", rec[ipaUnkRecordCategoryIndex])
				return
			}
			if _, ok := dic.UnkIndex[catid]; !ok {
				dic.UnkIndex[catid] = len(dic.UnkContents)
			} else {
				if _, ok := dic.UnkIndexDup[catid]; !ok {
					dic.UnkIndexDup[catid]++
				}
				dic.UnkIndexDup[catid]++
			}
			var l, r, w int
			if l, err = strconv.Atoi(rec[ipaUnkRecordLeftIdIndex]); err != nil {
				return
			}
			if r, err = strconv.Atoi(rec[ipaUnkRecordRightIndex]); err != nil {
				return
			}
			if w, err = strconv.Atoi(rec[ipaUnkRecordWeigthIndex]); err != nil {
				return
			}
			m := kagome.Morph{LeftId: int16(l), RightId: int16(r), Weight: int16(w)}
			dic.UnkMorphs = append(dic.UnkMorphs, m)
			dic.UnkContents = append(dic.UnkContents, rec[ipaUnkRecordOtherContentsStartIndex:])
		}
	}
	return
}

func loadIpaMorphFile(path string) (records [][]string, err error) {
	var f *os.File
	f, err = os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	r := csv.NewReader(transform.NewReader(f, japanese.EUCJP.NewDecoder()))
	r.Comma = ','
	for {
		record, e := r.Read()
		if e == io.EOF {
			break
		} else if e != nil {
			err = e
			return
		}
		records = append(records, record)
	}
	return
}

func loadIpaMatrixDefFile(path string) (rowSize, colSize int, vec []int16, err error) {
	var file *os.File
	file, err = os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	dim := strings.Split(line, " ")
	if len(dim) != 2 {
		err = fmt.Errorf("invalid format: %s", line)
		return
	}
	rowSize, err = strconv.Atoi(dim[0])
	if err != nil {
		err = fmt.Errorf("invalid format: %s, %s", err, line)
		return
	}
	colSize, err = strconv.Atoi(dim[1])
	if err != nil {
		err = fmt.Errorf("invalid format: %s, %s", err, line)
		return
	}
	vec = make([]int16, rowSize*colSize)
	for scanner.Scan() {
		line := scanner.Text()
		ary := strings.Split(line, " ")
		if len(ary) != 3 {
			err = fmt.Errorf("invalid format: %s", line)
			return
		}
		row, e := strconv.Atoi(ary[0])
		if e != nil {
			err = fmt.Errorf("invalid format: %s, %s", e, line)
			return
		}
		col, e := strconv.Atoi(ary[1])
		if e != nil {
			err = fmt.Errorf("invalid format: %s, %s", e, line)
			return
		}
		val, e := strconv.Atoi(ary[2])
		if e != nil {
			err = fmt.Errorf("invalid format: %s, %s", e, line)
			return
		}
		vec[row*colSize+col] = int16(val)
	}
	if err = scanner.Err(); err != nil {
		err = fmt.Errorf("invalid format: %s, %s", err, line)
		return
	}
	return
}

func loadIpaCharClassDefFile(path string) (charClass []string, charCategory []byte, invokeMap, groupMap []bool, err error) {
	var file *os.File
	file, err = os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	charCategory = make([]byte, 65536)
	charClass = make([]string, 0)

	regCharClass := regexp.MustCompile("^(\\w+)\\s+(\\d+)\\s+(\\d+)\\s+(\\d+)")
	regCharCategory := regexp.MustCompile("^(0x[0-9A-F]+)(?:\\s+([^#\\s]+))(?:\\s+([^#\\s]+))?")
	regCharCategoryRange := regexp.MustCompile("^(0x[0-9A-F]+)..(0x[0-9A-F]+)(?:\\s+([^#\\s]+))(?:\\s+([^#\\s]+))?")

	scanner := bufio.NewScanner(file)
	cc2id := make(map[string]byte)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		line = strings.TrimSpace(line)

		if matches := regCharClass.FindStringSubmatch(line); len(matches) > 0 {
			invokeMap = append(invokeMap, matches[2] == "1")
			groupMap = append(groupMap, matches[3] == "1")
			cc2id[matches[1]] = byte(len(charClass))
			charClass = append(charClass, matches[1])
		} else if matches := regCharCategory.FindStringSubmatch(line); len(matches) > 0 {
			ch, _ := strconv.ParseInt(matches[1], 0, 0)
			charCategory[ch] = cc2id[matches[2]]
		} else if matches := regCharCategoryRange.FindStringSubmatch(line); len(matches) > 0 {
			start, _ := strconv.ParseInt(matches[1], 0, 0)
			end, _ := strconv.ParseInt(matches[2], 0, 0)
			for x := start; x <= end; x++ {
				charCategory[x] = cc2id[matches[3]]
			}
		} else {
			err = fmt.Errorf("invalid format error: %v", line)
			return
		}

	}
	err = scanner.Err()
	return
}

func loadIpaUnkFile(path string) (records [][]string, err error) {
	var file *os.File
	file, err = os.Open(path)
	if err != nil {
		return
	}
	r := csv.NewReader(transform.NewReader(file, japanese.EUCJP.NewDecoder()))
	r.Comma = ','
	for {
		rec, e := r.Read()
		if e == io.EOF {
			break
		} else if e != nil {
			err = e
			return
		} else if len(rec) != ipaUnkRecordSize {
			err = fmt.Errorf("invalid format csv: %v, %v", file, rec)
			return
		}
		records = append(records, rec)
	}
	return
}
