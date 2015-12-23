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

package tokenizer

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/ikawaha/kagome/internal/lattice"
)

const (
	testUserDicPath = "../_sample/userdic.txt"
)

func TestAnalyze01(t *testing.T) {
	tnz := New()
	tokens := tnz.Analyze("", Normal)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: -1, Surface: "EOS"},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestAnalyze02(t *testing.T) {
	tnz := New()
	tokens := tnz.Analyze("関西国際空港", Normal)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 372977, Surface: "関西国際空港", Start: 0, End: 6, Class: TokenClass(lattice.KNOWN)},
		{ID: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestAnalyze03(t *testing.T) {
	tnz := New()
	udic, e := NewUserDic(testUserDicPath)
	if e != nil {
		t.Fatalf("new user dic: unexpected error")
	}
	tnz.SetUserDic(udic)
	tokens := tnz.Analyze("関西国際空港", Normal)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 2, Surface: "関西国際空港", Start: 0, End: 6, Class: TokenClass(lattice.USER)},
		{ID: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestAnalyze04(t *testing.T) {
	tnz := New()
	tokens := tnz.Analyze("ポポピ", Normal)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 34, Surface: "ポポピ", Start: 0, End: 3, Class: TokenClass(lattice.UNKNOWN)},
		{ID: -1, Surface: "EOS", Start: 3, End: 3},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestTokenize(t *testing.T) {
	const input = "すもももももももものうち"
	tnz := New()
	x := tnz.Tokenize(input)
	y := tnz.Analyze(input, Normal)
	if !reflect.DeepEqual(x, y) {
		t.Errorf("got %v, expected %v", x, y)
	}
}

func TestSearcModeAnalyze01(t *testing.T) {
	tnz := New()
	tokens := tnz.Analyze("", Search)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: -1, Surface: "EOS"},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestSearchModeAnalyze02(t *testing.T) {
	tnz := New()
	tokens := tnz.Analyze("関西国際空港", Search)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 372968, Surface: "関西", Start: 0, End: 2, Class: TokenClass(lattice.KNOWN)},
		{ID: 168541, Surface: "国際", Start: 2, End: 4, Class: TokenClass(lattice.KNOWN)},
		{ID: 307133, Surface: "空港", Start: 4, End: 6, Class: TokenClass(lattice.KNOWN)},
		{ID: -1, Surface: "EOS", Start: 6, End: 6},
	}

	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestSearchModeAnalyze03(t *testing.T) {
	tnz := New()
	udic, e := NewUserDic(testUserDicPath)
	if e != nil {
		t.Fatalf("new user dic: unexpected error")
	}
	tnz.SetUserDic(udic)
	tokens := tnz.Analyze("関西国際空港", Search)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 2, Surface: "関西国際空港", Start: 0, End: 6, Class: TokenClass(lattice.USER)},
		{ID: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestSearchModeAnalyze04(t *testing.T) {
	tnz := New()
	tokens := tnz.Analyze("ポポピ", Search)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 34, Surface: "ポポピ", Start: 0, End: 3, Class: TokenClass(lattice.UNKNOWN)},
		{ID: -1, Surface: "EOS", Start: 3, End: 3},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestExtendedModeAnalyze01(t *testing.T) {
	tnz := New()
	tokens := tnz.Analyze("", Extended)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: -1, Surface: "EOS"},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestExtendedModeAnalyze02(t *testing.T) {
	tnz := New()
	tokens := tnz.Analyze("関西国際空港", Extended)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 372968, Surface: "関西", Start: 0, End: 2, Class: TokenClass(lattice.KNOWN)},
		{ID: 168541, Surface: "国際", Start: 2, End: 4, Class: TokenClass(lattice.KNOWN)},
		{ID: 307133, Surface: "空港", Start: 4, End: 6, Class: TokenClass(lattice.KNOWN)},
		{ID: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestExtendedModeAnalyze03(t *testing.T) {
	tnz := New()
	udic, e := NewUserDic(testUserDicPath)
	if e != nil {
		t.Fatalf("new user dic: unexpected error")
	}
	tnz.SetUserDic(udic)
	tokens := tnz.Analyze("関西国際空港", Extended)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 2, Surface: "関西国際空港", Start: 0, End: 6, Class: TokenClass(lattice.USER)},
		{ID: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestExtendedModeAnalyze04(t *testing.T) {
	tnz := New()
	tokens := tnz.Analyze("ポポピ", Extended)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 34, Surface: "ポ", Start: 0, End: 1, Class: TokenClass(lattice.DUMMY)},
		{ID: 34, Surface: "ポ", Start: 1, End: 2, Class: TokenClass(lattice.DUMMY)},
		{ID: 34, Surface: "ピ", Start: 2, End: 3, Class: TokenClass(lattice.DUMMY)},
		{ID: -1, Surface: "EOS", Start: 3, End: 3},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestTokenizerSetDic(t *testing.T) {
	d := SysDic()
	tnz := NewWithDic(d)

	tnz.SetDic(d)
	if tnz.dic != d.dic {
		t.Errorf("got %v, expected %v", tnz.dic, d)
	}
}

func TestTokenizerDot(t *testing.T) {
	tnz := New()

	// test empty case
	var b bytes.Buffer
	tnz.Dot("", &b)
	if b.String() == "" {
		t.Errorf("got empty string")
	}

	// only idling
	b.Reset()
	tnz.Dot("わたしまけましたわ", &b)
	if b.String() == "" {
		t.Errorf("got empty string")
	}
}

func TestTokenizerAnalyzeGraph(t *testing.T) {
	tnz := New()

	// test empty case
	for _, mode := range []TokenizeMode{Normal, Search, Extended} {
		var b bytes.Buffer
		tnz.AnalyzeGraph("", mode, &b)
		if b.String() == "" {
			t.Errorf("got empty string")
		}

		// only idling
		b.Reset()
		tnz.Dot("わたしまけましたわ", &b)
		if b.String() == "" {
			t.Errorf("got empty string")
		}
	}
}

var benchSampleText = "人魚は、南の方の海にばかり棲んでいるのではありません。北の海にも棲んでいたのであります。北方の海の色は、青うございました。ある時、岩の上に、女の人魚があがって、あたりの景色を眺めながら休んでいました。"

func BenchmarkAnalyzeNormal(b *testing.B) {
	tnz := New()
	for i := 0; i < b.N; i++ {
		tnz.Analyze(benchSampleText, Normal)
	}
}

func BenchmarkAnalyzeSearch(b *testing.B) {
	tnz := New()
	for i := 0; i < b.N; i++ {
		tnz.Analyze(benchSampleText, Search)
	}
}

func BenchmarkAnalyzeExtended(b *testing.B) {
	tnz := New()
	for i := 0; i < b.N; i++ {
		tnz.Analyze(benchSampleText, Extended)
	}
}
