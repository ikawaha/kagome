package lattice

import (
	"bytes"
	"testing"
	"unicode/utf8"

	ipa "github.com/ikawaha/kagome-dict-ipa"
	"github.com/ikawaha/kagome/v2/dict"
)

func Test_LatticeBuildEmptyInput(t *testing.T) {
	la := New(ipa.Dict(), nil)
	if la == nil {
		t.Fatal("lattice new failed")
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

func Test_LatticeBuild(t *testing.T) {
	la := New(ipa.Dict(), nil)
	if la == nil {
		t.Fatal("lattice new failed")
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
		t.Fatalf("lattice initialize error: got %v, expected %v", len(la.list[1]), expected)
	}
	l := la.list[1]
	for _, v := range l {
		if v.Surface != inp {
			t.Errorf("lattice initialize error: got %+v, expected surface %s", v, inp)
		}
		if v.Class != KNOWN {
			t.Errorf("lattice initialize error: got %+v, expected class KNOWN", v)
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

func Test_LatticeBuildWithUserDict(t *testing.T) {

	const udictPath = "../../_sample/userdict.txt"

	udic, err := dict.NewUserDict(udictPath)
	if err != nil {
		t.Fatalf("unexpected error: cannot load user dic, %v", err)
	}
	la := New(ipa.Dict(), udic)
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

func Test_LatticeBuildUnknown(t *testing.T) {
	la := New(ipa.Dict(), nil)
	if la == nil {
		t.Fatal("lattice new failed")
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
		t.Fatalf("lattice initialize error: got %v, expected %v", len(la.list[1]), expected)
	}
	l := la.list[1]
	var known, unknown, undef int
	for _, v := range l {
		if v.Surface != string([]rune(inp)[0:1]) {
			t.Errorf("lattice initialize error: got %+v, expected surface %c", v, []rune(inp)[0])
		}
		switch v.Class {
		case KNOWN:
			known++
		case UNKNOWN:
			unknown++
		default:
			undef++
		}
	}
	if known != 1 {
		t.Errorf("lattice initialize error: got KNOWN %d, expected 1, %+v", known, l)
	}
	if unknown != 6 {
		t.Errorf("lattice initialize error: got UNKNOWN %d, expected 6, %+v", unknown, l)
	}
	if undef != 0 {
		t.Errorf("lattice initialize error: got unexpected class %d, %+v", undef, l)
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

func Test_LatticeBuildMaximumUnknownWordLength(t *testing.T) {

	la := New(ipa.Dict(), nil)
	if la == nil {
		t.Fatal("cannot new a lattice")
	}
	defer la.Free()

	inp := "ポポピポンポコナーノ"
	var b bytes.Buffer
	for i, step := 0, utf8.RuneCountInString(inp); i < maximumUnknownWordLength; i = i + step {
		if _, err := b.WriteString(inp); err != nil {
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

func Test_LatticeBuildInvalidInput(t *testing.T) {

	la := New(ipa.Dict(), nil)
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

func Test_KanjiOnly(t *testing.T) {
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

func Test_LatticeString(t *testing.T) {
	la := New(ipa.Dict(), nil)
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

func Test_LatticeDot(t *testing.T) {
	la := New(ipa.Dict(), nil)
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

func Test_LatticeNewAndFree(t *testing.T) {
	for i := 0; i < 100; i++ {
		la := New(ipa.Dict(), nil)
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
		la = New(ipa.Dict(), nil)
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

func Test_Forward(t *testing.T) {
	la := New(ipa.Dict(), nil)
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

func Test_Backward(t *testing.T) {
	la := New(ipa.Dict(), nil)
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
