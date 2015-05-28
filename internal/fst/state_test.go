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

package fst

import (
	"fmt"
	"testing"
)

func TestStateEq01(t *testing.T) {
	type pair struct {
		x *state
		y *state
	}

	s := &state{}

	crs := []struct {
		call pair
		resp bool
	}{
		{pair{x: s, y: s}, true},
		{pair{x: nil, y: nil}, false},
		{pair{x: nil, y: &state{}}, false},
		{pair{x: &state{}, y: nil}, false},
		{pair{&state{ID: 1}, &state{ID: 2}}, true},
		{pair{&state{IsFinal: true}, &state{IsFinal: false}}, false},
		{pair{&state{Output: map[byte]int32{1: 555}}, &state{}}, false},
		{pair{&state{Output: map[byte]int32{1: 555}}, &state{Output: map[byte]int32{1: 555}}},
			true},
		{pair{&state{Output: map[byte]int32{1: 555}}, &state{Output: map[byte]int32{1: 444}}},
			false},
		{pair{&state{Output: map[byte]int32{1: 555}}, &state{Output: map[byte]int32{2: 555}}},
			false},
		{pair{&state{Tail: map[int32]bool{555: true}}, &state{Tail: map[int32]bool{555: true}}}, true},
	}
	for _, cr := range crs {
		if rst := cr.call.x.eq(cr.call.y); rst != cr.resp {
			t.Errorf("got %v, expected %v, %v\n", rst, cr.resp, cr)
		}
		if rst := cr.call.y.eq(cr.call.x); rst != cr.resp {
			t.Errorf("got %v, expected %v, %v\n", rst, cr.resp, cr)
		}
	}
}

func TestStateEq02(t *testing.T) {
	x := &state{ID: 1}
	y := &state{ID: 2}
	a := &state{
		Trans:  map[byte]*state{1: x, 2: y},
		Output: make(map[byte]int32),
	}
	b := &state{
		Trans:  map[byte]*state{1: x, 2: y},
		Output: make(map[byte]int32),
	}
	c := &state{
		Trans: map[byte]*state{1: y, 2: y},
	}
	d := &state{
		Trans: map[byte]*state{1: x, 2: y, 3: x},
	}
	if rst, exp := a.eq(b), true; rst != exp {
		t.Errorf("got %v, expected %v\n", rst, exp)
	}

	a.setOutput('a', 1)
	b.setOutput('a', 2)
	if rst, exp := a.eq(b), false; rst != exp {
		t.Errorf("got %v, expected %v\n", rst, exp)
	}

	if rst, exp := a.eq(c), false; rst != exp {
		t.Errorf("got %v, expected %v\n", rst, exp)
	}
	if rst, exp := a.eq(c), false; rst != exp {
		t.Errorf("got %v, expected %v\n", rst, exp)
	}
	if rst, exp := a.eq(d), false; rst != exp {
		t.Errorf("got %v, expected %v\n", rst, exp)
	}

}

func TestStateString01(t *testing.T) {
	crs := []struct {
		call *state
		resp string
	}{
		{nil, "<nil>"},
	}
	for _, cr := range crs {
		if rst := cr.call.String(); rst != cr.resp {
			t.Errorf("got %v, expected %v, %v\n", rst, cr.resp, cr)
		}
	}
	r := &state{}
	s := state{
		ID:      1,
		Trans:   map[byte]*state{1: nil, 2: r},
		Output:  map[byte]int32{3: 555, 4: 888},
		Tail:    int32Set{1111: true},
		IsFinal: true,
	}
	fmt.Println(s.String())
}
