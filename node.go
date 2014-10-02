package kagome

import ()

const (
	BOSEOS_ID int       = -1
	DUMMY     NodeClass = iota
	KNOWN
	UNKNOWN
	USER
)

type NodeClass int

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
	np = new(nodePool)
	np.buf = make([]*node, size)
	for i := range np.buf {
		np.buf[i] = new(node)
	}
	return
}

const _MINIMUM_NODE_POOL_SIZE = 128

func (np *nodePool) get() (n *node) {
	if np.usage == len(np.buf) {
		var neoCap int
		if np.usage != 0 {
			neoCap = np.usage * 2
		} else {
			neoCap = _MINIMUM_NODE_POOL_SIZE
		}
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
	np.usage = 0
}
