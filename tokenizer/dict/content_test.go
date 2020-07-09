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

func TestContentsSave(t *testing.T) {
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

func TestNewContents(t *testing.T) {
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
