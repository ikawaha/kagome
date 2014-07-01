package tokenizer

type NodePool struct {
	usage int
	buf   []*Node
}

func NewNodePool(a_size int) *NodePool {
	ret := new(NodePool)
	ret.buf = make([]*Node, a_size)
	for i, size := 0, len(ret.buf); i < size; i++ {
		ret.buf[i] = new(Node)
	}
	return ret
}

func (this *NodePool) get() *Node {
	var p *Node
	if this.usage >= len(this.buf) {
		dst := make([]*Node, len(this.buf)*2)
		copy(dst, this.buf)
		for i, size := len(this.buf), len(dst); i < size; i++ {
			dst[i] = new(Node)
		}
		this.buf = dst
	}
	p = this.buf[this.usage]
	this.usage++
	return p
}

func (this *NodePool) clear() {
	this.usage = 0
}
