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

package lattice

import (
	"bytes"
	"testing"
	"unicode/utf8"

	"github.com/ikawaha/kagome/internal/dic"
)

func TestLatticeBuild01(t *testing.T) {
	la := New(dic.SysDic(), nil)
	if la == nil {
		t.Error("cannot new a lattice")
	}
	defer la.Free()

	inp := ""
	la.Build(inp)
	if la.Input != inp {
		t.Errorf("got %v, expected %v", la.Input, inp)
	}
	boseos := node{ID: -1}
	if len(la.list) != 2 {
		t.Errorf("lattice initialize error: got %v, expected has 2 eos/bos nodes", la.list)
	} else if len(la.list[0]) != 1 || *la.list[0][0] != boseos {
		t.Errorf("lattice initialize error: got %v, expected %v", *la.list[0][0], boseos)
	} else if len(la.list[1]) != 1 || *la.list[1][0] != boseos {
		t.Errorf("lattice initialize error: got %v, expected %v", *la.list[1][0], boseos)
	}
	if len(la.Output) != 0 {
		t.Errorf("lattice initialize error: got %v, expected empty", la.Output)
	}

	if la.dic == nil {
		t.Errorf("lattice initialize error: dic is nil")
	}
	if la.udic != nil {
		t.Errorf("lattice initialize error: got %v, expected empty", la.udic)
	}
}

func TestLatticeBuild02(t *testing.T) {
	la := New(dic.SysDic(), nil)
	if la == nil {
		t.Fatal("cannot new a lattice")
	}
	defer la.Free()

	inp := "あ"
	la.Build(inp)
	if la.Input != inp {
		t.Errorf("got %v, expected %v", la.Input, inp)
	}
	bos := node{ID: -1}
	eos := node{ID: -1, Start: 1}
	if len(la.list) != 3 {
		t.Errorf("lattice initialize error: got %v, expected has 2 eos/bos nodes", la.list)
	} else if len(la.list[0]) != 1 || *la.list[0][0] != bos {
		t.Errorf("lattice initialize error: got %v, expected %v", *la.list[0][0], bos)
	} else if len(la.list[2]) != 1 || *la.list[2][0] != eos {
		t.Errorf("lattice initialize error: got %v, expected %v", *la.list[2][0], eos)
	}

	expected := 4
	if len(la.list[1]) != expected {
		t.Errorf("lattice initialize error: got %v, expected %v", len(la.list[1]), expected)
	} else {
		l := la.list[1]
		callAndResponse := []struct {
			in  int
			out node
		}{
			{in: 0, out: node{122, 0, KNOWN, 0, 3, 3, 5549, "あ", nil}},
			{in: 1, out: node{123, 0, KNOWN, 0, 776, 776, 6690, "あ", nil}},
			{in: 2, out: node{124, 0, KNOWN, 0, 2, 2, 4262, "あ", nil}},
			{in: 3, out: node{125, 0, KNOWN, 0, 1118, 1118, 9035, "あ", nil}},
		}
		for _, cr := range callAndResponse {
			if *l[cr.in] != cr.out {
				t.Errorf("lattice initialize error: got %v, expected %v", l[cr.in], cr.out)
			}
		}
	}
	if len(la.Output) != 0 {
		t.Errorf("lattice initialize error: got %v, expected empty", la.Output)
	}
	if la.dic == nil {
		t.Errorf("lattice initialize error: dic is nil")
	}
	if la.udic != nil {
		t.Errorf("lattice initialize error: got %v, expected empty", la.udic)
	}
}

func TestLatticeBuild03(t *testing.T) {

	const udicPath = "../../_sample/userdic.txt"

	udic, e := dic.NewUserDic(udicPath)
	if e != nil {
		t.Fatalf("unexpected error: cannot load user dic, %v", e)
	}
	la := New(dic.SysDic(), udic)
	if la == nil {
		t.Fatal("cannot new a lattice")
	}
	defer la.Free()

	inp := "朝青龍"
	la.Build(inp)
	if la.Input != inp {
		t.Errorf("got %v, expected %v", la.Input, inp)
	}

	if la.list[3][0].Class != USER {
		t.Errorf("%+v", la)
	}

	if len(la.Output) != 0 {
		t.Errorf("lattice initialize error: got %v, expected empty", la.Output)
	}
	if la.dic == nil {
		t.Errorf("lattice initialize error: dic is nil")
	}
	if la.udic == nil {
		t.Errorf("lattice initialize error: got %v, expected not empty", la.udic)
	}
}

