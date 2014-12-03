//  Copyright (c) 2014 ikawaha.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package kagome

import (
	"reflect"
	"testing"
)

func TestNewNodePool01(t *testing.T) {
	np := newNodePool(0)
	expected := nodePool{}
	if reflect.DeepEqual(np, expected) {
		t.Errorf("got %v, expected %v\n", np, expected)
	}
}

func TestNewNodePool02(t *testing.T) {
	size := 10
	np := newNodePool(size)
	if np.usage != 0 {
		t.Errorf("usage: got %v, expected %v\n", np.usage, 0)
	}
	if cap(np.buf) != size {
		t.Errorf("buf capacity: got %v, expected %v\n", cap(np.buf), size)
	}
}

func TestNewNodePoolGetAndClear01(t *testing.T) {
	np := newNodePool(0)
	n := np.get()
	if n == nil {
		t.Error("node alloc error\n")
	}
	if np.usage != 1 {
		t.Errorf("usage: got %v, expected 1\n", np.usage)
	}
	np.clear()
	if np.usage != 0 {
		t.Errorf("usage: got %v, expected 0\n", np.usage)
	}
}

func TestNewNodePoolGetAndClear02(t *testing.T) {
	np := newNodePool(1)
	np.get()
	np.get()
	n := np.get()
	if n == nil {
		t.Error("node alloc error\n")
	}
	if np.usage != 3 {
		t.Errorf("usage: got %v, expected 3\n", np.usage)
	}
	np.clear()
	if np.usage != 0 {
		t.Errorf("usage: got %v, expected 0\n", np.usage)
	}
}

func TestNilNodePoolGetAndClear01(t *testing.T) {
	var np *nodePool
	if np != nil {
		t.Fatalf("initialize error, expected nil, got %v\n", np)
	}
	np.get()
	np.get()
	n := np.get()
	if n == nil {
		t.Error("node alloc error\n")
	}
	np.clear()
}
