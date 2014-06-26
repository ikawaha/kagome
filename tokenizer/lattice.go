package tokenizer

import (
	"github.com/ikawaha/kagome/dic"

	"fmt"
	"unicode/utf8"
)

const (
	_MAX_INT32               = 1<<31 - 1
	_MAX_UNKNOWN_WORD_LENGTH = 1024
)

type Lattice struct {
	input  []byte
	list   [][]Node
	output []*Node
}

func (this *Lattice) addNode(a_id, a_pos, a_start, a_end int, a_class NodeType) {
	var cost dic.Cost
	if a_class == KNOWN {
		cost = dic.Costs[a_id]
	} else {
		cost = dic.UnkCosts[a_id]
	}
	node := Node{
		id:      a_id,
		start:   a_pos,
		class:   a_class,
		left:    int32(cost.Left),
		right:   int32(cost.Right),
		weight:  int32(cost.Weight),
		surface: this.input[a_start:a_end],
	}
	p := a_pos + utf8.RuneCount(node.surface)
	this.list[p] = append(this.list[p], node)
}

func (this *Lattice) build(a_input *string) (err error) {
	this.input = []byte(*a_input)

	runeCount := utf8.RuneCount(this.input)
	this.list = make([][]Node, runeCount+2)

	this.list[0] = append(this.list[0], Node{id: BOSEOS, class: KNOWN, start: 0})
	this.list[len(this.list)-1] = append(this.list[len(this.list)-1], Node{id: BOSEOS, class: KNOWN, start: len(this.list) - 2})

	chPos := -1
	for bufPos, ch := range *a_input {
		chPos++

		// (1) TODO: USER DIC
		anyMatches := false
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
				this.addNode(id+x, chPos, bufPos, bufPos+len(substr), KNOWN)
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
					this.addNode(pair[0]+x, chPos, bufPos, end, UNKNOWN)
				}
			}
		}
	}
	return
}

func (this *Lattice) String() string {
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

func (this *Lattice) forward() (err error) {
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
				c, err = dic.Connection.At(int(n.right), int(target.left))
				if err != nil {
					err = fmt.Errorf("Lattice.forward(): dic.Connection.At(%d, %d), %v", n.right, target.left, err)
					return
				}
				totalCost := int64(c) + int64(target.weight) + int64(n.cost)
				if totalCost > _MAX_INT32 {
					totalCost = _MAX_INT32
				}
				if j == 0 || int32(totalCost) < this.list[i][index].cost {
					this.list[i][index].cost = int32(totalCost)
					this.list[i][index].prev = &this.list[target.start][j]
				}
			}
		}
	}
	return
}

func (this *Lattice) backward() {
	size := len(this.list)
	stack := make([]*Node, 0, size)
	p := &this.list[size-1][0]

	stack = append(stack, p)
	for p != nil {
		stack = append(stack, p)
		p = p.prev
	}

	this.output = make([]*Node, 0, len(stack))
	for i := len(stack) - 1; i > 0; i-- {
		this.output = append(this.output, stack[i])
	}
}
