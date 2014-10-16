package kagome

import (
	"fmt"
	"io"
	"unicode/utf8"
)

const (
	maximumCost              = 1<<31 - 1
	maximumUnknownWordLength = 1024
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
		class := la.dic.CharCategory[ch]
		if !anyMatches || la.dic.InvokeList[class] {
			endPos := pos + utf8.RuneLen(ch)
			unkWordLen := 1
			for i, w, size := endPos, 1, len(input); i < size; i += w {
				var c rune
				c, w = utf8.DecodeRuneInString(input[i:])
				if la.dic.CharCategory[c] != class {
					break
				}
				endPos += w
				unkWordLen++
				if unkWordLen >= maximumUnknownWordLength {
					break
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

func (la *lattice) forward() {
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

func (la *lattice) backward() {
	size := len(la.list)
	if cap(la.output) < size {
		la.output = make([]*node, 0, size)
	} else {
		la.output = la.output[:0]
	}
	for p := la.list[size-1][0]; p != nil; p = p.prev {
		la.output = append(la.output, p)
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
