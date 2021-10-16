package lattice

import (
	"reflect"
	"testing"
)

func Test_NodeClassString(t *testing.T) {
	testdata := []struct {
		class NodeClass
		want  string
	}{
		{class: DUMMY, want: "DUMMY"},
		{class: KNOWN, want: "KNOWN"},
		{class: UNKNOWN, want: "UNKNOWN"},
		{class: USER, want: "USER"},
		{class: NodeClass(999), want: "UNDEF"},
	}
	for _, p := range testdata {
		if got := p.class.String(); got != p.want {
			t.Errorf("got %v, expected %v", got, p.want)
		}
	}
}

func TestNodeHeap_PushPop(t *testing.T) {
	idSorter := func(x, y *Node) bool {
		return x.ID < y.ID
	}
	heap := NodeHeap{
		less: idSorter,
	}
	testdata := []struct {
		name string
		ids  []int
		want []int
	}{
		{
			name: "ascending order",
			ids:  []int{1, 2, 3, 4, 5, 6, 7},
			want: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "descending order",
			ids:  []int{7, 6, 5, 4, 3, 2, 1},
			want: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "random order",
			ids:  []int{3, 6, 4, 1, 7, 5, 2},
			want: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "list /w duplicate items",
			ids:  []int{3, 6, 3, 4, 1, 3, 6, 2, 7, 5, 2},
			want: []int{1, 2, 2, 3, 3, 3, 4, 5, 6, 6, 7},
		},
	}
	for _, data := range testdata {
		for _, v := range data.ids {
			heap.Push(&Node{ID: v})
		}
		got := make([]int, 0, heap.Size())
		for !heap.Empty() {
			n := heap.Pop()
			if n == nil {
				t.Fatalf("unexpected nil node, heap=%+v", heap)
			}
			got = append(got, n.ID)
		}
		if !reflect.DeepEqual(got, data.want) {
			t.Errorf("got %+v, want %+v", got, data.want)
		}
	}
}
