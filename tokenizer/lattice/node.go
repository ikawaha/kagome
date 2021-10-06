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
