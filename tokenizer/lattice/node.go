package lattice

import "sync"

// BosEosID represents Reserved identifier of node id.
const BosEosID int = -1

// NodeClass codes.
const (
	DUMMY NodeClass = iota
	KNOWN
	UNKNOWN
	USER
)

// NodeClass represents a node type.
type NodeClass int

// String returns a string representation of a node class.
func (nc NodeClass) String() string {
	switch nc {
	case DUMMY:
		return "DUMMY"
	case KNOWN:
		return "KNOWN"
	case UNKNOWN:
		return "UNKNOWN"
	case USER:
		return "USER"
	}
	return "UNDEF"
}

// Node is a lattice node.
type Node struct {
	ID       int
	Position int // byte position
	Start    int // rune position
	Class    NodeClass
	Cost     int32
	Left     int32
	Right    int32
	Weight   int32
	Surface  string
	prev     *Node
}

var nodePool = sync.Pool{
	New: func() interface{} {
		return new(Node)
	},
}

type NodeHeap struct {
	list []*Node
	less func(x, y *Node) bool
}

// Push adds a node to the heap.
func (h *NodeHeap) Push(n *Node) {
	i := len(h.list)
	h.list = append(h.list, n)
	for i != 0 {
		p := (i - 1) / 2
		if !h.less(h.list[p], h.list[i]) {
			h.list[p], h.list[i] = h.list[i], h.list[p]
		}
		i = p
	}
}

// Pop returns the highest priority node of the heap. If the heap is empty, Pop returns nil.
func (h *NodeHeap) Pop() *Node {
	if len(h.list) < 1 {
		return nil
	}
	ret := h.list[0]
	if len(h.list) > 1 {
		h.list[0] = h.list[len(h.list)-1]
	}
	h.list[len(h.list)-1] = nil
	h.list = h.list[:len(h.list)-1]

	for i := 0; ; {
		min := i
		if left := (i+1)*2 - 1; left < len(h.list) && !h.less(h.list[min], h.list[left]) {
			min = left
		}
		if right := (i + 1) * 2; right < len(h.list) && !h.less(h.list[min], h.list[right]) {
			min = right
		}
		if min == i {
			break
		}
		h.list[i], h.list[min] = h.list[min], h.list[i]
		i = min
	}

	return ret
}

// Empty returns true if the heap is empty.
func (h NodeHeap) Empty() bool {
	return len(h.list) == 0
}

// Size returns the size of the heap.
func (h NodeHeap) Size() int {
	return len(h.list)
}
