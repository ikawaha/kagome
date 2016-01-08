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

package dic

import (
	"bytes"
	"reflect"
	"testing"
)

func TestConnectionTableAt(t *testing.T) {
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
				t.Errorf("got %v, expected %v\n", r, expected)
			}
		}
	}

}

func TestConnectionTableWriteTo(t *testing.T) {
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

func TestLoadConnectionTable(t *testing.T) {
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

	dst, err := LoadConnectionTable(&b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(src, dst) {
		t.Errorf("got %v, expected %v", dst, src)
	}
}
