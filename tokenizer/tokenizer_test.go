//  Copyright (c) 2015 ikawaha.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package tokenizer

import (
	"testing"

	"github.com/ikawaha/kagome/internal/lattice"
)

func TestTokenize01(t *testing.T) {
	tnz := New(SysDic())
	tokens := tnz.Tokenize("", Normal)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: -1, Surface: "EOS"},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestTokenize02(t *testing.T) {
	tnz := New(SysDic())
	tokens := tnz.Tokenize("関西国際空港", Normal)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 372977, Surface: "関西国際空港", Start: 0, End: 6, Class: TokenClass(lattice.KNOWN)},
		{ID: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestTokenize03(t *testing.T) {
	tnz := New(SysDic())
	udic, e := NewUserDic("../_sample/userdic.txt")
	if e != nil {
		t.Fatalf("new user dic: unexpected error\n")
	}
	tnz.SetUserDic(udic)
	tokens := tnz.Tokenize("関西国際空港", Normal)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 2, Surface: "関西国際空港", Start: 0, End: 6, Class: TokenClass(lattice.USER)},
		{ID: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestTokenize04(t *testing.T) {
	tnz := New(SysDic())
	tokens := tnz.Tokenize("ポポピ", Normal)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 34, Surface: "ポポピ", Start: 0, End: 3, Class: TokenClass(lattice.UNKNOWN)},
		{ID: -1, Surface: "EOS", Start: 3, End: 3},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestSearcModeTokenize01(t *testing.T) {
	tnz := New(SysDic())
	tokens := tnz.Tokenize("", Search)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: -1, Surface: "EOS"},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestSearchModeTokenize02(t *testing.T) {
	tnz := New(SysDic())
	tokens := tnz.Tokenize("関西国際空港", Search)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 372968, Surface: "関西", Start: 0, End: 2, Class: TokenClass(lattice.KNOWN)},
		{ID: 168541, Surface: "国際", Start: 2, End: 4, Class: TokenClass(lattice.KNOWN)},
		{ID: 307133, Surface: "空港", Start: 4, End: 6, Class: TokenClass(lattice.KNOWN)},
		{ID: -1, Surface: "EOS", Start: 6, End: 6},
	}

	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestSearchModeTokenize03(t *testing.T) {
	tnz := New(SysDic())
	udic, e := NewUserDic("../_sample/userdic.txt")
	if e != nil {
		t.Fatalf("new user dic: unexpected error\n")
	}
	tnz.SetUserDic(udic)
	tokens := tnz.Tokenize("関西国際空港", Search)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 2, Surface: "関西国際空港", Start: 0, End: 6, Class: TokenClass(lattice.USER)},
		{ID: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestSearchModeTokenize04(t *testing.T) {
	tnz := New(SysDic())
	tokens := tnz.Tokenize("ポポピ", Search)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 34, Surface: "ポポピ", Start: 0, End: 3, Class: TokenClass(lattice.UNKNOWN)},
		{ID: -1, Surface: "EOS", Start: 3, End: 3},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestExtendedModeTokenize01(t *testing.T) {
	tnz := New(SysDic())
	tokens := tnz.Tokenize("", Extended)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: -1, Surface: "EOS"},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestExtendedModeTokenize02(t *testing.T) {
	tnz := New(SysDic())
	tokens := tnz.Tokenize("関西国際空港", Extended)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 372968, Surface: "関西", Start: 0, End: 2, Class: TokenClass(lattice.KNOWN)},
		{ID: 168541, Surface: "国際", Start: 2, End: 4, Class: TokenClass(lattice.KNOWN)},
		{ID: 307133, Surface: "空港", Start: 4, End: 6, Class: TokenClass(lattice.KNOWN)},
		{ID: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestExtendedModeTokenize03(t *testing.T) {
	tnz := New(SysDic())
	udic, e := NewUserDic("../_sample/userdic.txt")
	if e != nil {
		t.Fatalf("new user dic: unexpected error\n")
	}
	tnz.SetUserDic(udic)
	tokens := tnz.Tokenize("関西国際空港", Extended)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 2, Surface: "関西国際空港", Start: 0, End: 6, Class: TokenClass(lattice.USER)},
		{ID: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestExtendedModeTokenize04(t *testing.T) {
	tnz := New(SysDic())
	tokens := tnz.Tokenize("ポポピ", Extended)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 34, Surface: "ポ", Start: 0, End: 1, Class: TokenClass(lattice.DUMMY)},
		{ID: 34, Surface: "ポ", Start: 1, End: 2, Class: TokenClass(lattice.DUMMY)},
		{ID: 34, Surface: "ピ", Start: 2, End: 3, Class: TokenClass(lattice.DUMMY)},
		{ID: -1, Surface: "EOS", Start: 3, End: 3},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func BenchmarkTokenize01(b *testing.B) {
	tnz := New(SysDic())
	for i := 0; i < b.N; i++ {
		tnz.Tokenize("人魚は、南の方の海にばかり棲んでいるのではありません。北の海にも棲んでいたのであります。北方の海の色は、青うございました。ある時、岩の上に、女の人魚があがって、あたりの景色を眺めながら休んでいました。", Normal)
	}
}

func BenchmarkTokenize02(b *testing.B) {
	tnz := New(SysDic())
	for i := 0; i < b.N; i++ {
		tnz.Tokenize("人魚は、南の方の海にばかり棲んでいるのではありません。北の海にも棲んでいたのであります。北方の海の色は、青うございました。ある時、岩の上に、女の人魚があがって、あたりの景色を眺めながら休んでいました。", Search)
	}
}

func BenchmarkTokenize03(b *testing.B) {
	tnz := New(SysDic())
	for i := 0; i < b.N; i++ {
		tnz.Tokenize("人魚は、南の方の海にばかり棲んでいるのではありません。北の海にも棲んでいたのであります。北方の海の色は、青うございました。ある時、岩の上に、女の人魚があがって、あたりの景色を眺めながら休んでいました。", Extended)
	}
}
