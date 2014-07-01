package tokenizer

import (
	"fmt"

	"github.com/ikawaha/kagome/dic"
)

const (
	BOSEOS = -1
)

type NodeClass byte

const (
	KNOWN NodeClass = iota
	UNKNOWN
	INSERTED
	USER
	DUMMY
)

type Node struct {
	id      int
	start   int
	class   NodeClass
	cost    int32
	left    int32
	right   int32
	weight  int32
	surface []byte
	prev    *Node
}

func (this Node) String() string {
	if this.id == BOSEOS {
		ret := fmt.Sprintf("[BOSEOS], start:%v, cost:%v", this.start, this.cost)
		if this.prev != nil {
			ret += fmt.Sprintf(" --> %v", this.prev.id)
		}
		return ret
	}
	content := ""
	if this.id >= 0 && this.id < len(dic.Contents) {
		content = fmt.Sprintf("%v", dic.Contents[this.id])
	}
	ret := fmt.Sprintf("id:%v, start:%v, cost:%v, class:%v, left:%v, right:%v, weight:%v, surface:%v, %v\n",
		this.id, this.start, this.cost, this.class, this.left, this.right, this.weight, string(this.surface), content)
	if this.prev != nil {
		ret += fmt.Sprintf(" --> %v", this.prev.id)
	}
	return ret
}

func (this NodeClass) String() string {
	switch this {
	case KNOWN:
		return "KNOWN"
	case UNKNOWN:
		return "UNKNOWN"
	case INSERTED:
		return "INSERTED"
	case USER:
		return "USER"
	case DUMMY:
		return "DUMMY"
	}
	return "<UNKNOWN_NODE_CLASS>"
}
