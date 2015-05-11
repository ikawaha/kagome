package dic

import (
	"bytes"
	"reflect"
	"testing"
)

func TestMorphsSave(t *testing.T) {
	m := []Morph{
		{1, 1, 1},
		{2, 2, 2},
		{3, 3, 3},
	}
	var b bytes.Buffer
	n, err := MorphSlice(m).WriteTo(&b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != int64(b.Len()) {
		t.Errorf("got %v, expected %v", n, b.Len())
	}
}

func TestLoadMorphSlice(t *testing.T) {
	src := []Morph{
		{1, 1, 1},
		{2, 2, 2},
		{3, 3, 3},
	}
	var b bytes.Buffer
	_, err := MorphSlice(src).WriteTo(&b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	dst, err := LoadMorphSlice(&b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(src, dst) {
		t.Errorf("got %v, expected %v", dst, src)
	}
}
