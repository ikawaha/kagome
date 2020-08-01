package dict

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_POSTableAdd(t *testing.T) {
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
			t.Errorf("%d, input %v, got %+v, expected %+v", i, d.In, pos, d.Exp)
		}
	}
}

func Test_POSString(t *testing.T) {
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
	for i, want := range data {
		pos := m.Add(want)
		got := make([]string, 0, len(pos))
		for _, id := range pos {
			got = append(got, table.NameList[id])
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%d, input %v, got %+v, expected %+v", i, want, got, want)
		}
	}
}

func Test_POSTableReadAndWrite(t *testing.T) {
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
		t.Errorf("unexpected error: %v", err)
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
