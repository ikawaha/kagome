// Copyright 2015 ikawaha
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// 	You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package uni

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ikawaha/kagome/internal/dic"
)

const (
	uniMatrixDefFileName = "matrix.def"
	uniCharDefFileName   = "char.def"
	uniUnkDefFileName    = "unk.def"

	uniDicArchiveFileName    = "uni.dic"
	uniDicMorphFileName      = "morph.dic"
	uniDicContentFileName    = "content.dic"
	uniDicIndexFileName      = "index.dic"
	uniDicConnectionFileName = "connection.dic"
	uniDicCharDefFileName    = "chardef.dic"
	uniDicUnkFileName        = "unk.dic"

	uniMorphCsvColSize                    = 21
	uniMrophRecordSurfaceIndex            = 0
	uniMorphRecordLeftIDIndex             = 1
	uniMorphRecordRightIDIndex            = 2
	uniMorphRecordWeightIndex             = 3
	uniMorphRecordOtherContentsStartIndex = 4

	uniUnkRecordSize                    = 10
	uniUnkRecordCategoryIndex           = 0
	uniUnkRecordLeftIDIndex             = 1
	uniUnkRecordRightIndex              = 2
	uniUnkRecordWeigthIndex             = 3
	uniUnkRecordOtherContentsStartIndex = 4
)

type UniDic struct {
	Morphs       []dic.Morph
	Contents     [][]string
	Index        dic.IndexTable
	Connection   dic.ConnectionTable
	CharClass    []string
	CharCategory []byte
	InvokeList   []bool
	GroupList    []bool

	UnkMorphs   []dic.Morph
	UnkIndex    map[int32]int32
	UnkIndexDup map[int32]int32
	UnkContents [][]string
}

type uniDicPath struct {
	Morph      string
	Index      string
	Connection string
	CharDef    string
	Unk        string
}

