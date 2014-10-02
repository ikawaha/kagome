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
