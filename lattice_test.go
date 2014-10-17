package kagome

import (
	"testing"
)

func TestLatticeBuild01(t *testing.T) {
	la := newLattice()
	if la == nil {
		t.Error("cannot new a lattice\n")
	}
	inp := ""
	la.build(inp)
	if la.input != inp {
		t.Errorf("got %v, expected %v\n", la.input, inp)
	}
	boseos := node{id: -1}
	if len(la.list) != 2 {
		t.Errorf("lattice initialize error: got %v, expected has 2 eos/bos nodes\n", la.list)
	} else if len(la.list[0]) != 1 || *la.list[0][0] != boseos {
		t.Errorf("lattice initialize error: got %v, expected %v\n", *la.list[0][0], boseos)
	} else if len(la.list[1]) != 1 || *la.list[1][0] != boseos {
		t.Errorf("lattice initialize error: got %v, expected %v\n", *la.list[1][0], boseos)
	}
	if len(la.output) != 0 {
		t.Errorf("lattice initialize error: got %v, expected empty\n", la.output)
	}
	if la.pool != nil {
		t.Errorf("lattice initialize error: node pool is not nil: %v\n", la.pool)
	}

	if la.dic == nil {
		t.Errorf("lattice initialize error: dic is nil\n", la.udic)
	}
	if la.udic != nil {
		t.Errorf("lattice initialize error: got %v, expected empty\n", la.udic)
	}
}

func TestLatticeBuild02(t *testing.T) {
	la := newLattice()
	if la == nil {
		t.Error("cannot new a lattice\n")
	}
	inp := "あ"
	la.build(inp)
	if la.input != inp {
		t.Errorf("got %v, expected %v\n", la.input, inp)
	}
	bos := node{id: -1}
	eos := node{id: -1, start: 1}
	if len(la.list) != 3 {
		t.Errorf("lattice initialize error: got %v, expected has 2 eos/bos nodes\n", la.list)
	} else if len(la.list[0]) != 1 || *la.list[0][0] != bos {
		t.Errorf("lattice initialize error: got %v, expected %v\n", *la.list[0][0], bos)
	} else if len(la.list[2]) != 1 || *la.list[2][0] != eos {
		t.Errorf("lattice initialize error: got %v, expected %v\n", *la.list[2][0], eos)
	}

	expected := 4
	if len(la.list[1]) != expected {
		t.Errorf("lattice initialize error: got %v, expected %v\n", len(la.list[1]), expected)
	} else {
		l := la.list[1]
		callAndResponse := []struct {
			in  int
			out node
		}{
			{in: 0, out: node{122, 0, 1, 0, 3, 3, 5549, "あ", nil}},
			{in: 1, out: node{123, 0, 1, 0, 776, 776, 6690, "あ", nil}},
			{in: 2, out: node{124, 0, 1, 0, 2, 2, 4262, "あ", nil}},
			{in: 3, out: node{125, 0, 1, 0, 1118, 1118, 9035, "あ", nil}},
		}
		for _, cr := range callAndResponse {
			if *l[cr.in] != cr.out {
				t.Errorf("lattice initialize error: got %v, expected %v\n", l[cr.in], cr.out)
			}
		}
	}
	if len(la.output) != 0 {
		t.Errorf("lattice initialize error: got %v, expected empty\n", la.output)
	}
	if la.pool != nil {
		t.Errorf("lattice initialize error: node pool is not nil: %v\n", la.pool)
	}
	if la.dic == nil {
		t.Errorf("lattice initialize error: dic is nil\n", la.udic)
	}
	if la.udic != nil {
		t.Errorf("lattice initialize error: got %v, expected empty\n", la.udic)
	}
}

func TestSetDic01(t *testing.T) {
	la := newLattice()
	la.setDic(nil)
	if la.dic == nil {
		t.Error("dic is nil\n")
	}
	d := NewSysDic()
	if la.dic != d {
		t.Errorf("got %p, expected %p\n", la.dic, d)
	}
}

func TestSetUserDic01(t *testing.T) {
	la := newLattice()
	udic, e := NewUserDic("_sample/userdic.txt")
	if e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	la.setUserDic(udic)
	if la.udic != udic {
		t.Error("got %p, expected %p", la.udic, udic)
	}
	la.setUserDic(nil)
	if la.udic != nil {
		t.Errorf("got %p, expected nil\n", la.dic)
	}
}

func TestSetNodePool01(t *testing.T) {
	la := newLattice()
	la.setNodePool(0)
	if la.pool == nil {
		t.Error("lattice initialize error: node pool is nil\n")
	}
	la.build("")
	la.build("すもももももももものうち")
}

func TestKanjiOnly01(t *testing.T) {
	callAndResponse := []struct {
		in  string
		out bool
	}{
		{in: "ひらがな", out: false},
		{in: "カタカナ", out: false},
		{in: "漢字", out: true},
		{in: "かな漢字交じり", out: false},
		{in: "123", out: false},
		{in: "#$%", out: false},
		{in: "", out: false},
	}
	for _, cr := range callAndResponse {
		if rsp := kanjiOnly(cr.in); rsp != cr.out {
			t.Errorf("in: %v, got %v, expected: %v", cr.in, rsp, cr.out)
		}
	}
}
