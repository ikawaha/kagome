package tokenizer

import (
	"testing"

	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome/v2/tokenizer/lattice"
)

const (
	testUserDictPath = "../_sample/userdict.txt"
)

func Test_AnalyzeWithUserDict(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	udict, err := dict.NewUserDict(testUserDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d, UserDict(udict))
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
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

func Test_AnalyzeWithSearchModeWithUserDict(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	udict, err := dict.NewUserDict(testUserDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d, UserDict(udict))
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

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

func Test_AnalyzeWithExtendedModeWithUserDict(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	udict, err := dict.NewUserDict(testUserDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d, UserDict(udict))
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

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

func TestTokenizer_Analyze_OmitBOSEOS(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d, OmitBosEos())
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tokens := tnz.Analyze("関西国際空港", Normal)
	expected := []Token{
		{ID: 372978, Surface: "関西国際空港", Start: 0, End: 6, Class: TokenClass(lattice.KNOWN)},
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
