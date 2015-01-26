package kagome

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"sort"
	"unsafe"
)

// Pair implements a pair of input and output.
type Pair struct {
	In  string
	Out int32
}

// PairSlice implements a slice of input and output pairs.
type PairSlice []Pair

func (ps PairSlice) Len() int {
	return len(ps)
}

func (ps PairSlice) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (ps PairSlice) Less(i, j int) bool {
	return ps[i].In < ps[j].In
}

func (ps PairSlice) maxInputWordLen() (max int) {
	for _, pair := range ps {
		if size := len(pair.In); size > max {
			max = size
		}
	}
	return
}

type int32Set map[int32]bool

type state struct {
	ID      int
	Trans   map[byte]*state
	Output  map[byte]int32
	Tail    int32Set
	IsFinal bool
	hcode   int64
}

func newState() (n *state) {
	n = new(state)
	n.Trans = make(map[byte]*state)
	n.Output = make(map[byte]int32)
	n.Tail = make(int32Set)
	return
}

func (n *state) hasTail() bool {
	return len(n.Tail) != 0
}

func (n *state) addTail(t int32) {
	n.Tail[t] = true
}

func (n *state) tails() []int32 {
	t := make([]int32, 0, len(n.Tail))
	for item := range n.Tail {
		t = append(t, item)
	}
	return t
}

func (n *state) removeOutput(ch byte) {
	const magic = 8191
	if out, ok := n.Output[ch]; ok && out != 0 {
		n.hcode -= (int64(ch) + int64(out)) * magic
	}
	delete(n.Output, ch)
}

func (n *state) setOutput(ch byte, out int32) {
	if out == 0 {
		return
	}
	n.Output[ch] = out

	const magic = 8191
	n.hcode += (int64(ch) + int64(out)) * magic
}

func (n *state) setTransition(ch byte, next *state) {
	n.Trans[ch] = next

	const magic = 1001
	n.hcode += (int64(ch) + int64(next.ID)) * magic
}

func (n *state) renew() {
	n.Trans = make(map[byte]*state)
	n.Output = make(map[byte]int32)
	n.Tail = make(int32Set)
	n.IsFinal = false
	n.hcode = 0
}

func (n *state) eq(dst *state) bool {
	if n == nil || dst == nil {
		return false
	}
	if n == dst {
		return true
	}
	if n.hcode != dst.hcode {
		return false
	}
	if len(n.Trans) != len(dst.Trans) ||
		len(n.Output) != len(dst.Output) ||
		len(n.Tail) != len(dst.Tail) ||
		n.IsFinal != dst.IsFinal {
		return false
	}
	for ch, next := range n.Trans {
		if dst.Trans[ch] != next {
			return false
		}
	}
	for ch, out := range n.Output {
		if dst.Output[ch] != out {
			return false
		}
	}
	for item := range n.Tail {
		if !dst.Tail[item] {
			return false
		}
	}
	return true
}

// String returns a string representaion of a node for debug.
func (n *state) String() string {
	ret := ""
	if n == nil {
		return "<nil>"
	}
	ret += fmt.Sprintf("%d[%p]:", n.ID, n)
	for ch := range n.Trans {
		ret += fmt.Sprintf("%X02/%v -->%p, ", ch, n.Output[ch], n.Trans[ch])
	}
	if n.IsFinal {
		ret += fmt.Sprintf(" (tail:%v) ", n.tails())
	}
	return ret
}

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

