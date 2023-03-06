package tokenizer

import (
	"testing"

	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome/v2/tokenizer/lattice"
)

const (
	testUserDictPath = "../sample/dict/userdict.txt"
)

func TestTokenizer_Analyze_Nop(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	_, err = New(d, Nop())
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
}

func Test_AnalyzeWithUserDict(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	t.Run("invalid nil user dict", func(t *testing.T) {
		_, err := New(d, UserDict(nil))
		if err == nil {
			t.Fatal("expected empty user dictionary error")
		} else if err.Error() != "empty user dictionary" {
			t.Errorf("expected empty user dictionary error, got %v", err)
		}
	})
	t.Run("user dict", func(t *testing.T) {
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
			{Index: 0, ID: -1, Surface: "BOS"},
			{Index: 1, ID: 2, Surface: "関西国際空港", Start: 0, End: 6, Class: TokenClass(lattice.USER)},
			{Index: 2, ID: -1, Surface: "EOS", Start: 6, End: 6, Position: len("関西国際空港")},
		}
		if len(tokens) != len(expected) {
			t.Fatalf("got %v, expected %v", tokens, expected)
		}
		for i, tok := range tokens {
			if !equalTokens(tok, expected[i]) {
				t.Errorf("%dth token, expected %v, got %v", i, expected[i], tok)
			}
		}
	})
}

func Test_AnalyzeWithSearchModeWithUserDict(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	t.Run("invalid nil user dict", func(t *testing.T) {
		_, err := New(d, UserDict(nil))
		if err == nil {
			t.Fatal("expected empty user dictionary error")
		} else if err.Error() != "empty user dictionary" {
			t.Errorf("expected empty user dictionary error, got %v", err)
		}
	})
	t.Run("user dict", func(t *testing.T) {
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
			{Index: 0, ID: -1, Surface: "BOS"},
			{Index: 1, ID: 2, Surface: "関西国際空港", Start: 0, End: 6, Class: TokenClass(lattice.USER)},
			{Index: 2, ID: -1, Surface: "EOS", Start: 6, End: 6, Position: len("関西国際空港")},
		}
		if len(tokens) != len(expected) {
			t.Fatalf("expected %v, got %v", expected, tokens)
		}
		for i, tok := range tokens {
			if !equalTokens(tok, expected[i]) {
				t.Errorf("%dth token, expected %v, got %v", i, expected[i], tok)
			}
		}
	})
}

func Test_AnalyzeWithExtendedModeWithUserDict(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	t.Run("invalid nil user dict", func(t *testing.T) {
		_, err := New(d, UserDict(nil))
		if err == nil {
			t.Fatal("expected empty user dictionary error")
		} else if err.Error() != "empty user dictionary" {
			t.Errorf("expected empty user dictionary error, got %v", err)
		}
	})
	t.Run("user dict", func(t *testing.T) {
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
			{Index: 0, ID: -1, Surface: "BOS"},
			{Index: 1, ID: 2, Surface: "関西国際空港", Start: 0, End: 6, Class: TokenClass(lattice.USER)},
			{Index: 2, ID: -1, Surface: "EOS", Start: 6, End: 6, Position: len("関西国際空港")},
		}
		if len(tokens) != len(expected) {
			t.Fatalf("expected %v, got %v", expected, tokens)
		}
		for i, tok := range tokens {
			if !equalTokens(tok, expected[i]) {
				t.Errorf("%dth token, expected %v, got %v", i, expected[i], tok)
			}
		}
	})
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
		if !equalTokens(tok, expected[i]) {
			t.Errorf("%dth token, expected %v, got %v", i, expected[i], tok)
		}
	}
}
