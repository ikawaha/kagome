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
	"bytes"
	"sync"

	"github.com/ikawaha/kagome/internal/dic/data"
)

const (
	IPADicPath = "dic/ipa"
)

var (
	sysDicIPA     *Dic
	initSysDicIPA sync.Once
)

// SysDic returns the kagome system dictionary.
func SysDic() *Dic {
	return SysDicIPA()
}

// SysDicIPA returns the IPA system dictionary.
func SysDicIPA() *Dic {
	initSysDicIPA.Do(func() {
		sysDicIPA = loadInternalSysDic(IPADicPath)
	})
	return sysDicIPA
}

func loadInternalSysDic(path string) (d *Dic) {
	d = new(Dic)
	var (
		buf []byte
		err error
	)
	// morph.dic
	if buf, err = data.Asset(path + "/morph.dic"); err != nil {
		panic(err)
	}
	if err = d.loadMorphDicPart(bytes.NewBuffer(buf)); err != nil {
		panic(err)
	}
	// index.dic
	if buf, err = data.Asset(path + "/index.dic"); err != nil {
		panic(err)
	}
	if err = d.loadIndexDicPart(bytes.NewBuffer(buf)); err != nil {
		panic(err)
	}
	// connection.dic
	if buf, err = data.Asset(path + "/connection.dic"); err != nil {
		panic(err)
	}
	if err = d.loadConnectionDicPart(bytes.NewBuffer(buf)); err != nil {
		panic(err)
	}
	// chardef.dic
	if buf, err = data.Asset(path + "/chardef.dic"); err != nil {
		panic(err)
	}
	if err = d.loadCharDefDicPart(bytes.NewBuffer(buf)); err != nil {
		panic(err)
	}
	// unk.dic
	if buf, err = data.Asset(path + "/unk.dic"); err != nil {
		panic(err)
	}
	if err = d.loadUnkDicPart(bytes.NewBuffer(buf)); err != nil {
		panic(err)
	}
	return
}

/*
func loadInternalSysDic(path string) (d *Dic) {
	d = new(Dic)
	if err := func() error {
		buf, e := data.Asset(path + "/morph.dic")
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
		buf, e := data.Asset(path + "/index.dic")
		if e != nil {
			return e
		}
		t, e := fst.Read(bytes.NewReader(buf))
		if e != nil {
			return fmt.Errorf("dic initializer, Index: %v", e)
		}
		d.Index = t
		return nil
	}(); err != nil {
		panic(err)
	}

	if err := func() error {
		buf, e := data.Asset(path + "/connection.dic")
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
		buf, e := data.Asset(path + "/chardef.dic")
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
		buf, e := data.Asset(path + "/unk.dic")
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
*/