func loadUniDic(path uniDicPath) (d *UniDic, err error) {
	d = new(UniDic)
	if err = func() error {
		f, e := os.Open(path.Morph)
		if e != nil {
			return e
		}
		dec := gob.NewDecoder(f)
		if e = dec.Decode(&d.Morphs); e != nil {
			return e
		}
		if e = dec.Decode(&d.Contents); e != nil {
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
		idx, e := dic.ReadIndexTable(f)
		if e != nil {
			return e
		}
		d.Index = idx
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
		if e = dec.Decode(&d.Connection); e != nil {
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
		if e = dec.Decode(&d.CharClass); e != nil {
			return e
		}
		if e = dec.Decode(&d.CharCategory); e != nil {
			return e
		}
		if e = dec.Decode(&d.InvokeList); e != nil {
			return e
		}
		if e = dec.Decode(&d.GroupList); e != nil {
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
		if e = dec.Decode(&d.UnkMorphs); e != nil {
			return e
		}
		if e = dec.Decode(&d.UnkIndex); e != nil {
			return e
		}
		if e = dec.Decode(&d.UnkIndexDup); e != nil {
			return e
		}
		if e = dec.Decode(&d.UnkContents); e != nil {
			return e
		}
		return nil
	}(); err != nil {
		return
	}
	return
}

type uniMorphRecordSlice [][]string

func (p uniMorphRecordSlice) Len() int           { return len(p) }
func (p uniMorphRecordSlice) Less(i, j int) bool { return p[i][0] < p[j][0] }
func (p uniMorphRecordSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func saveUniDic(d *UniDic, base string, archive bool) (err error) {
	var zw *zip.Writer
	if archive {
		p := path.Join(base, uniDicArchiveFileName)
		f, e := os.Create(p)
		if e != nil {
			return e
		}
		defer f.Close()
		zw = zip.NewWriter(f)
	}

	if err = func() (e error) {
		p := path.Join(base, uniDicMorphFileName)
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
		if _, e = dic.MorphSlice(d.Morphs).WriteTo(out); e != nil {
			return
		}
		return
	}(); err != nil {
		return
	}

	if err = func() (e error) {
		p := path.Join(base, uniDicContentFileName)
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
		if _, e = dic.Contents(d.Contents).WriteTo(out); e != nil {
			return
		}
		return
	}(); err != nil {
		return
	}

	if err = func() (e error) {
		p := path.Join(base, uniDicIndexFileName)
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
		if _, e := d.Index.WriteTo(out); e != nil {
			return e
		}
		return nil
	}(); err != nil {
		return
	}

	if err = func() (e error) {
		p := path.Join(base, uniDicConnectionFileName)
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
		if _, e = d.Connection.WriteTo(out); e != nil {
			return e
		}
		// var buf bytes.Buffer
		// enc := gob.NewEncoder(&buf)
		// if e = enc.Encode(d.Connection); e != nil {
		// 	return e
		// }
		// if _, e = buf.WriteTo(out); e != nil {
		// 	return e
		// }
		return e
	}(); err != nil {
		return
	}

	if err = func() (e error) {
		p := path.Join(base, uniDicCharDefFileName)
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
		if e = enc.Encode(d.CharClass); e != nil {
			return e
		}
		if _, e = buf.WriteTo(out); e != nil {
			return e
		}
		if e = enc.Encode(d.CharCategory); e != nil {
			return e
		}
		if _, e = buf.WriteTo(out); e != nil {
			return e
		}
		if e = enc.Encode(d.InvokeList); e != nil {
			return e
		}
		if _, e = buf.WriteTo(out); e != nil {
			return e
		}
		if e = enc.Encode(d.GroupList); e != nil {
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
		p := path.Join(base, uniDicUnkFileName)
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
		if e = enc.Encode(d.UnkMorphs); e != nil {
			return e
		}
		if _, e = buf.WriteTo(out); e != nil {
			return e
		}
		if e = enc.Encode(d.UnkIndex); e != nil {
			return e
		}
		if _, e = buf.WriteTo(out); e != nil {
			return e
		}
		if err = enc.Encode(d.UnkIndexDup); err != nil {
			return e
		}
		if _, e = buf.WriteTo(out); e != nil {
			return e
		}
		if e = enc.Encode(d.UnkContents); e != nil {
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

func buildUniDic(mecabPath, neologdPath string) (d *UniDic, err error) {
	// Morphs, Contents, Index
	var files []string
	files, err = filepath.Glob(mecabPath + "/*.csv")
	if err != nil {
		return
	}
	var records uniMorphRecordSlice
	for _, file := range files {
		if err = func() error {
			f, e := os.Open(file)
			if e != nil {
				return e
			}
			defer f.Close()
			r := csv.NewReader(f)
			r.Comma = ','
			for {
				rec, e := r.Read()
				if e == io.EOF {
					break
				} else if e != nil {
					return e
				} else if len(rec) != uniMorphCsvColSize {
					return fmt.Errorf("invalid format csv: %v, %v", file, rec)
				}
				records = append(records, rec)
			}
			return nil
		}(); err != nil {
			return
		}
	}
	if err = func() error {
		if neologdPath == "" {
			return nil
		}
		f, e := os.Open(neologdPath)
		if e != nil {
			return e
		}
		defer f.Close()
		r := csv.NewReader(f)
		r.Comma = ','
		r.LazyQuotes = true
		for {
			rec, e := r.Read()
			if e == io.EOF {
				break
			} else if e != nil {
				return e
			} else if len(rec) != uniMorphCsvColSize {
				return fmt.Errorf("invalid format csv: %v, %v", neologdPath, rec)
			}
			records = append(records, rec)
		}
		return nil
	}(); err != nil {
		return
	}

	sort.Sort(records)
	d = new(UniDic)
	d.Morphs = make([]dic.Morph, 0, len(records))
	d.Contents = make([][]string, 0, len(records))
	var keywords []string
	for _, rec := range records {
		keywords = append(keywords, rec[uniMrophRecordSurfaceIndex])
		var l, r, w int
		if l, err = strconv.Atoi(rec[uniMorphRecordLeftIDIndex]); err != nil {
			return
		}
		if r, err = strconv.Atoi(rec[uniMorphRecordRightIDIndex]); err != nil {
			return
		}
		if w, err = strconv.Atoi(rec[uniMorphRecordWeightIndex]); err != nil {
			return
		}
		m := dic.Morph{LeftID: int16(l), RightID: int16(r), Weight: int16(w)}
		d.Morphs = append(d.Morphs, m)
		d.Contents = append(d.Contents, rec[uniMorphRecordOtherContentsStartIndex:])
	}

	if d.Index, err = dic.BuildIndexTable(keywords); err != nil {
		return
	}

	// ConnectionTable
	if r, c, v, e := loadUniMatrixDefFile(mecabPath + "/" + uniMatrixDefFileName); e != nil {
		err = e
		return
	} else {
		d.Connection.Row = r
		d.Connection.Col = c
		d.Connection.Vec = v
	}

	// CharDef
	if cc, cm, inv, grp, e := loadUniCharClassDefFile(mecabPath + "/" + uniCharDefFileName); e != nil {
		err = e
		return
	} else {
		d.CharClass = cc
		d.CharCategory = cm
		d.InvokeList = inv
		d.GroupList = grp
	}

	// Unk
	if records, e := loadUniUnkFile(mecabPath + "/" + uniUnkDefFileName); e != nil {
		err = e
		return
	} else {
		d.UnkIndex = make(map[int32]int32)
		d.UnkIndexDup = make(map[int32]int32)
		sort.Sort(uniMorphRecordSlice(records))
		for _, rec := range records {
			catid := int32(-1)
			for id, cat := range d.CharClass {
				if cat == rec[uniUnkRecordCategoryIndex] {
					catid = int32(id)
					break
				}
			}
			if catid < 0 {
				err = fmt.Errorf("unknown unk category: %v", rec[uniUnkRecordCategoryIndex])
				return
			}
			if _, ok := d.UnkIndex[catid]; !ok {
				d.UnkIndex[catid] = int32(len(d.UnkContents))
			} else {
				d.UnkIndexDup[catid]++
			}
			var l, r, w int
			if l, err = strconv.Atoi(rec[uniUnkRecordLeftIDIndex]); err != nil {
				return
			}
			if r, err = strconv.Atoi(rec[uniUnkRecordRightIndex]); err != nil {
				return
			}
			if w, err = strconv.Atoi(rec[uniUnkRecordWeigthIndex]); err != nil {
				return
			}
			m := dic.Morph{LeftID: int16(l), RightID: int16(r), Weight: int16(w)}
			d.UnkMorphs = append(d.UnkMorphs, m)
			d.UnkContents = append(d.UnkContents, rec[uniUnkRecordOtherContentsStartIndex:])
		}
	}
	return
}

func loadUniMorphFile(path string) (records [][]string, err error) {
	var f *os.File
	f, err = os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	r := csv.NewReader(f)
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

func loadUniMatrixDefFile(path string) (rowSize, colSize int64, vec []int16, err error) {
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
	rowSize, err = strconv.ParseInt(dim[0], 10, 0)
	if err != nil {
		err = fmt.Errorf("invalid format: %s, %s", err, line)
		return
	}
	colSize, err = strconv.ParseInt(dim[1], 10, 0)
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
		row, e := strconv.ParseInt(ary[0], 10, 0)
		if e != nil {
			err = fmt.Errorf("invalid format: %s, %s", e, line)
			return
		}
		col, e := strconv.ParseInt(ary[1], 10, 0)
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

func loadUniCharClassDefFile(path string) (charClass []string, charCategory []byte, invokeMap, groupMap []bool, err error) {
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

func loadUniUnkFile(path string) (records [][]string, err error) {
	var f *os.File
	f, err = os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.Comma = ','
	for {
		rec, e := r.Read()
		if e == io.EOF {
			break
		} else if e != nil {
			err = e
			return
		} else if len(rec) != uniUnkRecordSize {
			err = fmt.Errorf("invalid format csv: %v, %v", f, rec)
			return
		}
		records = append(records, rec)
	}
	return
}
