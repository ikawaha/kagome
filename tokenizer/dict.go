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

package tokenizer

import "github.com/ikawaha/kagome/tokenizer/dict"

// Dict represents a dictionary.
type Dic struct {
	dic *dict.Dict
}

// SysDic returns the system dictionary (IPA dictionary).
func SysDic() Dic {
	return Dic{dict.SysDic()}
}

// SysDicSimple returns the system dictionary (IPA dictionary w/o contents).
func SysDicSimple() Dic {
	return Dic{dict.SysDicSimple()}
}

// SysDicIPA returns the IPA dictionary as the system dictionary.
func SysDicIPA() Dic {
	return Dic{dict.SysDicIPA()}
}

// SysDicIPASimple returns the simple IPA dictionary as the system dictionary (w/o contents).
func SysDicIPASimple() Dic {
	return Dic{dict.SysDicIPASimple()}
}

// SysDicUni returns the UniDic dictionary as the system dictionary.
func SysDicUni() Dic {
	return Dic{dict.SysDicUni()}
}

// SysDicUniSimple returns the simple UniDic dictionary as the system dictionary (w/o contents).
func SysDicUniSimple() Dic {
	return Dic{dict.SysDicUniSimple()}
}

// NewDic loads a dictionary from a file.
func NewDic(path string) (Dic, error) {
	d, err := dict.Load(path)
	return Dic{d}, err
}

// NewDicSimple loads a dictionary from a file w/o contents.
func NewDicSimple(path string) (Dic, error) {
	d, err := dict.LoadSimple(path)
	return Dic{d}, err
}
