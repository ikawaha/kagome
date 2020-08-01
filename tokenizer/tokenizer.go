package tokenizer

import (
	"errors"
	"fmt"
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

// Option represents an option for the tokenizer.
type Option func(*Tokenizer) error

// Nop represents a no operation option.
func Nop() Option {
	return func(t *Tokenizer) error {
		return nil
	}
}

// UserDict is a tokenizer option to sets a user dictionary.
func UserDict(d *dict.UserDict) Option {
	return func(t *Tokenizer) error {
		if d == nil {
			return errors.New("empty user dictionary")
		}
		t.userDict = d
		return nil
	}
}

// Tokenizer represents morphological analyzer.
type Tokenizer struct {
	dict     *dict.Dict     // system dictionary
	userDict *dict.UserDict // user dictionary
}

// New creates a tokenizer.
func New(d *dict.Dict, opts ...Option) (*Tokenizer, error) {
	if d == nil {
		return nil, errors.New("empty dictionary")
	}
	t := &Tokenizer{dict: d}
	for _, opt := range opts {
		if err := opt(t); err != nil {
			return nil, fmt.Errorf("invalid option: %v", err)
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
func (t Tokenizer) Analyze(input string, mode TokenizeMode) (tokens []Token) {
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
	return
}

// Dot returns morphs of a sentence and exports a lattice graph to dot format in standard tokenize mode.
func (t Tokenizer) Dot(w io.Writer, input string) (tokens []Token) {
	return t.AnalyzeGraph(w, input, Normal)
}

// AnalyzeGraph returns morphs of a sentence and exports a lattice graph to dot format.
func (t Tokenizer) AnalyzeGraph(w io.Writer, input string, mode TokenizeMode) (tokens []Token) {
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
	return
}
