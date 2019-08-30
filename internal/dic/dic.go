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

package dic

import (
	"archive/zip"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
)

// Dic represents a dictionary of a tokenizer.
type Dic struct {
	Morphs       []Morph
	POSTable     POSTable
	Contents     [][]string
	Connection   ConnectionTable
	Index        IndexTable
	CharClass    []string
	CharCategory []byte
	InvokeList   []bool
	GroupList    []bool

	UnkDic
}

// CharacterCategory returns the category of a rune.
func (d Dic) CharacterCategory(r rune) byte {
	if int(r) < len(d.CharCategory) {
		return d.CharCategory[r]
	}
	return d.CharCategory[0] // default
}

func (d *Dic) loadMorphDicPart(r io.Reader) error {
	m, err := LoadMorphSlice(r)
	if err != nil {
		return fmt.Errorf("dic initializer, Morphs: %v", err)
	}
	d.Morphs = m
	return nil
}

func (d *Dic) loadPOSDicPart(r io.Reader) error {
	p, err := ReadPOSTable(r)
	if err != nil {
		return fmt.Errorf("dic initializer, POSs: %v", err)
	}
	d.POSTable = p
	return nil
}

func (d *Dic) loadContentDicPart(r io.Reader) error {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("dic initializer, Contents: %v", err)
	}
	d.Contents = NewContents(buf)
	return nil
}

func (d *Dic) loadIndexDicPart(r io.Reader) error {
	idx, err := ReadIndexTable(r)
	if err != nil {
		return fmt.Errorf("dic initializer, Index: %v", err)
	}
	d.Index = idx
	return nil
}

func (d *Dic) loadConnectionDicPart(r io.Reader) error {
	t, err := LoadConnectionTable(r)
	if err != nil {
		return fmt.Errorf("dic initializer, Connection: %v", err)
	}
	d.Connection = t
	return nil
}

func (d *Dic) loadCharDefDicPart(r io.Reader) error {
	dec := gob.NewDecoder(r)
	if err := dec.Decode(&d.CharClass); err != nil {
		return fmt.Errorf("dic initializer, CharClass: %v", err)
	}
	if err := dec.Decode(&d.CharCategory); err != nil {
		return fmt.Errorf("dic initializer, CharCategory: %v", err)
	}
	if err := dec.Decode(&d.InvokeList); err != nil {
		return fmt.Errorf("dic initializer, InvokeList: %v", err)
	}
	if err := dec.Decode(&d.GroupList); err != nil {
		return fmt.Errorf("dic initializer, GroupList: %v", err)
	}
	return nil
}

func (d *Dic) loadUnkDicPart(r io.Reader) error {
	unk, err := ReadUnkDic(r)
	if err != nil {
		return fmt.Errorf("dic initializer, UnkDic: %v", err)
	}
	d.UnkDic = unk
	return nil
}

// Load loads a dictionary from a file.
func Load(path string) (d *Dic, err error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return d, err
	}
	defer r.Close()
	return load(&r.Reader, true)
}

// LoadSimple loads a dictionary from a file without contents.
func LoadSimple(path string) (d *Dic, err error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return d, err
	}
	defer r.Close()
	return load(&r.Reader, false)
}

func load(r *zip.Reader, full bool) (*Dic, error) {
	var d Dic
	type dictionaryPartLoader func(io.Reader) error
	loaders := map[string]dictionaryPartLoader{
		"morph.dic":      d.loadMorphDicPart,
		"pos.dic":        d.loadPOSDicPart,
		"content.dic":    d.loadContentDicPart,
		"index.dic":      d.loadIndexDicPart,
		"connection.dic": d.loadConnectionDicPart,
		"chardef.dic":    d.loadCharDefDicPart,
		"unk.dic":        d.loadUnkDicPart,
	}

	for _, f := range r.File {
		if err := func() error {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()
			if !full && f.Name == "content.dic" {
				return nil
			}
			loader, ok := loaders[f.Name]
			if !ok {
				return fmt.Errorf("unknown file, %v", f.Name)
			}
			return loader(rc)
		}(); err != nil {
			return nil, err
		}
	}
	return &d, nil
}
