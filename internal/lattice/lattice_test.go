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
	"testing"

	"github.com/ikawaha/kagome/internal/dic"
)

func TestLatticeBuild01(t *testing.T) {
	la := New(dic.SysDic(), nil)
	if la == nil {
		t.Error("cannot new a lattice\n")
	}
	inp := ""
	la.Build(inp)
	if la.Input != inp {
		t.Errorf("got %v, expected %v\n", la.Input, inp)
	}
	boseos := node{ID: -1}
	if len(la.list) != 2 {
		t.Errorf("lattice initialize error: got %v, expected has 2 eos/bos nodes\n", la.list)
	} else if len(la.list[0]) != 1 || *la.list[0][0] != boseos {
		t.Errorf("lattice initialize error: got %v, expected %v\n", *la.list[0][0], boseos)
	} else if len(la.list[1]) != 1 || *la.list[1][0] != boseos {
		t.Errorf("lattice initialize error: got %v, expected %v\n", *la.list[1][0], boseos)
	}
	if len(la.Output) != 0 {
		t.Errorf("lattice initialize error: got %v, expected empty\n", la.Output)
	}

	if la.dic == nil {
		t.Errorf("lattice initialize error: dic is nil\n")
	}
	if la.udic != nil {
		t.Errorf("lattice initialize error: got %v, expected empty\n", la.udic)
	}
}

func TestLatticeBuild02(t *testing.T) {
	la := New(dic.SysDic(), nil)
	if la == nil {
		t.Fatal("cannot new a lattice\n")
	}
	inp := "あ"
	la.Build(inp)
	if la.Input != inp {
		t.Errorf("got %v, expected %v\n", la.Input, inp)
	}
	bos := node{ID: -1}
	eos := node{ID: -1, Start: 1}
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
			{in: 0, out: node{122, 0, KNOWN, 0, 3, 3, 5549, "あ", nil}},
			{in: 1, out: node{123, 0, KNOWN, 0, 776, 776, 6690, "あ", nil}},
			{in: 2, out: node{124, 0, KNOWN, 0, 2, 2, 4262, "あ", nil}},
			{in: 3, out: node{125, 0, KNOWN, 0, 1118, 1118, 9035, "あ", nil}},
		}
		for _, cr := range callAndResponse {
			if *l[cr.in] != cr.out {
				t.Errorf("lattice initialize error: got %v, expected %v\n", l[cr.in], cr.out)
			}
		}
	}
	if len(la.Output) != 0 {
		t.Errorf("lattice initialize error: got %v, expected empty\n", la.Output)
	}
	if la.dic == nil {
		t.Errorf("lattice initialize error: dic is nil\n")
	}
	if la.udic != nil {
		t.Errorf("lattice initialize error: got %v, expected empty\n", la.udic)
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
