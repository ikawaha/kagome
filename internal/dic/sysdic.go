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
	"bytes"
	"sync"

	"github.com/ikawaha/kagome/internal/dic/data"
)

const (
	// IPADicPath represents the internal IPA dictionary path.
	IPADicPath = "dic/ipa"
	// UniDicPath represents the internal UniDic dictionary path.
	UniDicPath = "dic/uni"
)

var (
	sysDicIPA     *Dic
	initSysDicIPA sync.Once
	sysDicUni     *Dic
	initSysDicUni sync.Once
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

// SysDicUni returns the UniDic system dictionary.
func SysDicUni() *Dic {
	initSysDicUni.Do(func() {
		sysDicUni = loadInternalSysDic(UniDicPath)
	})
	return sysDicUni
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
	// content.dic
	if buf, err = data.Asset(path + "/content.dic"); err != nil {
		panic(err)
	}
	d.Contents = NewContents(buf)
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
