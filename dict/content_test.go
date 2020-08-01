package dict

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_ContentsSave(t *testing.T) {
	m := [][]string{
		{"11", "12", "13"},
		{"21", "22", "23"},
		{"31", "32", "33"},
	}
	var b bytes.Buffer
	n, err := Contents(m).WriteTo(&b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != int64(b.Len()) {
		t.Errorf("got %v, expected %v", n, b.Len())
	}
}

func Test_NewContents(t *testing.T) {
	src := [][]string{
		{"11", "12", "13"},
		{"21", "22", "23"},
		{"31", "32", "33"},
	}
	var b bytes.Buffer
	_, err := Contents(src).WriteTo(&b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	dst := NewContents(b.Bytes())
	if !reflect.DeepEqual(src, dst) {
		t.Errorf("got %v, expected %v", dst, src)
	}
}
