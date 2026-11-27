package lattice

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

// nodeBlockSize is the number of nodes in one block of a nodeArena.
const nodeBlockSize = 4096

// nodeArena hands out nodes from fixed-size blocks owned by one lattice.
// Blocks are never reallocated, so the pointers it returns stay valid until
// reset, which makes every node available again at once.
type nodeArena struct {
	blocks [][]Node
	block  int // index of the block the next node comes from
	next   int // index of the next node in that block
}

func (a *nodeArena) alloc() *Node {
	if a.block == len(a.blocks) {
		a.blocks = append(a.blocks, make([]Node, nodeBlockSize))
	}
	n := &a.blocks[a.block][a.next]
	a.next++
	if a.next == nodeBlockSize {
		a.block++
		a.next = 0
	}
	return n
}

func (a *nodeArena) reset() {
	a.block = 0
	a.next = 0
}
