package tokenizer

import (
	"errors"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome/v2/tokenizer/lattice"
)

// TokenizeMode represents a mode of tokenize.
//
// Kagome has segmentation mode for search such as Kuromoji.
//    Normal: Regular segmentation
//    Search: Use a heuristic to do additional segmentation useful for search
//    Extended: Similar to search mode, but also unigram unknown words
type TokenizeMode int

func (m TokenizeMode) String() string {
	switch m {
	case Normal:
		return "normal"
	case Search:
		return "search"
	case Extended:
		return "extend"
	}
	return fmt.Sprintf("unknown tokenize mode (%d)", m)
}

const (
	// Normal is the normal tokenize mode.
	Normal TokenizeMode = iota + 1
	// Search is the tokenize mode for search.
	Search
	// Extended is the experimental tokenize mode.
	Extended
	// BosEosID means the beginning a sentence (BOS) or the end of a sentence (EOS).
	BosEosID = lattice.BosEosID
)

// Tokenizer represents morphological analyzer.
type Tokenizer struct {
	dict       *dict.Dict     // system dictionary
	userDict   *dict.UserDict // user dictionary
	omitBosEos bool           // omit BOS/EOS
}

// New creates a tokenizer.
func New(d *dict.Dict, opts ...Option) (*Tokenizer, error) {
	if d == nil {
		return nil, errors.New("empty dictionary")
	}
	t := &Tokenizer{dict: d}
	for _, opt := range opts {
		if err := opt(t); err != nil {
			return nil, err
		}
	}
	return t, nil
}

// Tokenize analyzes a sentence in standard tokenize mode.
func (t Tokenizer) Tokenize(input string) []Token {
	return t.Analyze(input, Normal)
}

// Wakati tokenizes a sentence and returns its divided surface strings.
func (t Tokenizer) Wakati(input string) []string {
	ts := t.Analyze(input, Normal)
	ret := make([]string, 0, len(ts))
	for _, v := range ts {
		if v.Class != DUMMY && v.Surface != "" {
			ret = append(ret, v.Surface)
		}
	}
	return ret
}

// Analyze tokenizes a sentence in the specified mode.
func (t Tokenizer) Analyze(input string, mode TokenizeMode) []Token {
	la := lattice.New(t.dict, t.userDict)
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
	tokens := make([]Token, 0, size)
	for i := range la.Output {
		n := la.Output[size-1-i]
		if t.omitBosEos && n.ID == BosEosID {
			continue
		}
		tok := Token{
			Index:    len(tokens),
			ID:       n.ID,
			Class:    TokenClass(n.Class),
			Position: n.Position,
			Start:    n.Start,
			End:      n.Start + utf8.RuneCountInString(n.Surface),
			Surface:  n.Surface,
			dict:     t.dict,
			udict:    t.userDict,
		}
		if tok.ID == BosEosID {
			if i == 0 {
				tok.Surface = "BOS"
			} else {
				tok.Surface = "EOS"
			}
		}
		tokens = append(tokens, tok)
	}
	return tokens
}

// Dot returns morphs of a sentence and exports a lattice graph to dot format in standard tokenize mode.
func (t Tokenizer) Dot(w io.Writer, input string) (tokens []Token) {
	return t.AnalyzeGraph(w, input, Normal)
}

// AnalyzeGraph returns morphs of a sentence and exports a lattice graph to dot format.
func (t Tokenizer) AnalyzeGraph(w io.Writer, input string, mode TokenizeMode) []Token {
	la := lattice.New(t.dict, t.userDict)
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
	tokens := make([]Token, 0, size)
	for i := range la.Output {
		n := la.Output[size-1-i]
		tok := Token{
			ID:      n.ID,
			Class:   TokenClass(n.Class),
			Start:   n.Start,
			End:     n.Start + utf8.RuneCountInString(n.Surface),
			Surface: n.Surface,
			dict:    t.dict,
			udict:   t.userDict,
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
	return tokens
}
