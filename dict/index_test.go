package dict

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_BuildIndexTable(t *testing.T) {
	sortedKeywords := []string{
		"aaa", //0
		"bbb", //1
		"bbb", //2
		"ccc", //3
		"ddd", //4
		"ddd", //5
	}
	idx, err := BuildIndexTable(sortedKeywords)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	checkList := []struct {
		id  int32
		ok  bool
		dup int32
	}{
		{0, false, 0},
		{1, true, 1},
		{2, false, 0},
		{3, false, 0},
	}
	for _, v := range checkList {
		x, ok := idx.Dup[v.id]
		if ok != v.ok {
			t.Errorf("got %v, expected %v", ok, v.ok)
		}
		if !v.ok {
			continue
		}
		if x != v.dup {
			t.Errorf("got %v, expected %v", x, v.dup)
		}
	}
}

func Test_CommonPrefixSearch(t *testing.T) {
	sortedKeywords := []string{
		"す",    //0
		"すし",   //1
		"すし",   //2
		"すし",   //3
		"すし",   //4
		"すしめし", //5
		"すしめし", //6
	}
	idx, err := BuildIndexTable(sortedKeywords)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	lens, outputs := idx.CommonPrefixSearch("すしめしたべた")
	expectedLens := []int{3, 6, 12} // byte length
	if !reflect.DeepEqual(lens, expectedLens) {
		t.Errorf("common prefix search lens, got %v, expected %v", lens, expectedLens)
	}
	expectdOutputs := [][]int{{0}, {1, 2, 3, 4}, {5, 6}}
	if !reflect.DeepEqual(outputs, expectdOutputs) {
		t.Errorf("common prefix search outputs, got %v, expected %v", outputs, expectdOutputs)
	}
}

func Test_CommonPrefixSearchCallback(t *testing.T) {
	sortedKeywords := []string{
		"す",    //0
		"すし",   //1
		"すし",   //2
		"すし",   //3
		"すし",   //4
		"すしめし", //5
		"すしめし", //6
	}
	expected := []struct {
		id, l int
	}{
		{0, 3},
		{1, 6},
		{2, 6},
		{3, 6},
		{4, 6},
		{5, 12},
		{6, 12},
	}
	idx, err := BuildIndexTable(sortedKeywords)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	var i int
	idx.CommonPrefixSearchCallback("すしめしたべた", func(id, l int) {
		if expected[i].id != id {
			t.Errorf("common prefix search callback id, got %v, expected %v", id, expected[i].id)
		}
		if expected[i].l != l {
			t.Errorf("common prefix search callback len, got %v, expected %v", l, expected[i].l)
		}
		i++
	})
}

func Test_Search(t *testing.T) {
	sortedKeywords := []string{
		"す",    //0
		"すし",   //1
		"すし",   //2
		"すし",   //3
		"すし",   //4
		"すしめし", //5
		"すしめし", //6
	}
	idx, err := BuildIndexTable(sortedKeywords)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	ts := []struct {
		word string
		ids  []int
	}{
		{"す", []int{0}},
		{"すし", []int{1, 2, 3, 4}},
		{"すしめし", []int{5, 6}},
	}
	for _, v := range ts {
		ids := idx.Search(v.word)
		if !reflect.DeepEqual(ids, v.ids) {
			t.Errorf("search ids, got %v, expected %v", ids, v.ids)
		}
	}
}

func Test_IndexTableReadAndWrite(t *testing.T) {
	sortedKeywords := []string{
		"す",    //0
		"すし",   //1
		"すし",   //2
		"すし",   //3
		"すし",   //4
		"すしめし", //5
		"すしめし", //6
	}
	org, err := BuildIndexTable(sortedKeywords)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	var b bytes.Buffer
	n, err := org.WriteTo(&b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != int64(b.Len()) {
		t.Errorf("write len: got %v, expected %v", n, b.Len())
	}

	cpy, err := ReadIndexTable(&b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(org, cpy) {
		t.Errorf("got %v, expected %v", cpy, org)
	}

}
