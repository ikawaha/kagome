package dict

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_ConnectionTableAt(t *testing.T) {
	var ct ConnectionTable
	const (
		row = 4
		col = 5
	)
	ct.Row = row
	ct.Col = col
	ct.Vec = make([]int16, ct.Row*ct.Col)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			ct.Vec[i*col+j] = int16(i*col + j)
		}
	}
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			expected := int16(i*col + j)
			if r := ct.At(i, j); r != expected {
				t.Errorf("got %v, expected %v", r, expected)
			}
		}
	}

}

func Test_ConnectionTableWriteTo(t *testing.T) {
	ct := ConnectionTable{
		Row: 2,
		Col: 3,
		Vec: []int16{11, 12, 13, 21, 22, 23},
	}
	var b bytes.Buffer
	n, err := ct.WriteTo(&b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != int64(b.Len()) {
		t.Errorf("got %v, expected %v", n, b.Len())
	}
}

func Test_LoadConnectionTable(t *testing.T) {
	src := ConnectionTable{
		Row: 2,
		Col: 3,
		Vec: []int16{11, 12, 13, 21, 22, 23},
	}
	var b bytes.Buffer
	_, err := src.WriteTo(&b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	dst, err := ReadConnectionTable(&b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(src, dst) {
		t.Errorf("got %v, expected %v", dst, src)
	}
}

func Test_ContentInfoWriteToRead(t *testing.T) {
	c := ContentsMeta{
		POSStartIndex:      2,
		POSHierarchy:       22,
		BaseFormIndex:      3,
		ReadingIndex:       4,
		PronunciationIndex: 5,
		"hoo":              7,
	}
	var buf bytes.Buffer
	n, err := c.WriteTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	if want, got := int64(buf.Len()), n; want != got {
		t.Errorf("write length, want=%d, got=%d", want, got)
	}
	got, err := ReadContentsMeta(&buf)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	if !reflect.DeepEqual(c, got) {
		t.Errorf("want=%+v, got=%+v", c, got)
	}
}
