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