func TestLatticeBuild04(t *testing.T) {
	la := New(dic.SysDic(), nil)
	if la == nil {
		t.Fatal("cannot new a lattice")
	}
	defer la.Free()

	inp := "ポポピ"
	la.Build(inp)
	if la.Input != inp {
		t.Errorf("got %v, expected %v", la.Input, inp)
	}
	bos := node{ID: -1}
	eos := node{ID: -1, Start: 3}
	if len(la.list) != 5 {
		t.Errorf("lattice initialize error: got %v, expected has 2 eos/bos nodes", la.list)
	} else if len(la.list[0]) != 1 || *la.list[0][0] != bos {
		t.Errorf("lattice initialize error: got %v, expected %v", *la.list[0][0], bos)
	} else if len(la.list[len(la.list)-1]) != 1 || *la.list[len(la.list)-1][0] != eos {
		t.Errorf("lattice initialize error: got %v, expected %v", *la.list[len(la.list)-1][0], eos)
	}

	expected := 7
	if len(la.list[1]) != expected {
		t.Errorf("lattice initialize error: got %v, expected %v", len(la.list[1]), expected)
	} else {
		l := la.list[1]
		callAndResponse := []struct {
			in  int
			out node
		}{
			{in: 0, out: node{98477, 0, KNOWN, 0, 1285, 1285, 4279, "ポ", nil}},
			{in: 1, out: node{31, 0, UNKNOWN, 0, 1289, 1289, 13581, "ポ", nil}},
			{in: 2, out: node{32, 0, UNKNOWN, 0, 1285, 1285, 9461, "ポ", nil}},
			{in: 3, out: node{33, 0, UNKNOWN, 0, 1293, 1293, 13661, "ポ", nil}},
			{in: 4, out: node{34, 0, UNKNOWN, 0, 1292, 1292, 10922, "ポ", nil}},
			{in: 5, out: node{35, 0, UNKNOWN, 0, 1288, 1288, 10521, "ポ", nil}},
			{in: 6, out: node{36, 0, UNKNOWN, 0, 3, 3, 14138, "ポ", nil}},
		}
		for _, cr := range callAndResponse {
			if *l[cr.in] != cr.out {
				t.Errorf("lattice initialize error: got %v, expected %v", l[cr.in], cr.out)
			}
		}
	}
	if len(la.Output) != 0 {
		t.Errorf("lattice initialize error: got %v, expected empty", la.Output)
	}
	if la.dic == nil {
		t.Errorf("lattice initialize error: dic is nil")
	}
	if la.udic != nil {
		t.Errorf("lattice initialize error: got %v, expected empty", la.udic)
	}
}

func TestLatticeBuild05(t *testing.T) {

	la := New(dic.SysDic(), nil)
	if la == nil {
		t.Fatal("cannot new a lattice")
	}
	defer la.Free()

	inp := "ポポピポンポコナーノ"
	var b bytes.Buffer
	for i, step := 0, utf8.RuneCountInString(inp); i < maximumUnknownWordLength; i = i + step {
		if _, e := b.WriteString(inp); e != nil {
			t.Fatalf("unexpected error: create the test input, %v", b.String())
		}
	}
	la.Build(b.String())
	for i := range la.list {
		for j := range la.list[i] {
			l := utf8.RuneCountInString(la.list[i][j].Surface)
			if l > maximumUnknownWordLength {
				t.Errorf("too long unknown word, %v", l)
			}
		}
	}
}

