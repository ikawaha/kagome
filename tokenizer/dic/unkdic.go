// Copyright 2018 ikawaha
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

package dic

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"sort"
)

// UnkDic represents an unknown word dictionary part.
type UnkDic struct {
	UnkMorphs   []Morph
	UnkIndex    map[int32]int32
	UnkIndexDup map[int32]int32
	UnkContents [][]string
}

func writeMap(w io.Writer, m map[int32]int32) (n int64, err error) {
	keys := make([]int32, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	sz := int64(len(keys))
	if err := binary.Write(w, binary.LittleEndian, sz); err != nil {
		return n, err
	}
	n += int64(binary.Size(sz))
	for _, k := range keys {
		if err := binary.Write(w, binary.LittleEndian, k); err != nil {
			return n, err
		}
		n += int64(binary.Size(k))
		v := m[k]
		if err := binary.Write(w, binary.LittleEndian, v); err != nil {
			return n, err
		}
		n += int64(binary.Size(v))
	}
	return n, err
}

// WriteTo implements the io.WriterTo interface.
func (u UnkDic) WriteTo(w io.Writer) (n int64, err error) {
	size, err := writeMap(w, u.UnkIndex)
	if err != nil {
		return n, err
	}
	n += size
	size, err = writeMap(w, u.UnkIndexDup)
	if err != nil {
		return n, err
	}
	n += size

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(u.UnkMorphs); err != nil {
		return n, err
	}
	size, err = buf.WriteTo(w)
	if err != nil {
		return n, err
	}
	n += size
	if err := enc.Encode(u.UnkContents); err != nil {
		return n, err
	}
	size, err = buf.WriteTo(w)
	if err != nil {
		return n, err
	}
	n += size

	return n, nil
}

func readMap(r io.Reader) (map[int32]int32, error) {
	var sz int64
	if err := binary.Read(r, binary.LittleEndian, &sz); err != nil {
		return nil, err
	}
	m := make(map[int32]int32, sz)
	for i := int64(0); i < sz; i++ {
		var k int32
		if err := binary.Read(r, binary.LittleEndian, &k); err != nil {
			return nil, err
		}
		var v int32
		if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
			return nil, err
		}
		m[k] = v
	}
	return m, nil
}

// ReadUnkDic loads an unknown word dictionary.
func ReadUnkDic(r io.Reader) (UnkDic, error) {
	d := UnkDic{}
	ui, err := readMap(r)
	if err != nil {
		return d, fmt.Errorf("UnkIndex: %v", err)
	}
	d.UnkIndex = ui
	ud, err := readMap(r)
	if err != nil {
		return d, fmt.Errorf("UnkIndexDup: %v", err)
	}
	d.UnkIndexDup = ud

	dec := gob.NewDecoder(r)
	if err := dec.Decode(&d.UnkMorphs); err != nil {
		return d, fmt.Errorf("UnkMorphs: %v", err)
	}
	if err := dec.Decode(&d.UnkContents); err != nil {
		return d, fmt.Errorf("UnkContents: %v", err)
	}
	return d, nil
}
