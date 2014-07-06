package tokenizer

import (
	"fmt"
	"unicode/utf8"

	"github.com/ikawaha/kagome/dic"
)

const (
	_MAX_INT32               = 1<<31 - 1
	_MAX_UNKNOWN_WORD_LENGTH = 1024
	_INIT_NODE_BUFFER_SIZE   = 512
)

type lattice struct {
	input  []byte
	list   [][]*Node
	output []*Node
	pool   *NodePool
	udic   *dic.UserDic
}

func NewLattice() *lattice {
	ret := new(lattice)
	ret.pool = NewNodePool(_INIT_NODE_BUFFER_SIZE)
	return ret
}

func (this *lattice) setUserDic(a_userdic *dic.UserDic) {
	this.udic = a_userdic
}

func (this *lattice) addNode(a_id, a_pos, a_start int, a_surface []byte, a_class NodeClass) {
	var cost dic.Cost
	switch a_class {
	case DUMMY:
		// use default cost
	case KNOWN:
		cost = dic.Costs[a_id]
	case UNKNOWN:
		cost = dic.UnkCosts[a_id]
	case USER:
		// use default cost
	}
	node := this.pool.get()
	node.id = a_id
	node.start = a_start
	node.class = a_class
	node.left, node.right, node.weight = int32(cost.Left), int32(cost.Right), int32(cost.Weight)
	node.surface = a_surface
	node.prev = nil
	p := a_pos + utf8.RuneCount(node.surface)
	this.list[p] = append(this.list[p], node)
}

func (this *lattice) build(a_input *string) (err error) {
	this.pool.clear()

	this.input = []byte(*a_input)

	runeCount := utf8.RuneCount(this.input)
	this.list = make([][]*Node, runeCount+2)

	this.addNode(BOSEOS, 0, 0, this.input[0:0], DUMMY)
	this.addNode(BOSEOS, runeCount+1, runeCount, this.input[runeCount:runeCount], DUMMY)

	chPos := -1
	for bufPos, ch := range *a_input {
		chPos++

		// (1) TODO: USER DIC
		anyMatches := false
		if this.udic != nil {
			prefixs, ids := this.udic.Index.CommonPrefixSearchBytes(this.input[bufPos:])
			anyMatches = len(prefixs) > 0
			for key, substr := range prefixs {
				id := ids[key]
				this.addNode(id, chPos, chPos, this.input[bufPos:bufPos+len(substr)], USER)
			}
		}
		if anyMatches {
			continue
		}
		// (2) KNOWN DIC
		prefixs, ids := dic.Index.CommonPrefixSearchBytes(this.input[bufPos:])
		anyMatches = len(prefixs) > 0
		for key, substr := range prefixs {
			id := ids[key]
			c, ok := dic.Counts[id]
			if !ok {
				c = 1
			}
			for x := 0; x < c; x++ {
				this.addNode(id+x, chPos, chPos, this.input[bufPos:bufPos+len(substr)], KNOWN)
			}
		}
		// (3) UNKNOWN DIC
		if !anyMatches || dic.InvokeList[dic.CharacterCategoryList[ch]] {
			class := dic.CharacterCategoryList[ch]
			endPos := bufPos + utf8.RuneLen(ch)
			unkWordLen := 1
			for i, w, size := endPos, 1, len(this.input); i < size; i += w {
				var c rune
				c, w = utf8.DecodeRune(this.input[i:])
				if dic.CharacterCategoryList[c] != class {
					break
				}
				endPos += w
				unkWordLen++
				if unkWordLen >= _MAX_UNKNOWN_WORD_LENGTH {
					break
				}
			}
			pair := dic.UnkIndex[class]
			for i, w := bufPos, 0; i < endPos; i += w {
				_, w = utf8.DecodeRune(this.input[i:])
				end := i + w
				for x := 0; x < pair[1]; x++ {
					this.addNode(pair[0]+x, chPos, chPos, this.input[bufPos:end], UNKNOWN)
				}
			}
		}
	}
	return
}

func (this *lattice) String() string {
	str := ""
	for i, nodes := range this.list {
		str += fmt.Sprintf("[%v] :\n", i)
		for _, node := range nodes {
			str += fmt.Sprintf("%v\n", node)
		}
		str += "\n"
	}
	return str
}

func (this *lattice) forward() (err error) {
	for i, size := 1, len(this.list); i < size; i++ {
		currentList := this.list[i]
		for index, target := range currentList {
			prevList := this.list[target.start]
			if len(prevList) == 0 {
				this.list[i][index].cost = _MAX_INT32
				continue
			}
			for j, n := range prevList {
				var c int16
				if n.class != USER && target.class != USER {
					c, err = dic.Connection.At(int(n.right), int(target.left))
				}
				if err != nil {
					err = fmt.Errorf("lattice.forward(): dic.Connection.At(%d, %d), %v", n.right, target.left, err)
					return
				}
				totalCost := int64(c) + int64(target.weight) + int64(n.cost)
				if totalCost > _MAX_INT32 {
					totalCost = _MAX_INT32
				}
				if j == 0 || int32(totalCost) < this.list[i][index].cost {
					this.list[i][index].cost = int32(totalCost)
					this.list[i][index].prev = this.list[target.start][j]
				}
			}
		}
	}
	return
}

func (this *lattice) backward() {
	size := len(this.list)
	this.output = make([]*Node, 0, size)
	for p := this.list[size-1][0]; p != nil; p = p.prev {
		this.output = append(this.output, p)
	}
}
