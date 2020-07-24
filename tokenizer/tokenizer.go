package tokenizer

import (
	"io"
	"unicode/utf8"

	"github.com/ikawaha/kagome/v2/dict"
	"github.com/ikawaha/kagome/v2/tokenizer/lattice"
)

// TokenizeMode represents a mode of tokenize.
type TokenizeMode int

const (
	// Normal is the normal tokenize mode.
	Normal TokenizeMode = iota + 1
	// Search is the tokenize mode for search.
	Search
	// Extended is the experimental tokenize mode.
	Extended
	// BosEosID means the beginning a sentence or the end of a sentence.
	BosEosID = lattice.BosEosID
)

// Tokenizer represents morphological analyzer.
type Tokenizer struct {
	dict  *dict.Dict     // system dictionary
	udict *dict.UserDict // user dictionary
}

// New create a default tokenize.
func New(dict *dict.Dict) (t Tokenizer) {
	return Tokenizer{dict: dict}
}

// SetUserDict sets user dictionary to udict.
func (t *Tokenizer) SetUserDict(d *dict.UserDict) {
	t.udict = d
}

// Tokenize analyze a sentence in standard tokenize mode.
func (t Tokenizer) Tokenize(input string) []Token {
	return t.Analyze(input, Normal)
}

// Analyze tokenizes a sentence in the specified mode.
func (t Tokenizer) Analyze(input string, mode TokenizeMode) (tokens []Token) {
	la := lattice.New(t.dict, t.udict)
	defer la.Free()
	la.Build(input)
	m := lattice.Normal
	switch mode {
	case Normal:
		m = lattice.Normal
	case Search:
		m = lattice.Search
	case Extended:
		m = lattice.Extended
	}
	la.Forward(m)
	la.Backward(m)
	size := len(la.Output)
	tokens = make([]Token, 0, size)
	for i := range la.Output {
		n := la.Output[size-1-i]
		tok := Token{
			ID:      n.ID,
			Class:   TokenClass(n.Class),
			Start:   n.Start,
			End:     n.Start + utf8.RuneCountInString(n.Surface),
			Surface: n.Surface,
			dict:    t.dict,
			udict:   t.udict,
		}
		if tok.ID == lattice.BosEosID {
			if i == 0 {
				tok.Surface = "BOS"
			} else {
				tok.Surface = "EOS"
			}
		}
		tokens = append(tokens, tok)
	}
	return
}

// Dot returns morphs of a sentence and exports a lattice graph to dot format in standard tokenize mode.
func (t Tokenizer) Dot(w io.Writer, input string) (tokens []Token) {
	return t.AnalyzeGraph(w, input, Normal)
}

// AnalyzeGraph returns morphs of a sentence and exports a lattice graph to dot format.
func (t Tokenizer) AnalyzeGraph(w io.Writer, input string, mode TokenizeMode) (tokens []Token) {
	la := lattice.New(t.dict, t.udict)
	defer la.Free()
	la.Build(input)
	m := lattice.Normal
	switch mode {
	case Normal:
		m = lattice.Normal
	case Search:
		m = lattice.Search
	case Extended:
		m = lattice.Extended
	}
	la.Forward(m)
	la.Backward(m)
	size := len(la.Output)
	tokens = make([]Token, 0, size)
	for i := range la.Output {
		n := la.Output[size-1-i]
		tok := Token{
			ID:      n.ID,
			Class:   TokenClass(n.Class),
			Start:   n.Start,
			End:     n.Start + utf8.RuneCountInString(n.Surface),
			Surface: n.Surface,
			dict:    t.dict,
			udict:   t.udict,
		}
		if tok.ID == lattice.BosEosID {
			if i == 0 {
				tok.Surface = "BOS"
			} else {
				tok.Surface = "EOS"
			}
		}
		tokens = append(tokens, tok)
	}
	la.Dot(w)
	return
}
