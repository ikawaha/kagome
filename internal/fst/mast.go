// Copyright 2015 ikawaha
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// 	You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fst

import (
	"fmt"
	"io"
	"sort"
)

// mast represents a Minimal Acyclic Subsequential Transeducer.
type mast struct {
	initialState *state
	states       []*state
	finalStates  []*state
}

func (m *mast) addState(n *state) {
	n.ID = len(m.states)
	m.states = append(m.states, n)
	if n.IsFinal {
		m.finalStates = append(m.finalStates, n)
	}
}

// Build constructs a virtual machine of a finite state transducer from a given inputs.
func Build(input PairSlice) (t FST, err error) {
	m := buildMAST(input)
	return m.buildMachine()
}

func commonPrefix(a, b string) string {
	end := len(a)
	if end > len(b) {
		end = len(b)
	}
	var i int
	for i < end && a[i] == b[i] {
		i++
	}
	return a[0:i]
}

func buildMAST(input PairSlice) (m mast) {
	sort.Sort(input)

	const initialMASTSize = 1024
	dic := make(map[int64][]*state)
	m.states = make([]*state, 0, initialMASTSize)
	m.finalStates = make([]*state, 0, initialMASTSize)

	buf := make([]*state, input.maxInputWordLen()+1)
	for i := range buf {
		buf[i] = newState()
	}
	prev := ""
	for _, pair := range input {
		in, out := pair.In, pair.Out
		fZero := (out == 0) // flag
		prefixLen := len(commonPrefix(in, prev))
		for i := len(prev); i > prefixLen; i-- {
			var s *state
			if cs, ok := dic[buf[i].hcode]; ok {
				for _, c := range cs {
					if c.eq(buf[i]) {
						s = c
						break
					}
				}
			}
			if s == nil {
				s = &state{}
				*s = *buf[i]
				m.addState(s)
				dic[s.hcode] = append(dic[s.hcode], s)
			}
			buf[i].renew()
			buf[i-1].setTransition(prev[i-1], s)
		}
		for i, size := prefixLen+1, len(in); i <= size; i++ {
			buf[i-1].setTransition(in[i-1], buf[i])
		}
		if in != prev {
			buf[len(in)].IsFinal = true
		}
		for j := 1; j < prefixLen+1; j++ {
			outSuff, ok := buf[j-1].Output[in[j-1]]
			if ok {
				if outSuff == out {
					out = 0
					break
				}
				buf[j-1].removeOutput(in[j-1]) // clear the prev edge
				for ch := range buf[j].Trans {
					buf[j].setOutput(ch, outSuff)
				}
				if buf[j].IsFinal && outSuff != 0 {
					buf[j].addTail(outSuff)
				}
			}
		}
		if in != prev {
			buf[prefixLen].setOutput(in[prefixLen], out)
		} else if fZero || out != 0 {
			buf[len(in)].addTail(out)
		}
		prev = in
	}
	// flush the buf
	for i := len(prev); i > 0; i-- {
		var s *state
		if cs, ok := dic[buf[i].hcode]; ok {
			for _, c := range cs {
				if c.eq(buf[i]) {
					s = c
					break
				}
			}
		}
		if s == nil {
			s = &state{}
			*s = *buf[i]
			buf[i].renew()
			m.addState(s)
			dic[s.hcode] = append(dic[s.hcode], s)
		}
		buf[i-1].setTransition(prev[i-1], s)
	}
	m.initialState = buf[0]
	m.addState(buf[0])

	return
}

func (m *mast) run(input string) (out []int32, ok bool) {
	s := m.initialState
	for i, size := 0, len(input); i < size; i++ {
		if o, ok := s.Output[input[i]]; ok {
			out = append(out, o)
		}
		if s, ok = s.Trans[input[i]]; !ok {
			return
		}
	}
	for _, t := range s.tails() {
		out = append(out, t)
	}
	return
}

func (m *mast) accept(input string) (ok bool) {
	s := m.initialState
	for i, size := 0, len(input); i < size; i++ {
		if s, ok = s.Trans[input[i]]; !ok {
			return
		}
	}
	return
}

func (m *mast) dot(w io.Writer) {
	fmt.Fprintln(w, "digraph G {")
	fmt.Fprintln(w, "\trankdir=LR;")
	fmt.Fprintln(w, "\tnode [shape=circle]")
	for _, s := range m.finalStates {
		fmt.Fprintf(w, "\t%d [peripheries = 2];\n", s.ID)
	}
	for _, from := range m.states {
		for in, to := range from.Trans {
			fmt.Fprintf(w, "\t%d -> %d [label=\"%02X/%v", from.ID, to.ID, in, from.Output[in])
			if to.hasTail() {
				fmt.Fprintf(w, " %v", to.tails())
			}
			fmt.Fprintln(w, "\"];")
		}
	}
	fmt.Fprintln(w, "}")
}
