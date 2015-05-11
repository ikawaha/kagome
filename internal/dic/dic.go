//  Copyright (c) 2015 ikawaha.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package dic

import (
	"archive/zip"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/ikawaha/kagome/internal/fst"
)

// Dic represents a dictionary of a tokenizer.
type Dic struct {
	Morphs       []Morph
	Contents     [][]string
	Connection   ConnectionTable
	Index        Trie
	CharClass    []string
	CharCategory []byte
	InvokeList   []bool
	GroupList    []bool

	UnkMorphs   []Morph
	UnkIndex    map[int]int
	UnkIndexDup map[int]int
	UnkContents [][]string
}

// CharactorCategory returns the category of a rune.
func (d Dic) CharactorCategory(r rune) byte {
	if int(r) <= len(d.CharCategory) {
		return d.CharCategory[r]
	}
	return d.CharCategory[0] // default
}

func (d *Dic) loadMorphDicPart(r io.Reader) error {
	m, e := LoadMorphSlice(r)
	if e != nil {
		return fmt.Errorf("dic initializer, Morphs: %v", e)
	}
	d.Morphs = m
	dec := gob.NewDecoder(r)
	if e := dec.Decode(&d.Contents); e != nil {
		return fmt.Errorf("dic initializer, Contents: %v", e)
	}
	return nil
}

func (d *Dic) loadIndexDicPart(r io.Reader) error {
	t, e := fst.Read(r)
	if e != nil {
		return fmt.Errorf("dic initializer, Index: %v", e)
	}
	d.Index = t
	return nil
}

func (d *Dic) loadConnectionDicPart(r io.Reader) error {
	t, e := LoadConnectionTable(r)
	if e != nil {
		return fmt.Errorf("dic initializer, Connection: %v", e)
	}
	d.Connection = t
	return nil
}

func (d *Dic) loadCharDefDicPart(r io.Reader) error {
	dec := gob.NewDecoder(r)
	if e := dec.Decode(&d.CharClass); e != nil {
		return fmt.Errorf("dic initializer, CharClass: %v", e)
	}
	if e := dec.Decode(&d.CharCategory); e != nil {
		return fmt.Errorf("dic initializer, CharCategory: %v", e)
	}
	if e := dec.Decode(&d.InvokeList); e != nil {
		return fmt.Errorf("dic initializer, InvokeList: %v", e)
	}
	if e := dec.Decode(&d.GroupList); e != nil {
		return fmt.Errorf("dic initializer, GroupList: %v", e)
	}
	return nil
}

func (d *Dic) loadUnkDicPart(r io.Reader) error {
	dec := gob.NewDecoder(r)
	if e := dec.Decode(&d.UnkMorphs); e != nil {
		return fmt.Errorf("dic initializer, UnkMorphs: %v", e)
	}
	if e := dec.Decode(&d.UnkIndex); e != nil {
		return fmt.Errorf("dic initializer, UnkIndex: %v", e)
	}
	if e := dec.Decode(&d.UnkIndexDup); e != nil {
		return fmt.Errorf("dic initializer, UnkIndexDup: %v", e)
	}
	if e := dec.Decode(&d.UnkContents); e != nil {
		return fmt.Errorf("dic initializer, UnkContents: %v", e)
	}
	return nil
}

// Load loads a dictionary from a file.
func Load(path string) (d *Dic, err error) {
	d = new(Dic)
	r, err := zip.OpenReader(path)
	if err != nil {
		return
	}
	defer r.Close()

	for _, f := range r.File {
		if err = func() error {
			rc, e := f.Open()
			if e != nil {
				return e
			}
			rc.Close()
			switch f.Name {
			case "morph.dic":
				if e = d.loadMorphDicPart(rc); err != nil {
					return e
				}
			case "index.dic":
				if e = d.loadIndexDicPart(rc); err != nil {
					return e
				}
			case "connection.dic":
				if e = d.loadConnectionDicPart(rc); err != nil {
					return e
				}
			case "chardef.dic":
				if e = d.loadCharDefDicPart(rc); err != nil {
					return e
				}
			case "unk.dic":
				if e = d.loadUnkDicPart(rc); err != nil {
					return e
				}
			}
			return nil
		}(); err != nil {
			return
		}
	}
	return
}
