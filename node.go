package kagome

// Reserved identifier of node id.
const BosEosId int = -1

// NodeClass codes.
const (
	DUMMY NodeClass = iota
	KNOWN
	UNKNOWN
	USER
)

// NodeClass represents a node type.
type NodeClass int

// String returns a string representation of a node class.1
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

type node struct {
	id      int
	start   int
	class   NodeClass
	cost    int32
	left    int32
	right   int32
	weight  int32
	surface string
	prev    *node
}

type nodePool struct {
	usage int
	buf   []*node
}

func newNodePool(size int) (np *nodePool) {
	const minimumNodePoolCapacity = 128
	if size <= 0 {
		size = minimumNodePoolCapacity
	}
	np = new(nodePool)
	np.buf = make([]*node, size)
	for i := range np.buf {
		np.buf[i] = new(node)
	}
	return
}

func (np *nodePool) get() (n *node) {
	if np == nil {
		return new(node)
	}
	if np.usage == len(np.buf) {
		neoCap := np.usage * 2
		dst := make([]*node, neoCap)
		copy(dst, np.buf)
		for i, end := len(np.buf), neoCap; i < end; i++ {
			dst[i] = new(node)
		}
		np.buf = dst
	}
	n = np.buf[np.usage]
	np.usage++
	return
}

func (np *nodePool) clear() {
	if np == nil {
		return
	}
	np.usage = 0
}
