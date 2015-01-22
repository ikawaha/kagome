//  Copyright (c) 2014 ikawaha.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package kagome

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"

	"github.com/ikawaha/kagome/data"
)

const (
	sysDicType = "IPA"
	sysDicPath = "dic/ipa"
)

// Morph represents part of speeches and an occurrence cost.
type Morph struct {
	LeftId, RightId, Weight int16
}

// Dic represents a dictionary of a tokenizer.
type Dic struct {
	Morphs       []Morph
	Contents     [][]string
	Connection   ConnectionTable
	Index        FST
	CharClass    []string
	CharCategory []byte
	InvokeList   []bool
	GroupList    []bool

	UnkMorphs   []Morph
	UnkIndex    map[int]int
	UnkIndexDup map[int]int
	UnkContents [][]string
}

var (
	sysDic     *Dic
	initSysDic sync.Once
)

// NewSysDic returns the kagome system dictionary.
func NewSysDic() (dic *Dic) {
	initSysDic.Do(func() { sysDic = loadSysDic() })
	return sysDic
}

func (d Dic) charCategory(r rune) byte {
	if int(r) <= len(d.CharCategory) {
		return d.CharCategory[r]
	}
	return d.CharCategory[0] //XXX
}

func loadSysDic() (d *Dic) {
	d = new(Dic)
	if err := func() error {
		buf, e := data.Asset(sysDicPath + "/morph.dic")
		if e != nil {
			return e
		}
		dec := gob.NewDecoder(bytes.NewBuffer(buf))
		if e = dec.Decode(&d.Morphs); e != nil {
			return fmt.Errorf("dic initializer, Morphs: %v", e)
		}
		if e = dec.Decode(&d.Contents); e != nil {
			return fmt.Errorf("dic initializer, Contents: %v", e)
		}
		return nil
	}(); err != nil {
		panic(err)
	}

	if err := func() error {
		buf, e := data.Asset(sysDicPath + "/index.dic")
		if e != nil {
			return e
		}
		if e = d.Index.Read(bytes.NewReader(buf)); e != nil {
			return fmt.Errorf("dic initializer, Index: %v", e)
		}
		return nil
	}(); err != nil {
		panic(err)
	}

	if err := func() error {
		buf, e := data.Asset(sysDicPath + "/connection.dic")
		if e != nil {
			return e
		}
		dec := gob.NewDecoder(bytes.NewBuffer(buf))
		if e = dec.Decode(&d.Connection); e != nil {
			return fmt.Errorf("dic initializer, Connection: %v", e)
		}
		return nil
	}(); err != nil {
		panic(err)
	}

	if err := func() error {
		buf, e := data.Asset(sysDicPath + "/chardef.dic")
		if e != nil {
			return e
		}
		dec := gob.NewDecoder(bytes.NewBuffer(buf))
		if e = dec.Decode(&d.CharClass); e != nil {
			return fmt.Errorf("dic initializer, CharClass: %v", e)
		}
		if e = dec.Decode(&d.CharCategory); e != nil {
			return fmt.Errorf("dic initializer, CharCategory: %v", e)
		}
		if e = dec.Decode(&d.InvokeList); e != nil {
			return fmt.Errorf("dic initializer, InvokeList: %v", e)
		}
		if e = dec.Decode(&d.GroupList); e != nil {
			return fmt.Errorf("dic initializer, GroupList: %v", e)
		}
		return nil
	}(); err != nil {
		panic(err)
	}

	if err := func() error {
		buf, e := data.Asset(sysDicPath + "/unk.dic")
		if e != nil {
			return e
		}
		dec := gob.NewDecoder(bytes.NewBuffer(buf))
		if e = dec.Decode(&d.UnkMorphs); e != nil {
			return fmt.Errorf("dic initializer, UnkMorphs: %v", e)
		}
		if e = dec.Decode(&d.UnkIndex); e != nil {
			return fmt.Errorf("dic initializer, UnkIndex: %v", e)
		}
		if e = dec.Decode(&d.UnkIndexDup); e != nil {
			return fmt.Errorf("dic initializer, UnkIndexDup: %v", e)
		}
		if e = dec.Decode(&d.UnkContents); e != nil {
			return fmt.Errorf("dic initializer, UnkContents: %v", e)
		}
		return nil
	}(); err != nil {
		panic(err)
	}
	return
}
