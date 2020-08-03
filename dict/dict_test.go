package dict

import (
	"archive/zip"
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func newTestDict(t *testing.T) *Dict {
	idx, err := BuildIndexTable([]string{"key1", "key2", "key3"})
	if err != nil {
		t.Fatalf("build index table failed, %v", err)
	}

	return &Dict{
		Morphs: Morphs{
			{LeftID: 111, RightID: 222, Weight: 333},
			{LeftID: 444, RightID: 555, Weight: 666},
		},
		POSTable: POSTable{
			POSs: []POS{
				{1, 2, 3, 3, 4, 5},
				{1, 2, 3, 6, 7, 8},
			},
			NameList: []string{
				"str1", "str2", "str3", "str4", "str5", "str6", "str7", "str8",
			},
		},
		ContentsMeta: ContentsMeta{
			"meta": 47,
		},
		Contents: Contents{
			[]string{"a1", "a2", "a3"},
			[]string{"b1", "b2", "b3"},
		},
		Connection: ConnectionTable{
			Row: 2,
			Col: 3,
			Vec: []int16{0, 1, 2, 3, 4, 5},
		},
		Index: idx,
		CharClass: []string{
			"class1", "class2", "class3",
		},
		CharCategory: []byte{'a', 'b', 'c'},
		InvokeList:   []bool{true, false, true},
		GroupList:    []bool{false, true, false},
		UnkDict: UnkDict{
			Morphs: Morphs{
				{LeftID: 11, RightID: 22, Weight: 33},
				{LeftID: 44, RightID: 55, Weight: 66},
			},
			Index: map[int32]int32{
				1: 1111,
				2: 2222,
			},
			IndexDup: map[int32]int32{
				1: 0,
				2: 3,
			},
			ContentsMeta: ContentsMeta{
				"unkmeta": 47,
			},
			Contents: Contents{
				[]string{"aa1", "aa2", "aa3"},
				[]string{"bb1", "bb2", "bb3"},
			},
		},
	}
}

// save <--> load
func Test_DictSaveLoad(t *testing.T) {
	dict := newTestDict(t)

	var b bytes.Buffer
	zw := zip.NewWriter(&b)
	if err := dict.Save(zw); err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	r := bytes.NewReader(b.Bytes())
	zr, err := zip.NewReader(r, r.Size())
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	got, err := Load(zr, true)
	if err != nil {
		t.Fatalf("zipped dict loading failed, %v", err)
	}

	if !reflect.DeepEqual(dict, got) {
		t.Errorf("want %+v, got %+v", dict, got)
		fmt.Printf("%T\n", got.ContentsMeta)
		fmt.Printf("%T\n", dict.ContentsMeta)
	}
}
