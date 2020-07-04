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
	"compress/flate"
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

	"github.com/ikawaha/kagome/cmd/_dictool/splitfile"
	"github.com/ikawaha/kagome/tokenizer/dic"
)

const (
	uniMatrixDefFileName = "matrix.def"
	uniCharDefFileName   = "char.def"
	uniUnkDefFileName    = "unk.def"

	uniDicArchiveFileName    = "uni.dic"
	uniDicMorphFileName      = "morph.dic"
	uniDicPOSFileName        = "pos.dic"
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
	uniMorphRecordPOSRecordStartIndex     = 4
	uniMorphRecordOtherContentsStartIndex = 10

	uniUnkRecordSize                    = 10
	uniUnkRecordCategoryIndex           = 0
	uniUnkRecordLeftIDIndex             = 1
	uniUnkRecordRightIndex              = 2
	uniUnkRecordWeigthIndex             = 3
	uniUnkRecordOtherContentsStartIndex = 4
)

type UniDic struct {
	Morphs       []dic.Morph
	POSTable     dic.POSTable
	Contents     [][]string
	Index        dic.IndexTable
	Connection   dic.ConnectionTable
	CharClass    []string
	CharCategory []byte
	InvokeList   []bool
	GroupList    []bool

	dic.UnkDic
}

type uniDicPath struct {
	Morph      string
	Index      string
	Connection string
	CharDef    string
	Unk        string
}

type uniMorphRecordSlice [][]string

func (p uniMorphRecordSlice) Len() int           { return len(p) }
func (p uniMorphRecordSlice) Less(i, j int) bool { return p[i][0] < p[j][0] }
func (p uniMorphRecordSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func saveUniDic(d *UniDic, base string, archive bool) error {
	var zw *zip.Writer
	p := path.Join(base, uniDicArchiveFileName)
	if archive {
		f, err := os.Create(p)
		if err != nil {
			return err
		}
		defer f.Close()
		zw = zip.NewWriter(f)
	} else {
		f, err := splitfile.Open(p, 10*1024*1024) // 10MB
		if err != nil {
			return err
		}
		defer f.Close()
		zw = zip.NewWriter(f)
		zw.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
			return flate.NewWriter(out, flate.NoCompression)
		})
		if err != nil {
			return err
		}

	}

	if err := func() error {
		out, err := zw.Create(uniDicMorphFileName)
		if err != nil {
			return err
		}
		_, err = dic.MorphSlice(d.Morphs).WriteTo(out)
		return err
	}(); err != nil {
		return err
	}

	if err := func() (e error) {
		out, err := zw.Create(uniDicPOSFileName)
		if err != nil {
			return err
		}
		_, err = dic.POSTable(d.POSTable).WriteTo(out)
		return err
	}(); err != nil {
		return err
	}

	if err := func() error {
		out, err := zw.Create(uniDicContentFileName)
		if err != nil {
			return err
		}
		_, err = dic.Contents(d.Contents).WriteTo(out)
		return err
	}(); err != nil {
		return err
	}

	if err := func() error {
		out, err := zw.Create(uniDicIndexFileName)
		if err != nil {
			return err
		}
		_, err = d.Index.WriteTo(out)
		return err
	}(); err != nil {
		return err
	}

	if err := func() error {
		out, err := zw.Create(uniDicConnectionFileName)
		if err != nil {
			return err
		}
		_, err = d.Connection.WriteTo(out)
		return err
	}(); err != nil {
		return err
	}

	if err := func() error {
		out, err := zw.Create(uniDicCharDefFileName)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(d.CharClass); err != nil {
			return err
		}
		if _, err := buf.WriteTo(out); err != nil {
			return err
		}
		if err := enc.Encode(d.CharCategory); err != nil {
			return err
		}
		if _, err := buf.WriteTo(out); err != nil {
			return err
		}
		if err := enc.Encode(d.InvokeList); err != nil {
			return err
		}
		if _, err := buf.WriteTo(out); err != nil {
			return err
		}
		if err := enc.Encode(d.GroupList); err != nil {
			return err
		}
		if _, err := buf.WriteTo(out); err != nil {
			return err
		}
		return nil
	}(); err != nil {
		return err
	}

	if err := func() error {
		out, err := zw.Create(uniDicUnkFileName)
		if err != nil {
			return err
		}
		if _, err := d.UnkDic.WriteTo(out); err != nil {
			return err
		}
		return nil
	}(); err != nil {
		return err
	}

	return zw.Close()
}

