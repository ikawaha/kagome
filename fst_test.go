package kagome

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestStateEq01(t *testing.T) {
	type pair struct {
		x *state
		y *state
	}

	s := &state{}

	crs := []struct {
		call pair
		resp bool
	}{
		{pair{x: s, y: s}, true},
		{pair{x: nil, y: nil}, false},
		{pair{x: nil, y: &state{}}, false},
		{pair{x: &state{}, y: nil}, false},
		{pair{&state{ID: 1}, &state{ID: 2}}, true},
		{pair{&state{IsFinal: true}, &state{IsFinal: false}}, false},
		{pair{&state{Output: map[byte]int32{1: 555}}, &state{}}, false},
		{pair{&state{Output: map[byte]int32{1: 555}}, &state{Output: map[byte]int32{1: 555}}},
			true},
		{pair{&state{Output: map[byte]int32{1: 555}}, &state{Output: map[byte]int32{1: 444}}},
			false},
		{pair{&state{Output: map[byte]int32{1: 555}}, &state{Output: map[byte]int32{2: 555}}},
			false},
		{pair{&state{Tail: map[int32]bool{555: true}}, &state{Tail: map[int32]bool{555: true}}}, true},
	}
	for _, cr := range crs {
		if rst := cr.call.x.eq(cr.call.y); rst != cr.resp {
			t.Errorf("got %v, expected %v, %v\n", rst, cr.resp, cr)
		}
		if rst := cr.call.y.eq(cr.call.x); rst != cr.resp {
			t.Errorf("got %v, expected %v, %v\n", rst, cr.resp, cr)
		}
	}
}

func TestStateEq02(t *testing.T) {
	x := &state{ID: 1}
	y := &state{ID: 2}
	a := &state{
		Trans: map[byte]*state{1: x, 2: y},
	}
	b := &state{
		Trans: map[byte]*state{1: x, 2: y},
	}
	c := &state{
		Trans: map[byte]*state{1: y, 2: y},
	}
	d := &state{
		Trans: map[byte]*state{1: x, 2: y, 3: x},
	}
	if rst, exp := a.eq(b), true; rst != exp {
		t.Errorf("got %v, expected %v\n", rst, exp)
	}
	if rst, exp := a.eq(c), false; rst != exp {
		t.Errorf("got %v, expected %v\n", rst, exp)
	}
	if rst, exp := a.eq(c), false; rst != exp {
		t.Errorf("got %v, expected %v\n", rst, exp)
	}
	if rst, exp := a.eq(d), false; rst != exp {
		t.Errorf("got %v, expected %v\n", rst, exp)
	}

}

func TestStateString01(t *testing.T) {
	crs := []struct {
		call *state
		resp string
	}{
		{nil, "<nil>"},
	}
	for _, cr := range crs {
		if rst := cr.call.String(); rst != cr.resp {
			t.Errorf("got %v, expected %v, %v\n", rst, cr.resp, cr)
		}
	}
	r := &state{}
	s := state{
		ID:      1,
		Trans:   map[byte]*state{1: nil, 2: r},
		Output:  map[byte]int32{3: 555, 4: 888},
		Tail:    int32Set{1111: true},
		IsFinal: true,
		//		Prev:    []*state{nil, r},
	}
	fmt.Println(s.String())
}

func TestMASTBuildMAST01(t *testing.T) {
	inp := PairSlice{}
	m := buildMAST(inp)
	if m.initialState.ID != 0 {
		t.Errorf("got initial state id %v, expected 0\n", m.initialState.ID)
	}
	if len(m.states) != 1 {
		t.Errorf("expected: initial state only, got %v\n", m.states)
	}
	if len(m.finalStates) != 0 {
		t.Errorf("expected: final state is empty, got %v\n", m.finalStates)
	}
}

func TestMASTAccept01(t *testing.T) {
	inp := PairSlice{
		{"hello", 111},
		{"hello", 222},
		{"111", 111},
		{"112", 112},
		{"112", 122},
		{"211", 345},
	}
	m := buildMAST(inp)
	for _, pair := range inp {
		if ok := m.accept(pair.In); !ok {
			t.Errorf("expected: accept [%v]\n", pair.In)
		}
	}
	if ok := m.accept("aloha"); ok {
		t.Errorf("expected: reject \"aloha\"\n")
	}
}

func TestMASTRun01(t *testing.T) {
	inp := PairSlice{
		{"hello", 1111},
		{"hell", 2222},
		{"111", 111},
		{"112", 112},
		{"113", 122},
		{"211", 111},
	}
	m := buildMAST(inp)
	for _, pair := range inp {
		out, ok := m.run(pair.In)
		if !ok {
			t.Errorf("expected: accept [%v]\n", pair.In)
		}
		if len(out) != 1 {
			t.Errorf("input: %v, output size: got %v, expected 1\n", pair.In, len(out))
		}
		if out[0] != pair.Out {
			t.Errorf("input: %v, output: got %v, expected %v\n", pair.In, pair.Out, out[0])
		}
	}
	if out, ok := m.run("aloha"); ok {
		t.Errorf("expected: reject \"aloha\", %v\n", out)
	}
}

