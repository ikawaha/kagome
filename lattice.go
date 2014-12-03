//  Copyright (c) 2014 ikawaha.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package kagome

import (
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

const (
	maximumCost              = 1<<31 - 1
	maximumUnknownWordLength = 1024
	searchModeKanjiLength    = 2
	searchModeKanjiPenalty   = 3000
	searchModeOtherLength    = 7
	searchModeOtherPenalty   = 1700
)

type lattice struct {
	input  string
	list   [][]*node
	output []*node
	pool   *nodePool
	dic    *Dic
	udic   *UserDic
}

func newLattice() (la *lattice) {
	la = new(lattice)
	la.dic = NewSysDic()
	return
}

func (la *lattice) setNodePool(size int) {
	la.pool = newNodePool(size)
}

func (la *lattice) setDic(dic *Dic) {
	if dic != nil {
		la.dic = dic
	}
}

func (la *lattice) setUserDic(udic *UserDic) {
	la.udic = udic
}

func (la *lattice) addNode(pos, id, start int, class NodeClass, surface string) {
	var m Morph
	switch class {
	case DUMMY:
		//use default cost
	case KNOWN:
		m = la.dic.Morphs[id]
	case UNKNOWN:
		m = la.dic.UnkMorphs[id]
	case USER:
		// use default cost
	}
	n := la.pool.get()
	n.id = id
	n.start = start
	n.class = class
	n.left, n.right, n.weight = int32(m.LeftId), int32(m.RightId), int32(m.Weight)
	n.surface = surface
	n.prev = nil
	p := pos + utf8.RuneCountInString(surface)
	la.list[p] = append(la.list[p], n)
}

func (la *lattice) build(input string) {
	rc := utf8.RuneCountInString(input)
	la.pool.clear()
	la.input = input
	la.list = make([][]*node, rc+2)

	la.addNode(0, BosEosId, 0, DUMMY, input[0:0])
	la.addNode(rc+1, BosEosId, rc, DUMMY, input[rc:rc])

	runePos := -1
	for pos, ch := range input {
		runePos++
		anyMatches := false

		// (1) TODO: USER DIC
		if la.udic != nil {
			ids, lens := la.udic.Index.CommonPrefixSearchString(input[pos:])
			for i := range ids {
				la.addNode(runePos, ids[i], runePos, USER, input[pos:pos+lens[i]])
			}
			anyMatches = (len(ids) > 0)
		}
		if anyMatches {
			continue
		}

		// (2) KNOWN DIC
		if ids, lens := la.dic.Index.CommonPrefixSearchString(input[pos:]); len(ids) > 0 {
			anyMatches = true
			for i, id := range ids {
				dup, ok := la.dic.IndexDup[id]
				if !ok {
					dup = 1
				}
				for x := 0; x < dup; x++ {
					la.addNode(runePos, id+x, runePos, KNOWN, input[pos:pos+lens[i]])
				}
			}
		}
		// (3) UNKNOWN DIC
		class := la.dic.charCategory(ch)
		if !anyMatches || la.dic.InvokeList[class] {
			endPos := pos + utf8.RuneLen(ch)
			unkWordLen := 1
			if la.dic.GroupList[class] {
				for i, w, size := endPos, 1, len(input); i < size; i += w {
					var c rune
					c, w = utf8.DecodeRuneInString(input[i:])
					if la.dic.charCategory(c) != class {
						break
					}
					endPos += w
					unkWordLen++
					if unkWordLen >= maximumUnknownWordLength {
						break
					}
				}
			}
			id := la.dic.UnkIndex[int(class)]
			for i, w := pos, 0; i < endPos; i += w {
				_, w = utf8.DecodeRuneInString(input[i:])
				end := i + w
				dup, ok := la.dic.UnkIndexDup[int(class)]
				if !ok {
					dup = 1
				}
				for x := 0; x < dup; x++ {
					la.addNode(runePos, id+x, runePos, UNKNOWN, input[pos:end])
				}
			}
		}
	}
	return
}

// String returns a debug string of a lattice.
func (la *lattice) String() string {
	str := ""
	for i, nodes := range la.list {
		str += fmt.Sprintf("[%v] :\n", i)
		for _, node := range nodes {
			str += fmt.Sprintf("%v\n", node)
		}
		str += "\n"
	}
	return str
}
func kanjiOnly(s string) bool {
	for _, r := range s {
		if !unicode.In(r, unicode.Ideographic) {
			return false
		}
	}
	return s != ""
}

func additionalCost(n *node) int {
	l := utf8.RuneCountInString(n.surface)
	if l > searchModeKanjiLength && kanjiOnly(n.surface) {
		return (l - searchModeKanjiLength) * searchModeKanjiPenalty
	}
	if l > searchModeOtherLength {
		return (l - searchModeOtherLength) * searchModeOtherPenalty
	}
	return 0
}

func (la *lattice) forward(mode tokenizeMode) {
	for i, size := 1, len(la.list); i < size; i++ {
		currentList := la.list[i]
		for index, target := range currentList {
			prevList := la.list[target.start]
			if len(prevList) == 0 {
				la.list[i][index].cost = maximumCost
				continue
			}
			for j, n := range prevList {
				var c int16
				if n.class != USER && target.class != USER {
					c = la.dic.Connection.At(int(n.right), int(target.left))
				}
				totalCost := int64(c) + int64(target.weight) + int64(n.cost)
				if mode != normalModeTokenize {
					totalCost += int64(additionalCost(n))
				}
				if totalCost > maximumCost {
					totalCost = maximumCost
				}
				if j == 0 || int32(totalCost) < la.list[i][index].cost {
					la.list[i][index].cost = int32(totalCost)
					la.list[i][index].prev = la.list[target.start][j]
				}
			}
		}
	}
	return
}

func (la *lattice) backward(mode tokenizeMode) {
	const bufferExpandRatio = 2
	size := len(la.list)
	if cap(la.output) < size {
		la.output = make([]*node, 0, size*bufferExpandRatio)
	} else {
		la.output = la.output[:0]
	}
	for p := la.list[size-1][0]; p != nil; p = p.prev {
		if mode != extendedModeTokenize || p.class != UNKNOWN {
			la.output = append(la.output, p)
			continue
		}
		runeLen := utf8.RuneCountInString(p.surface)
		stack := make([]*node, 0, runeLen)
		i := 0
		for _, r := range p.surface {
			n := la.pool.get()
			n.id = p.id
			n.start = p.start + i
			n.class = DUMMY
			n.surface = string(r)
			stack = append(stack, n)
			i++
		}
		for j, end := 0, len(stack); j < end; j++ {
			la.output = append(la.output, stack[runeLen-1-j])
		}
	}
}

func (la *lattice) dot(w io.Writer) {
	type edge struct {
		from *node
		to   *node
	}
	edges := make([]edge, 0, 1024)
	for i, size := 1, len(la.list); i < size; i++ {
		currents := la.list[i]
		for _, to := range currents {
			prevs := la.list[to.start]
			if len(prevs) == 0 {
				continue
			}
			for _, from := range prevs {
				edges = append(edges, edge{from, to})
			}
		}
	}
	bests := make(map[*node]struct{})
	for _, n := range la.output {
		bests[n] = struct{}{}
	}
	fmt.Fprintln(w, "graph lattice {")
	fmt.Fprintln(w, "\tdpi=48;")
	fmt.Fprintln(w, "\tgraph [style=filled, rankdir=LR]")
	for i, list := range la.list {
		for _, n := range list {
			surf := n.surface
			if n.id == BosEosId {
				if i == 0 {
					surf = "BOS"
				} else {
					surf = "EOS"
				}
			}
			fmt.Fprintf(w, "\t\"%p\" [label=\"%s\\n%d\"];\n", n, surf, n.weight)
		}
	}
	for _, e := range edges {
		var c int16
		if e.from.class != USER && e.to.class != USER {
			c = la.dic.Connection.At(int(e.from.right), int(e.to.left))
		}
		_, l := bests[e.from]
		_, r := bests[e.to]
		if l && r {
			fmt.Fprintf(w, "\t\"%p\" -- \"%p\" [label=\"%d\",color=blue,style=bold];\n", e.from, e.to, c)
		} else {
			fmt.Fprintf(w, "\t\"%p\" -- \"%p\" [label=\"%d\"];\n", e.from, e.to, c)
		}
	}

	fmt.Fprintln(w, "}")
}
