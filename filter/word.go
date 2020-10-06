package filter

import (
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// WordFilter represents a word filter.
type WordFilter struct {
	words map[string]struct{}
}

// NewWordFilter returns a word filter.
func NewWordFilter(words []string) *WordFilter {
	ret := &WordFilter{
		words: make(map[string]struct{}, len(words)),
	}
	for _, v := range words {
		ret.words[v] = struct{}{}
	}
	return ret
}

// Match returns true if a filter matches a given word.
func (f WordFilter) Match(w string) bool {
	_, ok := f.words[w]
	return ok
}

// Drop drops a token if a filter matches token's surface.
func (f WordFilter) Drop(tokens *[]tokenizer.Token) {
	f.apply(tokens, true)
}

// Keep keeps a token if a filter matches token's surface.
func (f WordFilter) Keep(tokens *[]tokenizer.Token) {
	f.apply(tokens, false)
}

func (f WordFilter) apply(tokens *[]tokenizer.Token, drop bool) {
	if tokens == nil {
		return
	}
	tail := 0
	for i, v := range *tokens {
		if f.Match(v.Surface) == drop {
			continue
		}
		if i != tail {
			(*tokens)[tail] = v
		}
		tail++
	}
	*tokens = (*tokens)[:tail]
}