func TestMASTRun02(t *testing.T) {
	inp := PairSlice{
		{"hello", 1111},
		{"hello", 2222},
	}
	m := buildMAST(inp)
	for _, pair := range inp {
		out, ok := m.run(pair.In)
		if !ok {
			t.Errorf("expected: accept [%v]\n", pair.In)
		}
		if len(out) != 2 {
			t.Errorf("input: %v, output size: got %v, expected 2\n", pair.In, len(out))
		}
		expected := []int32{1111, 2222}
		sort.Sort(int32Slice(out))
		sort.Sort(int32Slice(expected))
		if !reflect.DeepEqual(out, expected) {
			t.Errorf("input: %v, output: got %v, expected %v\n", pair.In, out, expected)
		}
	}
	if out, ok := m.run("aloha"); ok {
		t.Errorf("expected: reject \"aloha\", %v\n", out)
	}
}

func TestMASTDot01(t *testing.T) {
	inp := PairSlice{
		{"apr", 30},
		{"aug", 31},
		{"dec", 31},
		{"feb", 28},
		{"feb", 29},
	}
	m := buildMAST(inp)
	m.dot(os.Stdout)
}

func TestMASTDot02(t *testing.T) {
	inp := PairSlice{
		{"apr", 30},
		{"aug", 31},
		{"dec", 31},
		{"feb", 28},
		{"feb", 29},
		{"lucene", 1},
		{"lucid", 2},
		{"lucifer", 666},
	}
	m := buildMAST(inp)
	m.dot(os.Stdout)
}

func TestFSTRun01(t *testing.T) {
	inp := PairSlice{
		{"feb", 28},
		{"feb", 29},
		{"feb", 30},
		{"dec", 31},
	}
	m := buildMAST(inp)
	m.dot(os.Stdout)

	fst, _ := m.buildMachine()
	fmt.Println(fst)

	config, ok := fst.run("feb")
	if !ok {
		t.Errorf("input:feb, config:%v, accept:%v", config, ok)
	}
	fmt.Println(config)

}

func TestFSTRun02(t *testing.T) {
	inp := PairSlice{
		{"feb", 28},
		{"feb", 29},
		{"feb", 30},
		{"dec", 31},
	}
	m := buildMAST(inp)
	m.dot(os.Stdout)

	fst, _ := m.buildMachine()
	fmt.Println(fst)

	config, ok := fst.run("dec")
	if !ok {
		t.Errorf("input:feb, config:%v, accept:%v", config, ok)
	}
	fmt.Println(config)

}

func TestFSTSearch01(t *testing.T) {
	inp := PairSlice{
		{"1a22xss", 111},
		{"1b22yss", 222},
	}
	vm, e := BuildFST(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)
	for _, p := range inp {
		outs := vm.Search(p.In)
		if !reflect.DeepEqual(outs, []int32{p.Out}) {
			t.Errorf("input: %v, got %v, expected %v\n", p.In, outs, []int32{p.Out})
		}
	}
}

func TestFSTSearch02(t *testing.T) {
	inp := PairSlice{
		{"1a22", 111},
		{"1a22xss", 222},
		{"1a22yss", 333},
	}
	vm, e := BuildFST(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)
	for _, p := range inp {
		outs := vm.Search(p.In)
		if !reflect.DeepEqual(outs, []int32{p.Out}) {
			t.Errorf("input: %v, got %v, expected %v\n", p.In, outs, []int32{p.Out})
		}
	}
}

func TestFSTSearch03(t *testing.T) {
	inp := PairSlice{
		{"1a22", 111},
		{"1a22xss", 222},
		{"1a22xss", 333},
	}
	vm, e := BuildFST(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)

	in := "1a22"
	exp := []int32{111}
	outs := vm.Search(in)
	if !reflect.DeepEqual(outs, []int32{111}) {
		t.Errorf("input: %v, got %v, expected %v\n", in, outs, exp)
	}
}

func TestFSTSearch04(t *testing.T) {
	inp := PairSlice{
		{"1a22", 111},
		{"1a22xss", 222},
		{"1a22xss", 333},
	}
	vm, e := BuildFST(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)

	in := "1a22xss"
	exp := []int32{222, 333}
	outs := vm.Search(in)
	sort.Sort(int32Slice(outs))
	if !reflect.DeepEqual(outs, exp) {
		t.Errorf("input: %v, got %v, expected %v\n", in, outs, exp)
	}
}