func TestLatticeBuildInvalidInput(t *testing.T) {

	la := New(dic.SysDic(), nil)
	if la == nil {
		t.Fatal("cannot new a lattice")
	}
	defer la.Free()

	inp := "\x96\x7b\x93\xfa" // sjis encoding for '日本'
	la.Build(inp)
	if la.Input != inp {
		t.Errorf("got %v, expected %v", la.Input, inp)
	}
	bos := node{ID: -1}
	eos := node{ID: -1, Start: 4}
	if len(la.list) != 6 {
		t.Errorf("lattice initialize error: got %v, expected has 2 eos/bos nodes", la.list)
	} else if len(la.list[0]) != 1 || *la.list[0][0] != bos {
		t.Errorf("lattice initialize error: got %v, expected %v", *la.list[0][0], bos)
	} else if len(la.list[len(la.list)-1]) != 1 || *la.list[len(la.list)-1][0] != eos {
		t.Errorf("lattice initialize error: got %v, expected %v", *la.list[len(la.list)-1][0], eos)
	}
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

func TestLatticeString(t *testing.T) {
	la := New(dic.SysDic(), nil)
	if la == nil {
		t.Fatal("cannot new a lattice")
	}
	defer la.Free()

	expected := ""
	str := la.String()
	if str != expected {
		t.Errorf("got %v, expected: %v", str, expected)
	}

	la.Build("わたしまけましたわ")
	str = la.String()
	if str == "" {
		t.Errorf("got empty string")
	}
}

func TestLatticeDot(t *testing.T) {
	la := New(dic.SysDic(), nil)
	if la == nil {
		t.Fatal("cannot new a lattice")
	}
	defer la.Free()

	expected := `graph lattice {
dpi=48;
graph [style=filled, splines=true, overlap=false, fontsize=30, rankdir=LR]
edge [fontname=Helvetica, fontcolor=red, color="#606060"]
node [shape=box, style=filled, fillcolor="#e8e8f0", fontname=Helvetica]
}
`
	var b bytes.Buffer
	la.Dot(&b)
	if b.String() != expected {
		t.Errorf("got %v, expected: %v", b.String(), expected)
	}
	b.Reset()
	la.Build("わたしまけましたわポポピ")
	la.Forward(Normal)
	la.Backward(Normal)
	la.Dot(&b)
	if b.String() == "" {
		t.Errorf("got empty string")
	}
}

func TestLatticeNewAndFree(t *testing.T) {
	for i := 0; i < 100; i++ {
		la := New(dic.SysDic(), nil)
		if la == nil {
			t.Fatal("unexpected error: cannot new a lattice")
		}
		if la.Input != "" {
			t.Fatalf("unexpected error: lattice input initialize error, %+v", la.Input)
		}
		if len(la.Output) != 0 {
			t.Fatalf("unexpected error: lattice output initialize error, %+v", la.Output)
		}
		if len(la.list) != 0 {
			t.Fatalf("unexpected error: lattice list initialize error, %+v", la.list)
		}
		la.Build("すべては科学する心に宿るのだ")
		la.Free()

		// renew
		la = New(dic.SysDic(), nil)
		if la == nil {
			t.Fatal("unexpected error: cannot new a lattice")
		}
		if la.Input != "" {
			t.Fatalf("unexpected error: lattice input initialize error, %+v", la.Input)
		}
		if len(la.Output) != 0 {
			t.Fatalf("unexpected error: lattice output initialize error, %+v", la.Output)
		}
		if len(la.list) != 0 {
			t.Fatalf("unexpected error: lattice list initialize error, %+v", la.list)
		}
		la.Free()
	}
}

func TestForward(t *testing.T) {
	la := New(dic.SysDic(), nil)
	if la == nil {
		t.Fatal("unexpected error: cannot new a lattice")
	}

	la.Forward(Normal)
	la.Forward(Search)
	la.Forward(Extended)

	for _, m := range []TokenizeMode{Normal, Search, Extended} {
		la.Build("わたしまけましたわ．関西国際空港．ポポポポポポポポポポ．\U0001f363\U0001f363\U0001f363\U0001f363\U0001f363\U0001f363\U0001f363\U0001f363\U0001f363\U0001f363")
		la.Forward(m)
	}
}

func TestBackward01(t *testing.T) {
	la := New(dic.SysDic(), nil)
	if la == nil {
		t.Fatal("unexpected error: cannot new a lattice")
	}

	// only run
	la.Backward(Normal)
	la.Backward(Search)
	la.Backward(Extended)

	for _, m := range []TokenizeMode{Normal, Search, Extended} {
		la.Build("わたしまけましたわ．ポポピ")
		la.Forward(m)
		la.Backward(m)
	}
}
