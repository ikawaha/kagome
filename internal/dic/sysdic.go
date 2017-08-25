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
	"bytes"
	"fmt"
	"sync"

	"github.com/ikawaha/kagome/internal/dic/data"
)

const (
	// IPADicPath represents the internal IPA dictionary path.
	IPADicPath = "dic/ipa/ipa.dic"
	// UniDicPath represents the internal UniDic dictionary path.
	UniDicPath = "dic/uni/uni.dic"
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

// SysDicSimple returns the kagome system dictionary without contents.
func SysDicSimple() *Dic {
	return SysDicIPASimple()
}

// SysDicIPA returns the IPA system dictionary.
func SysDicIPA() *Dic {
	initSysDicIPA.Do(func() {
		sysDicIPA = loadInternalSysDicFull(IPADicPath)
	})
	return sysDicIPA
}

// SysDicUni returns the UniDic system dictionary.
func SysDicUni() *Dic {
	initSysDicUni.Do(func() {
		sysDicUni = loadInternalSysDicFull(UniDicPath)
	})
	return sysDicUni
}

// SysDicIPASimple returns the IPA system dictionary without contents.
func SysDicIPASimple() *Dic {
	initSysDicIPA.Do(func() {
		sysDicIPA = loadInternalSysDicSimple(IPADicPath)
	})
	return sysDicIPA
}

// SysDicUniSimple returns the IPA system dictionary without contents.
func SysDicUniSimple() *Dic {
	initSysDicUni.Do(func() {
		sysDicUni = loadInternalSysDicSimple(UniDicPath)
	})
	return sysDicUni
}

func loadInternalSysDicFull(path string) (d *Dic) {
	return loadInternalSysDic(path, true)
}

func loadInternalSysDicSimple(path string) (d *Dic) {
	return loadInternalSysDic(path, false)
}

func loadInternalSysDic(path string, full bool) (d *Dic) {
	buf := make([]byte, 0, 36*1024*1024) // 36MB
	defer func() { buf = nil }()

	for i := 0; ; i++ {
		b, err := data.Asset(path + fmt.Sprintf(".%03x", i))
		if err != nil {
			break
		}
		buf = append(buf, b...)
	}
	r := bytes.NewReader(buf)
	zr, err := zip.NewReader(r, r.Size())
	if err != nil {
		panic(err)
	}
	d, err = load(zr, full)
	if err != nil {
		panic(err)
	}
	return d
}
