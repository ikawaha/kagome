package dict

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_WriteMap(t *testing.T) {
	var b0, b1 bytes.Buffer

	m := map[int32]int32{
		1: 1,
		2: 4,
		3: 9,
	}

	sz0, err := writeMapInt32Int32(&b0, m)
	if err != nil {
		t.Errorf("unexpected error, %v", err)
	}
	sz1, err := writeMapInt32Int32(&b1, m)
	if err != nil {
		t.Errorf("unexpected error, %v", err)
	}
	if sz0 != sz1 {
		t.Errorf("different size, %d <> %d", sz0, sz1)
	}
	if !reflect.DeepEqual(b0.Bytes(), b1.Bytes()) {
		t.Errorf("different silialization, %v <> %v", b0, b1)
	}
}

func Test_UnkDic_WriteAndRead(t *testing.T) {
	d := UnkDict{
		Morphs: Morphs{
			{LeftID: 1, RightID: 2, Weight: 3},
			{LeftID: 11, RightID: 22, Weight: 33},
		},
		Index: map[int32]int32{
			1: 1,
			2: 4,
			3: 9,
		},
		IndexDup: map[int32]int32{
			1: 7,
			2: 8,
			3: 9,
		},
		ContentsMeta: ContentsMeta{
			"meta": 47,
		},
		Contents: Contents{
			{"hello", "goodbye"},
			{"こんにちは", "さようなら"},
		},
	}

	// write
	var b bytes.Buffer
	sz, err := d.WriteTo(&b)
	if err != nil {
		t.Fatalf("unexpected write error, %v", err)
	}
	if expected := int64(150); sz != expected {
		t.Fatalf("silialization size, got %v, expected %v", sz, expected)
	}

	// read
	unk, err := ReadUnkDic(&b)
	if err != nil {
		t.Errorf("unexpected read error, %v", err)
	}
	if !reflect.DeepEqual(d, unk) {
		t.Errorf("got %v, expected %v", unk, d)
	}
}
