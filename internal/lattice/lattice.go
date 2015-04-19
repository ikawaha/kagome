//  Copyright (c) 2015 ikawaha.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package lattice

import (
	"fmt"
	"io"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/ikawaha/kagome/internal/dic"
)

const (
	maximumCost              = 1<<31 - 1
	maximumUnknownWordLength = 1024
	searchModeKanjiLength    = 2
	searchModeKanjiPenalty   = 3000
	searchModeOtherLength    = 7
	searchModeOtherPenalty   = 1700
)

// TokenizeMode represents how to tokenize sentencse.
type TokenizeMode int

const (
	Normal   TokenizeMode = iota + 1 //Normal Mode
	Search                           // Search Mode
	Extended                         // Extended Mode ()
)

var latticePool = sync.Pool{
	New: func() interface{} {
		return new(lattice)
	},
}

type lattice struct {
	Input  string
	Output []*node
	list   [][]*node
	dic    *dic.Dic
	udic   *dic.UserDic
}

func New(d *dic.Dic, u *dic.UserDic) *lattice {
	la := latticePool.Get().(*lattice)
	la.dic = d
	la.udic = u
	return la
}

func (la *lattice) Free() {
	for i := range la.list {
		for j := range la.list[i] {
			nodePool.Put(la.list[i][j])
		}
		la.list[i] = la.list[i][:0]
	}
	la.list = la.list[:0]
	la.udic = nil
	latticePool.Put(la)
}

func (la *lattice) addNode(pos, id, start int, class NodeClass, surface string) {
	var m dic.Morph
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
	n := newNode()
	n.ID = id
	n.Start = start
	n.Class = class
	n.Left, n.Right, n.Weight = int32(m.LeftID), int32(m.RightID), int32(m.Weight)
	n.Surface = surface
	n.Prev = nil
	p := pos + utf8.RuneCountInString(surface)
	la.list[p] = append(la.list[p], n)
}