func buildUniDic(mecabPath, neologdPath string) (*UniDic, error) {
	// Morphs, Contents, Index
	files, err := filepath.Glob(mecabPath + "/*.csv")
	if err != nil {
		return nil, err
	}
	var records uniMorphRecordSlice
	for _, file := range files {
		if err := func() error {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()
			r := csv.NewReader(f)
			r.Comma = ','
			for {
				rec, err := r.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					return err
				} else if len(rec) != uniMorphCsvColSize {
					return fmt.Errorf("invalid format csv: %v, %v", file, rec)
				}
				records = append(records, rec)
			}
			return nil
		}(); err != nil {
			return nil, err
		}
	}
	if err := func() error {
		if neologdPath == "" {
			return nil
		}
		f, err := os.Open(neologdPath)
		if err != nil {
			return err
		}
		defer f.Close()
		r := csv.NewReader(f)
		r.Comma = ','
		r.LazyQuotes = true
		for {
			rec, err := r.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			} else if len(rec) != uniMorphCsvColSize {
				return fmt.Errorf("invalid format csv: %v, %v", neologdPath, rec)
			}
			records = append(records, rec)
		}
		return nil
	}(); err != nil {
		return nil, err
	}

	sort.Sort(records)
	var d UniDic
	d.Morphs = make([]dic.Morph, 0, len(records))
	d.POSTable = dic.POSTable{
		POSs: make([]dic.POS, 0, len(records)),
	}
	d.Contents = make([][]string, 0, len(records))
	var (
		keywords []string
		posMap   = make(dic.POSMap)
	)
	for _, rec := range records {
		keywords = append(keywords, rec[uniMrophRecordSurfaceIndex])
		l, err := strconv.Atoi(rec[uniMorphRecordLeftIDIndex])
		if err != nil {
			return nil, err
		}
		r, err := strconv.Atoi(rec[uniMorphRecordRightIDIndex])
		if err != nil {
			return nil, err
		}
		w, err := strconv.Atoi(rec[uniMorphRecordWeightIndex])
		if err != nil {
			return nil, err
		}
		m := dic.Morph{LeftID: int16(l), RightID: int16(r), Weight: int16(w)}
		d.Morphs = append(d.Morphs, m)
		d.POSTable.POSs = append(d.POSTable.POSs, posMap.Add(
			rec[uniMorphRecordPOSRecordStartIndex:uniMorphRecordOtherContentsStartIndex]),
		)
		d.Contents = append(d.Contents, rec[uniMorphRecordOtherContentsStartIndex:])
	}
	d.POSTable.NameList = posMap.List()

	if d.Index, err = dic.BuildIndexTable(keywords); err != nil {
		return nil, err
	}

	// ConnectionTable
	if r, c, v, err := loadUniMatrixDefFile(mecabPath + "/" + uniMatrixDefFileName); err != nil {
		return nil, err
	} else {
		d.Connection.Row = r
		d.Connection.Col = c
		d.Connection.Vec = v
	}

	// CharDef
	if cc, cm, inv, grp, err := loadUniCharClassDefFile(mecabPath + "/" + uniCharDefFileName); err != nil {
		return nil, err
	} else {
		d.CharClass = cc
		d.CharCategory = cm
		d.InvokeList = inv
		d.GroupList = grp
	}

	// Unk
	if records, err := loadUniUnkFile(mecabPath + "/" + uniUnkDefFileName); err != nil {
		return nil, err
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
				return nil, fmt.Errorf("unknown unk category: %v", rec[uniUnkRecordCategoryIndex])

			}
			if _, ok := d.UnkIndex[catid]; !ok {
				d.UnkIndex[catid] = int32(len(d.UnkContents))
			} else {
				d.UnkIndexDup[catid]++
			}
			l, err := strconv.Atoi(rec[uniUnkRecordLeftIDIndex])
			if err != nil {
				return nil, err
			}
			r, err := strconv.Atoi(rec[uniUnkRecordRightIndex])
			if err != nil {
				return nil, err
			}
			w, err := strconv.Atoi(rec[uniUnkRecordWeigthIndex])
			if err != nil {
				return nil, err
			}
			m := dic.Morph{LeftID: int16(l), RightID: int16(r), Weight: int16(w)}
			d.UnkMorphs = append(d.UnkMorphs, m)
			d.UnkContents = append(d.UnkContents, rec[uniUnkRecordOtherContentsStartIndex:])
		}
	}
	return &d, nil
}

func loadUniMorphFile(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.Comma = ','
	var records [][]string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func loadUniMatrixDefFile(path string) (rowSize, colSize int64, vec []int16, err error) {
	var file *os.File
	file, err = os.Open(path)
	if err != nil {
		return rowSize, colSize, vec, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	dim := strings.Split(line, " ")
	if len(dim) != 2 {
		return rowSize, colSize, vec, fmt.Errorf("invalid format: %s", line)
	}
	rowSize, err = strconv.ParseInt(dim[0], 10, 0)
	if err != nil {
		return rowSize, colSize, vec, fmt.Errorf("invalid format: %s, %s", err, line)
	}
	colSize, err = strconv.ParseInt(dim[1], 10, 0)
	if err != nil {
		return rowSize, colSize, vec, fmt.Errorf("invalid format: %s, %s", err, line)
	}
	vec = make([]int16, rowSize*colSize)
	for scanner.Scan() {
		line := scanner.Text()
		ary := strings.Split(line, " ")
		if len(ary) != 3 {
			return rowSize, colSize, vec, fmt.Errorf("invalid format: %s", line)
		}
		row, err := strconv.ParseInt(ary[0], 10, 0)
		if err != nil {
			return rowSize, colSize, vec, fmt.Errorf("invalid format: %s, %s", err, line)
		}
		col, err := strconv.ParseInt(ary[1], 10, 0)
		if err != nil {
			return rowSize, colSize, vec, fmt.Errorf("invalid format: %s, %s", err, line)
		}
		val, err := strconv.Atoi(ary[2])
		if err != nil {
			return rowSize, colSize, vec, fmt.Errorf("invalid format: %s, %s", err, line)
		}
		vec[row*colSize+col] = int16(val)
	}
	if err := scanner.Err(); err != nil {
		return rowSize, colSize, vec, fmt.Errorf("invalid format: %s, %s", err, line)
	}
	return rowSize, colSize, vec, nil
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

func loadUniUnkFile(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.Comma = ','
	var records [][]string
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else if len(rec) != uniUnkRecordSize {
			return nil, fmt.Errorf("invalid format csv: %v, %v", f, rec)
		}
		records = append(records, rec)
	}
	return records, nil
}
