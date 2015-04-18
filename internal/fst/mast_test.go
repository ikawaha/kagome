//  Copyright (c) 2015 ikawaha.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package fst

import (
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestMASTBuildMAST01(t *testing.T) {
	inp := PairSlice{}
	m := buildMAST(inp)
	if m.initialState.ID != 0 {
		t.Errorf("got initial state id %v, expected 0\n", m.initialState.ID)
	}
	if len(m.states) != 1 {
		t.Errorf("expected: initial state only, got %v\n", m.states)
	}
	if len(m.finalStates) != 0 {
		t.Errorf("expected: final state is empty, got %v\n", m.finalStates)
	}
}

func TestMASTAccept01(t *testing.T) {
	inp := PairSlice{
		{"hello", 111},
		{"hello", 222},
		{"111", 111},
		{"112", 112},
		{"112", 122},
		{"211", 345},
	}
	m := buildMAST(inp)
	for _, pair := range inp {
		if ok := m.accept(pair.In); !ok {
			t.Errorf("expected: accept [%v]\n", pair.In)
		}
	}
	if ok := m.accept("aloha"); ok {
		t.Errorf("expected: reject \"aloha\"\n")
	}
}

func TestMASTRun01(t *testing.T) {
	inp := PairSlice{
		{"hello", 1111},
		{"hell", 2222},
		{"111", 111},
		{"112", 112},
		{"113", 122},
		{"211", 111},
	}
	m := buildMAST(inp)
	for _, pair := range inp {
		out, ok := m.run(pair.In)
		if !ok {
			t.Errorf("expected: accept [%v]\n", pair.In)
		}
		if len(out) != 1 {
			t.Errorf("input: %v, output size: got %v, expected 1\n", pair.In, len(out))
		}
		if out[0] != pair.Out {
			t.Errorf("input: %v, output: got %v, expected %v\n", pair.In, pair.Out, out[0])
		}
	}
	if out, ok := m.run("aloha"); ok {
		t.Errorf("expected: reject \"aloha\", %v\n", out)
	}
}

func TestMASTRun02(t *testing.T) {
	inp := PairSlice{
		{"hello", 1111},
		{"hello", 2222},
	}
	m := buildMAST(inp)
	for _, pair := range inp {
		out, ok := m.run(pair.In)
		if !ok {
			t.Errorf("expected: accept [%v]\n", pair.In)
		}
		if len(out) != 2 {
			t.Errorf("input: %v, output size: got %v, expected 2\n", pair.In, len(out))
		}
		expected := []int32{1111, 2222}
		sort.Sort(int32Slice(out))
		sort.Sort(int32Slice(expected))
		if !reflect.DeepEqual(out, expected) {
			t.Errorf("input: %v, output: got %v, expected %v\n", pair.In, out, expected)
		}
	}
	if out, ok := m.run("aloha"); ok {
		t.Errorf("expected: reject \"aloha\", %v\n", out)
	}
}

func TestMASTDot01(t *testing.T) {
	inp := PairSlice{
		{"apr", 30},
		{"aug", 31},
		{"dec", 31},
		{"feb", 28},
		{"feb", 29},
	}
	m := buildMAST(inp)
	m.dot(os.Stdout)
}

func TestMASTDot02(t *testing.T) {
	inp := PairSlice{
		{"apr", 30},
		{"aug", 31},
		{"dec", 31},
		{"feb", 28},
		{"feb", 29},
		{"lucene", 1},
		{"lucid", 2},
		{"lucifer", 666},
	}
	m := buildMAST(inp)
	m.dot(os.Stdout)
}
