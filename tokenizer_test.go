package kagome

import (
	"testing"
)

func TestTokenize01(t *testing.T) {
	tnz := NewTokenizer()
	tokens := tnz.Tokenize("")
	expected := []Token{
		{Id: -1, Surface: "BOS"},
		{Id: -1, Surface: "EOS"},
	}
	if len(tokens) != len(expected) {
		t.Errorf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Id != expected[i].Id ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestTokenize02(t *testing.T) {
	tnz := NewTokenizer()
	tokens := tnz.Tokenize("関西国際空港")
	expected := []Token{
		{Id: -1, Surface: "BOS"},
		{Id: 372977, Surface: "関西国際空港", Start: 0, End: 6, Class: KNOWN},
		{Id: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Errorf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Id != expected[i].Id ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestTokenize03(t *testing.T) {
	tnz := NewTokenizer()
	udic, e := NewUserDic("_sample/userdic.txt")
	if e != nil {
		t.Fatalf("new user dic: unexpected error\n")
	}
	tnz.SetUserDic(udic)
	tokens := tnz.Tokenize("関西国際空港")
	expected := []Token{
		{Id: -1, Surface: "BOS"},
		{Id: 2, Surface: "関西国際空港", Start: 0, End: 6, Class: USER},
		{Id: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Errorf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Id != expected[i].Id ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestTokenize04(t *testing.T) {
	tnz := NewTokenizer()
	tokens := tnz.Tokenize("ポポピ")
	expected := []Token{
		{Id: -1, Surface: "BOS"},
		{Id: 34, Surface: "ポポピ", Start: 0, End: 3, Class: UNKNOWN},
		{Id: -1, Surface: "EOS", Start: 3, End: 3},
	}
	if len(tokens) != len(expected) {
		t.Errorf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Id != expected[i].Id ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}
