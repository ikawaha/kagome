package dict

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_MorphsSave(t *testing.T) {
	m := []Morph{
		{1, 1, 1},
		{2, 2, 2},
		{3, 3, 3},
	}
	var b bytes.Buffer
	n, err := Morphs(m).WriteTo(&b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != int64(b.Len()) {
		t.Errorf("got %v, expected %v", n, b.Len())
	}
}

func Test_LoadMorphSlice(t *testing.T) {
	src := Morphs{
		{1, 1, 1},
		{2, 2, 2},
		{3, 3, 3},
	}
	var b bytes.Buffer
	_, err := Morphs(src).WriteTo(&b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	dst, err := ReadMorphs(&b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(src, dst) {
		t.Errorf("got %v, expected %v", dst, src)
	}
}