func TestFSTSearch05(t *testing.T) {
	inp := PairSlice{
		{"1a22", 0},
		{"1a22xss", 0},
		{"1a22xss", 0},
	}
	vm, e := BuildFST(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)

	in := "1a22xss"
	exp := []int32{0}
	outs := vm.Search(in)
	if !reflect.DeepEqual(outs, exp) {
		t.Errorf("input: %v, got %v, expected %v\n", in, outs, exp)
	}
}

func TestFSTSearch06(t *testing.T) {
	inp := PairSlice{
		{"こんにちは", 111},
		{"世界", 222},
		{"すもももももも", 333},
		{"すもも", 333},
		{"すもも", 444},
	}
	vm, e := BuildFST(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)

	cr := []struct {
		in  string
		out []int32
	}{
		{"すもも", []int32{333, 444}},
		{"こんにちわ", nil},
		{"こんにちは", []int32{111}},
		{"世界", []int32{222}},
		{"すもももももも", []int32{333}},
		{"すももももももも", nil},
		{"すも", nil},
		{"すもう", nil},
	}

	for _, pair := range cr {
		outs := vm.Search(pair.in)
		sort.Sort(int32Slice(outs))
		sort.Sort(int32Slice(pair.out))
		if !reflect.DeepEqual(outs, pair.out) {
			t.Errorf("input:%v, got %v, expected %v\n", pair.in, outs, pair.out)
		}
	}
}

func TestFSTPrefixSearch01(t *testing.T) {
	inp := PairSlice{
		{"こんにちは", 111},
		{"世界", 222},
		{"すもももももも", 333},
		{"すもも", 333},
		{"すもも", 444},
	}
	vm, e := BuildFST(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)

	crs := []struct {
		in  string
		pos int
		out []int32
	}{
		{"すもも", 9, []int32{333, 444}},
		{"こんにちわ", -1, nil},
		{"こんにちは", 15, []int32{111}},
		{"世界", 6, []int32{222}},
		{"すもももももも", 21, []int32{333}},
		{"すもももももももものうち", 21, []int32{333}},
		{"すも", -1, nil},
		{"すもう", -1, nil},
	}

	for _, cr := range crs {
		pos, outs := vm.PrefixSearch(cr.in)
		sort.Sort(int32Slice(outs))
		sort.Sort(int32Slice(cr.out))
		if !reflect.DeepEqual(outs, cr.out) {
			t.Errorf("input:%v, got %v, expected %v\n", cr.in, outs, cr.out)
		}
		if pos != cr.pos {
			t.Errorf("input:%v, got %v, expected %v\n", cr.in, pos, cr.pos)
		}
	}
}

func TestFSTCommonPrefixSearch01(t *testing.T) {
	inp := PairSlice{
		{"こんにちは", 111},
		{"世界", 222},
		{"すもももももも", 333},
		{"すもも", 333},
		{"すもも", 444},
	}
	vm, e := BuildFST(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	fmt.Println(vm)

	crs := []struct {
		in   string
		lens []int
		outs [][]int32
	}{
		{"すもも", []int{9}, [][]int32{{333, 444}}},
		{"こんにちわ", nil, nil},
		{"こんにちは", []int{15}, [][]int32{{111}}},
		{"世界", []int{6}, [][]int32{{222}}},
		{"すもももももも", []int{9, 21}, [][]int32{{333, 444}, {333}}},
		{"すもももももももものうち", []int{9, 21}, [][]int32{{333, 444}, {333}}},
		{"すも", nil, nil},
		{"すもう", nil, nil},
	}

	for _, cr := range crs {
		lens, outs := vm.CommonPrefixSearch(cr.in)
		if !reflect.DeepEqual(lens, cr.lens) {
			t.Errorf("input:%v, got %v %v, expected %v %v\n", cr.in, lens, outs, cr.lens, cr.outs)
		}
		for i := range lens {
			sort.Sort(int32Slice(outs[i]))
			sort.Sort(int32Slice(cr.outs[i]))
			if !reflect.DeepEqual(outs[i], cr.outs[i]) {
				t.Errorf("input:%v, got %v %v, expected %v %v\n", cr.in, lens, outs, cr.lens, cr.outs)
			}
		}
	}
}

func TestFSTSaveAndLoad01(t *testing.T) {
	inp := PairSlice{
		{"feb", 28},
		{"feb", 29},
		{"apr", 30},
		{"jan", 31},
		{"jun", 30},
		{"jul", 31},
		{"dec", 31},
	}

	org, e := BuildFST(inp)
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}

	var b bytes.Buffer
	org.Write(&b)

	var rst FST
	e = rst.Read(&b)

	if !reflect.DeepEqual(org.data, rst.data) {
		t.Errorf("data:got %v, expected %v\n", rst.data, org.data)
	}
	if !reflect.DeepEqual(org.prog, rst.prog) {
		t.Errorf("prog:got %v, expected %v\n", rst.prog, org.prog)
	}
}