func (la *lattice) Build(Input string) {
	rc := utf8.RuneCountInString(Input)
	la.Input = Input
	if cap(la.list) < rc+2 {
		const expandRatio = 2
		la.list = make([][]*node, 0, (rc+2)*expandRatio)
	}
	la.list = la.list[0 : rc+2]

	la.addNode(0, BosEosID, 0, DUMMY, Input[0:0])
	la.addNode(rc+1, BosEosID, rc, DUMMY, Input[rc:rc])

	runePos := -1
	for pos, ch := range Input {
		runePos++
		anyMatches := false

		// (1) USER DIC
		if la.udic != nil {
			lens, outputs := la.udic.Index.CommonPrefixSearch(Input[pos:])
			for i, ids := range outputs {
				for j := range ids {
					la.addNode(runePos, int(ids[j]), runePos,
						USER, Input[pos:pos+int(lens[i])])
				}
			}
			anyMatches = (len(lens) > 0)
		}
		if anyMatches {
			continue
		}

		// (2) KNOWN DIC
		if lens, outputs := la.dic.Index.CommonPrefixSearch(Input[pos:]); len(lens) > 0 {
			anyMatches = true
			for i, ids := range outputs {
				for j := range ids {
					la.addNode(runePos, int(ids[j]), runePos,
						KNOWN, Input[pos:pos+lens[i]])
				}
			}
		}
		// (3) UNKNOWN DIC
		class := la.dic.CharactorCategory(ch)
		if !anyMatches || la.dic.InvokeList[class] {
			endPos := pos + utf8.RuneLen(ch)
			unkWordLen := 1
			if la.dic.GroupList[class] {
				for i, w, size := endPos, 1, len(Input); i < size; i += w {
					var c rune
					c, w = utf8.DecodeRuneInString(Input[i:])
					if la.dic.CharactorCategory(c) != class {
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
				_, w = utf8.DecodeRuneInString(Input[i:])
				end := i + w
				dup, ok := la.dic.UnkIndexDup[int(class)]
				if !ok {
					dup = 1
				}
				for x := 0; x < dup; x++ {
					la.addNode(runePos, id+x, runePos,
						UNKNOWN, Input[pos:end])
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
	l := utf8.RuneCountInString(n.Surface)
	if l > searchModeKanjiLength && kanjiOnly(n.Surface) {
		return (l - searchModeKanjiLength) * searchModeKanjiPenalty
	}
	if l > searchModeOtherLength {
		return (l - searchModeOtherLength) * searchModeOtherPenalty
	}
	return 0
}

func (la *lattice) Forward(m TokenizeMode) {
	for i, size := 1, len(la.list); i < size; i++ {
		currentList := la.list[i]
		for index, target := range currentList {
			prevList := la.list[target.Start]
			if len(prevList) == 0 {
				la.list[i][index].Cost = maximumCost
				continue
			}
			for j, n := range prevList {
				var c int16
				if n.Class != USER && target.Class != USER {
					c = la.dic.Connection.At(int(n.Right), int(target.Left))
				}
				totalCost := int64(c) + int64(target.Weight) + int64(n.Cost)
				if m != Normal {
					totalCost += int64(additionalCost(n))
				}
				if totalCost > maximumCost {
					totalCost = maximumCost
				}
				if j == 0 || int32(totalCost) < la.list[i][index].Cost {
					la.list[i][index].Cost = int32(totalCost)
					la.list[i][index].Prev = la.list[target.Start][j]
				}
			}
		}
	}
	return
}

func (la *lattice) Backward(m TokenizeMode) {
	const bufferExpandRatio = 2
	size := len(la.list)
	if cap(la.Output) < size {
		la.Output = make([]*node, 0, size*bufferExpandRatio)
	} else {
		la.Output = la.Output[:0]
	}
	for p := la.list[size-1][0]; p != nil; p = p.Prev {
		if m != Extended || p.Class != UNKNOWN {
			la.Output = append(la.Output, p)
			continue
		}
		runeLen := utf8.RuneCountInString(p.Surface)
		stack := make([]*node, 0, runeLen)
		i := 0
		for _, r := range p.Surface {
			n := nodePool.Get().(*node)
			n.ID = p.ID
			n.Start = p.Start + i
			n.Class = DUMMY
			n.Surface = string(r)
			stack = append(stack, n)
			i++
		}
		for j, end := 0, len(stack); j < end; j++ {
			la.Output = append(la.Output, stack[runeLen-1-j])
		}
	}
}

func (la *lattice) Dot(w io.Writer) {
	type edge struct {
		from *node
		to   *node
	}
	edges := make([]edge, 0, 1024)
	for i, size := 1, len(la.list); i < size; i++ {
		currents := la.list[i]
		for _, to := range currents {
			prevs := la.list[to.Start]
			if len(prevs) == 0 {
				continue
			}
			for _, from := range prevs {
				edges = append(edges, edge{from, to})
			}
		}
	}
	bests := make(map[*node]struct{})
	for _, n := range la.Output {
		bests[n] = struct{}{}
	}
	fmt.Fprintln(w, "graph lattice {")
	fmt.Fprintln(w, "\tdpi=48;")
	fmt.Fprintln(w, "\tgraph [style=filled, rankdir=LR]")
	for i, list := range la.list {
		for _, n := range list {
			surf := n.Surface
			if n.ID == BosEosID {
				if i == 0 {
					surf = "BOS"
				} else {
					surf = "EOS"
				}
			}
			fmt.Fprintf(w, "\t\"%p\" [label=\"%s\\n%d\"];\n", n, surf, n.Weight)
		}
	}
	for _, e := range edges {
		var c int16
		if e.from.Class != USER && e.to.Class != USER {
			c = la.dic.Connection.At(int(e.from.Right), int(e.to.Left))
		}
		_, l := bests[e.from]
		_, r := bests[e.to]
		if l && r {
			fmt.Fprintf(w, "\t\"%p\" -- \"%p\" [label=\"%d\",color=blue,style=bold];\n",
				e.from, e.to, c)
		} else {
			fmt.Fprintf(w, "\t\"%p\" -- \"%p\" [label=\"%d\"];\n",
				e.from, e.to, c)
		}
	}

	fmt.Fprintln(w, "}")
}