// BuildFST constructs a virtual machine of a finite state transducer from a given inputs.
func BuildFST(input PairSlice) (t FST, err error) {
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
			if buf[j-1].Output[in[j-1]] == out {
				out = 0
				break
			}
			var outSuff int32
			outSuff = buf[j-1].Output[in[j-1]]
			buf[j-1].removeOutput(in[j-1]) // clear the prev edge
			for ch := range buf[j].Trans {
				buf[j].setOutput(ch, outSuff)
			}
			if buf[j].IsFinal && outSuff != 0 {
				buf[j].addTail(outSuff)
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

type operation byte

const (
	opAccept      operation = 1
	opMatch       operation = 2
	opBreak       operation = 3
	opOutput      operation = 4
	opOutputBreak operation = 5
)

func (o operation) String() string {
	opName := []string{"OP0", "ACC", "MTC", "BRK", "OUT", "OUB", "OP6", "OP7"}
	if int(o) >= len(opName) {
		return fmt.Sprintf("NA[%d]", o)
	}
	return opName[o]
}

type instruction [4]byte

// FST represents a finite state transducer (virtual machine).
type FST struct {
	prog []instruction
	data []int32
}

// Configuration represents a FST (virtual machine) configuration.
type configuration struct {
	pc  int     // program counter
	hd  int     // input head
	out []int32 // outputs
}

const maxUint16 = 1<<16 - 1

type int32Slice []int32

func (p int32Slice) Len() int           { return len(p) }
func (p int32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p int32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type byteSlice []byte

func (p byteSlice) Len() int           { return len(p) }
func (p byteSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p byteSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func invert(prog []instruction) []instruction {
	size := len(prog)
	inv := make([]instruction, size)
	for i := range prog {
		inv[i] = prog[size-1-i]
	}
	return inv
}

func (m mast) buildMachine() (t FST, err error) {
	var (
		prog []instruction
		data []int32
		code instruction // tmp instruction
	)
	var edges []byte
	addrMap := make(map[int]int)
	for _, s := range m.states {
		edges = edges[:0]
		for ch := range s.Trans {
			edges = append(edges, ch)
		}
		if len(edges) > 0 {
			sort.Sort(byteSlice(edges))
		}
		for i, size := 0, len(edges); i < size; i++ {
			ch := edges[size-1-i]
			next := s.Trans[ch]
			out := s.Output[ch]
			addr, ok := addrMap[next.ID]
			if !ok && !next.IsFinal {
				err = fmt.Errorf("next addr is undefined: state(%v), input(%X)", s.ID, ch)
				return
			}
			jump := len(prog) - addr + 1
			var op operation
			if out != 0 {
				if i == 0 {
					op = opOutputBreak
				} else {
					op = opOutput
				}
			} else if i == 0 {
				op = opBreak
			} else {
				op = opMatch
			}

			if jump > maxUint16 {
				p := unsafe.Pointer(&code[0])
				(*(*int32)(p)) = int32(jump)
				prog = append(prog, code)
				jump = 0
			}
			if out != 0 {
				p := unsafe.Pointer(&code[0])
				(*(*int32)(p)) = int32(out)
				prog = append(prog, code)
			}

			code[0] = byte(op)
			code[1] = ch
			p := unsafe.Pointer(&code[2])
			(*(*uint16)(p)) = uint16(jump)
			prog = append(prog, code)
		}
		if s.IsFinal {
			if len(s.Tail) > 0 {
				p := unsafe.Pointer(&code[0])
				(*(*int32)(p)) = int32(len(data))
				prog = append(prog, code)
				var tmp int32Slice
				for t := range s.Tail {
					tmp = append(tmp, t)
				}
				sort.Sort(tmp)
				data = append(data, tmp...)
				p = unsafe.Pointer(&code[0])
				(*(*int32)(p)) = int32(len(data))
				prog = append(prog, code)
			}
			code[0] = byte(opAccept)
			code[1], code[2], code[3] = 0, 0, 0 // clear
			if len(s.Tail) > 0 {
				code[1] = 1
			}

			prog = append(prog, code)
		}
		addrMap[s.ID] = len(prog)
	}
	t = FST{prog: invert(prog), data: data}
	return
}

// String returns debug codes of a fst virtual machine.
func (t FST) String() string {
	var (
		pc   int
		code instruction
		op   operation
		ch   byte
		v16  uint16
		v32  int32
	)
	ret := ""
	for pc = 0; pc < len(t.prog); pc++ {
		code = t.prog[pc]
		op = operation(code[0])
		ch = code[1]
		v16 = (*(*uint16)(unsafe.Pointer(&code[2])))
		switch operation(op) {
		case opAccept:
			//fmt.Printf("%3d %v\t%X %d\n", pc, op, ch, v16) //XXX
			ret += fmt.Sprintf("%3d %v\t%d %d\n", pc, op, ch, v16)
			if ch == 0 {
				break
			}
			pc++
			code = t.prog[pc]
			to := (*(*int32)(unsafe.Pointer(&code[0])))
			ret += fmt.Sprintf("%3d [%d]\n", pc, to)
			pc++
			code = t.prog[pc]
			from := (*(*int32)(unsafe.Pointer(&code[0])))
			ret += fmt.Sprintf("%3d [%d] %v\n", pc, from, t.data[from:to]) //FIXME
		case opMatch:
			fallthrough
		case opBreak:
			//fmt.Printf("%3d %v\t%02X %d\n", pc, op, ch, v16) //XXX
			ret += fmt.Sprintf("%3d %v\t%02X(%c) %d\n", pc, op, ch, ch, v16)
			if v16 == 0 {
				pc++
				code = t.prog[pc]
				v32 = (*(*int32)(unsafe.Pointer(&code[0])))
				//fmt.Printf("%3d [%d]\n", pc, v32) //XXX
				ret += fmt.Sprintf("%3d jmp[%d]\n", pc, v32)
				//break
			}
		case opOutput:
			fallthrough
		case opOutputBreak:
			//fmt.Printf("%3d %v\t%02X %d\n", pc, op, ch, v16) //XXX
			ret += fmt.Sprintf("%3d %v\t%02X(%c) %d\n", pc, op, ch, ch, v16)
			if v16 == 0 {
				pc++
				code = t.prog[pc]
				v32 = (*(*int32)(unsafe.Pointer(&code[0])))
				//fmt.Printf("%3d [%d]\n", pc, v32) //XXX
				ret += fmt.Sprintf("%3d jmp[%d]\n", pc, v32)
				//break
			}
			pc++
			code = t.prog[pc]
			v32 = (*(*int32)(unsafe.Pointer(&code[0])))
			//fmt.Printf("%3d [%d]\n", pc, v32) //XXX
			ret += fmt.Sprintf("%3d [%d]\n", pc, v32)
		default:
			//fmt.Printf("%3d UNDEF %v\n", pc, code)
			ret += fmt.Sprintf("%3d UNDEF %v\n", pc, code)
		}
	}
	return ret
}

func (t *FST) run(input string) (snap []configuration, accept bool) {
	var (
		pc  int       // program counter
		op  operation // operation
		ch  byte      // char
		v16 uint16    // 16bit register
		v32 int32     // 32bit register
		hd  int       // input head
		out int32     // output

		code instruction // tmp instruction
	)
	for pc < len(t.prog) && hd <= len(input) {
		code = t.prog[pc]
		op = operation(code[0])
		ch = code[1]
		v16 = (*(*uint16)(unsafe.Pointer(&code[2])))
		//fmt.Printf("pc:%v,op:%v,hd:%v,v16:%v,out:%v\n", pc, op, hd, v16, out) //XXX
		switch op {
		case opMatch:
			fallthrough
		case opBreak:
			if hd == len(input) {
				goto L_END
			}
			if ch != input[hd] {
				if op == opBreak {
					return
				}
				if v16 == 0 {
					pc++
				}
				pc++
				continue
			}
			if v16 > 0 {
				pc += int(v16)
			} else {
				pc++
				code = t.prog[pc]
				v32 = (*(*int32)(unsafe.Pointer(&code[0])))
				//fmt.Printf("ex jump:%d\n", v32) //XXX
				pc += int(v32)
			}
			hd++
			continue
		case opOutput:
			fallthrough
		case opOutputBreak:
			if hd == len(input) {
				goto L_END
			}
			if ch != input[hd] {
				if op == opOutputBreak {
					return
				}
				if v16 == 0 {
					pc++
				}
				pc++
				pc++
				continue
			}
			pc++
			code = t.prog[pc]
			out = (*(*int32)(unsafe.Pointer(&code[0])))
			if v16 > 0 {
				pc += int(v16)
			} else {
				pc++
				code = t.prog[pc]
				v32 = (*(*int32)(unsafe.Pointer(&code[0])))
				//fmt.Printf("ex jump:%d\n", v32) //XXX
				pc += int(v32)
			}
			hd++
			continue
		case opAccept:
			c := configuration{pc: pc, hd: hd}
			pc++
			if ch == 0 {
				c.out = []int32{out}
			} else {
				code = t.prog[pc]
				to := (*(*int32)(unsafe.Pointer(&code[0])))
				pc++
				code = t.prog[pc]
				from := (*(*int32)(unsafe.Pointer(&code[0])))
				c.out = t.data[from:to]
				pc++
			}
			snap = append(snap, c)
			if hd == len(input) {
				goto L_END
			}
			continue
		default:
			//fmt.Printf("unknown op:%v\n", op) //XXX
			return
		}
	}
L_END:
	if hd != len(input) {
		return
	}
	if op != opAccept {
		//fmt.Printf("[[FINAL]]pc:%d, op:%s, ch:[%X], sz:%d, v:%d\n", pc, op, ch, sz, va) //XXX
		return

	}
	accept = true
	return
}

// Search runs a finite state transducer for a given input and returns outputs if accepted otherwise nil.
func (t FST) Search(input string) []int32 {
	snap, acc := t.run(input)
	if !acc || len(snap) == 0 {
		return nil
	}
	c := snap[len(snap)-1]
	return c.out
}

// PrefixSearch returns the longest commom prefix keyword and it's length in given input
// if detected otherwise -1, nil.
func (t FST) PrefixSearch(input string) (length int, output []int32) {
	snap, _ := t.run(input)
	if len(snap) == 0 {
		return -1, nil
	}
	c := snap[len(snap)-1]
	return c.hd, c.out

}

// CommonPrefixSearch finds keywords sharing common prefix in given input
// and returns it's lengths and outputs. Returns nil, nil if there does not common prefix keywords.
func (t FST) CommonPrefixSearch(input string) (lens []int, outputs [][]int32) {
	snap, _ := t.run(input)
	if len(snap) == 0 {
		return
	}
	for _, c := range snap {
		lens = append(lens, c.hd)
		outputs = append(outputs, c.out)
	}
	return

}

// Write saves a program of finite state transducer (virtual machine)
func (t FST) Write(w io.Writer) error {
	var (
		pc   int
		code instruction
		op   operation
		ch   byte
		v16  uint16
		v32  int32
	)
	dataLen := int64(len(t.data))
	//fmt.Println("data len:", dataLen)
	if e := binary.Write(w, binary.LittleEndian, dataLen); e != nil {
		return e
	}
	//fmt.Println("data len:", dataLen) //XXX
	for _, v := range t.data {
		if e := binary.Write(w, binary.LittleEndian, v); e != nil {
			return e
		}
	}

	progLen := int64(len(t.prog))
	if e := binary.Write(w, binary.LittleEndian, progLen); e != nil {
		return e
	}
	//fmt.Println("prog len:", progLen) //XXX
	for pc = 0; pc < len(t.prog); pc++ {
		code = t.prog[pc]
		op = operation(code[0])
		ch = code[1]
		v16 = (*(*uint16)(unsafe.Pointer(&code[2])))

		// write op and ch
		if _, e := w.Write(code[0:2]); e != nil {
			return e
		}
		//fmt.Printf("%3d %v\t%X %d\n", pc, op, ch, v16) //XXX
		switch operation(op) {
		case opAccept:
			if ch == 0 {
				break
			}
			pc++
			code = t.prog[pc]
			v32 = (*(*int32)(unsafe.Pointer(&code[0]))) //to addr
			if e := binary.Write(w, binary.LittleEndian, v32); e != nil {
				return e
			}
			//fmt.Printf("%3d \t[%d]\n", pc, v32) //XXX
			pc++
			code = t.prog[pc]
			v32 = (*(*int32)(unsafe.Pointer(&code[0]))) //from addr
			if e := binary.Write(w, binary.LittleEndian, v32); e != nil {
				return e
			}
			//fmt.Printf("%3d \t[%d]\n", pc, v32) //XXX
		case opMatch:
			fallthrough
		case opBreak:
			if e := binary.Write(w, binary.LittleEndian, v16); e != nil {
				return e
			}
			if v16 != 0 {
				break
			}
			pc++
			code = t.prog[pc]
			v32 = (*(*int32)(unsafe.Pointer(&code[0])))
			if e := binary.Write(w, binary.LittleEndian, v32); e != nil {
				return e
			}
			//fmt.Printf("%3d \t[%d]\n", pc, v32) //XXX
		case opOutput:
			fallthrough
		case opOutputBreak:
			if e := binary.Write(w, binary.LittleEndian, v16); e != nil {
				return e
			}
			pc++
			code = t.prog[pc]
			v32 = (*(*int32)(unsafe.Pointer(&code[0])))
			if e := binary.Write(w, binary.LittleEndian, v32); e != nil {
				return e
			}
			//fmt.Printf("%3d \t[%d]\n", pc, v32) //XXX

			if v16 != 0 {
				break
			}
			pc++
			code = t.prog[pc]
			v32 = (*(*int32)(unsafe.Pointer(&code[0])))
			if e := binary.Write(w, binary.LittleEndian, v32); e != nil {
				return e
			}
			//fmt.Printf("%3d \t[%d]\n", pc, v32) //XXX
		default:
			return fmt.Errorf("undefined operation error")
		}
	}
	return nil
}

// Read loads a program of finite state transducer (virtual machine)
func (t *FST) Read(r io.Reader) (e error) {
	var (
		code instruction
		op   byte
		ch   byte
		v16  uint16
		v32  int32
		p    unsafe.Pointer
		//pc   int //XXX
	)

	rd := bufio.NewReader(r)

	var dataLen int64
	if e = binary.Read(rd, binary.LittleEndian, &dataLen); e != nil {
		return
	}
	//fmt.Println("data len:", dataLen) //XXX
	t.data = make([]int32, 0, dataLen)
	for i := 0; i < int(dataLen); i++ {
		if e = binary.Read(rd, binary.LittleEndian, &v32); e != nil {
			return
		}
		t.data = append(t.data, v32)
	}

	var progLen int64
	if e = binary.Read(rd, binary.LittleEndian, &progLen); e != nil {
		return
	}
	//fmt.Println("prog len:", progLen) //XXX
	t.prog = make([]instruction, 0, progLen)

	for e == nil {
		if op, e = rd.ReadByte(); e != nil {
			break
		}
		if ch, e = rd.ReadByte(); e != nil {
			break
		}
		switch operation(op) {
		case opAccept:
			code[0], code[1], code[2], code[3] = op, ch, 0, 0
			t.prog = append(t.prog, code)
			//fmt.Printf("%3d %v\t%X %d\n", pc, operation(op), ch, 0) //XXX
			//pc++                                                    //XXX
			if ch == 0 {
				break
			}
			if e = binary.Read(rd, binary.LittleEndian, &v32); e != nil {
				break
			}
			p = unsafe.Pointer(&code[0])
			(*(*int32)(p)) = int32(v32)
			//fmt.Printf("%3d \t[%d]\n", pc, v32) //XXX
			//pc++                                //XXX
			t.prog = append(t.prog, code)

			if e = binary.Read(rd, binary.LittleEndian, &v32); e != nil {
				break
			}
			p = unsafe.Pointer(&code[0])
			(*(*int32)(p)) = int32(v32)
			//fmt.Printf("%3d \t[%d]\n", pc, v32) //XXX
			//pc++                                //XXX
			t.prog = append(t.prog, code)
		case opMatch:
			fallthrough
		case opBreak:
			code[0], code[1] = op, ch
			if e = binary.Read(rd, binary.LittleEndian, &v16); e != nil {
				break
			}
			p = unsafe.Pointer(&code[2])
			(*(*uint16)(p)) = uint16(v16)
			//fmt.Printf("%3d %v\t%X %d\n", pc, operation(op), ch, v16) //XXX
			//pc++                                                      //XXX
			t.prog = append(t.prog, code)

			if v16 != 0 {
				break
			}
			if e = binary.Read(rd, binary.LittleEndian, &v32); e != nil {
				break
			}
			p = unsafe.Pointer(&code[0])
			(*(*int32)(p)) = int32(v32)
			//fmt.Printf("%3d \t[%d]\n", pc, v32) //XXX
			//pc++                                //XXX
			t.prog = append(t.prog, code)
		case opOutput:
			fallthrough
		case opOutputBreak:
			code[0], code[1] = op, ch
			if e = binary.Read(rd, binary.LittleEndian, &v16); e != nil {
				break
			}
			p = unsafe.Pointer(&code[2])
			(*(*uint16)(p)) = uint16(v16)
			//fmt.Printf("%3d %v\t%X %d\n", pc, operation(op), ch, v16) //XXX
			//pc++                                                      //XXX
			t.prog = append(t.prog, code)
			if e = binary.Read(rd, binary.LittleEndian, &v32); e != nil {
				break
			}
			p = unsafe.Pointer(&code[0])
			(*(*int32)(p)) = int32(v32)
			//fmt.Printf("%3d \t[%d]\n", pc, v32) //XXX
			//pc++                                //XXX
			t.prog = append(t.prog, code)

			if v16 != 0 {
				break
			}
			if e = binary.Read(rd, binary.LittleEndian, &v32); e != nil {
				break
			}
			p = unsafe.Pointer(&code[0])
			(*(*int32)(p)) = int32(v32)
			//fmt.Printf("%3d \t[%d]\n", pc, v32) //XXX
			//pc++                                //XXX
			t.prog = append(t.prog, code)
		default:
			e = fmt.Errorf("invalid format: undefined operation error")
			break
		}
	}
	if e == io.EOF {
		e = nil
	}
	return
}
