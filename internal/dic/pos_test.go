// Copyright 2017 ikawaha
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
	"reflect"
	"testing"
)

func TestPOSTableAdd(t *testing.T) {
	data := []struct {
		In  []string
		Exp POS
	}{
		{In: []string{"動詞", "自立", "*", "*", "五段・マ行", "基本形"}, Exp: POS{1, 2, 3, 3, 4, 5}},
		{In: []string{"動詞", "接尾", "*", "*", "五段・サ行", "未然形"}, Exp: POS{1, 6, 3, 3, 7, 8}},
		{In: []string{"一般", "*", "*", "*", "*"}, Exp: POS{9, 3, 3, 3, 3}},
		{In: []string{"動詞", "自立", "*", "*", "五段・マ行", "未然形"}, Exp: POS{1, 2, 3, 3, 4, 8}},
	}
	m := POSMap{}
	for i, d := range data {
		pos := m.Add(d.In)
		if !reflect.DeepEqual(pos, d.Exp) {
			t.Errorf("%d, input %v, got %+v, expected %+v\n", i, d.In, pos, d.Exp)
		}
	}
}

func TestPOSString(t *testing.T) {
	data := [][]string{
		{"動詞", "接尾", "*", "*"},
		{"動詞", "接尾", "*", "*", "五段・サ行,未然形"},
		{"自立", "*", "*", "五段・マ行,基本形"},
		{"動詞", "自立", "*", "*", "五段・マ行,未然形"},
	}
	m := POSMap{}
	for _, d := range data {
		m.Add(d)
	}
	table := POSTable{
		NameList: m.List(),
	}
	for i, d := range data {
		pos := m.Add(d)
		p := table.GetPOSName(pos)
		if !reflect.DeepEqual(p, d) {
			t.Errorf("%d, input %v, got %+v, expected %+v\n", i, d, p, d)
		}
	}
}

func TestPOSTableReadAndWrite(t *testing.T) {
	data := [][]string{
		{"動詞", "接尾", "*", "*"},
		{"動詞", "接尾", "*", "*", "五段・サ行,未然形"},
		{"自立", "*", "*", "五段・マ行,基本形"},
		{"動詞", "自立", "*", "*", "五段・マ行,未然形"},
	}
	org := POSTable{
		POSs: []POS{},
	}
	m := POSMap{}
	for _, d := range data {
		org.POSs = append(org.POSs, m.Add(d))
	}
	org.NameList = m.List()

	var b bytes.Buffer
	n, err := org.WriteTo(&b)
	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
	}
	if n != int64(b.Len()) {
		t.Errorf("write len: got %v, expected %v", n, b.Len())
	}

	cpy, err := ReadPOSTable(&b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(org, cpy) {
		t.Errorf("got %v, expected %v", cpy, org)
	}

}
