package filter

import (
	"github.com/ikawaha/kagome/v2/tokenizer"
)

type (
	POS = []string
)

type POSFilter struct {
	filter *FeaturesFilter
}

// NewPOSFilter returns a part of speech filter.
func NewPOSFilter(stops ...POS) *POSFilter {
	return &POSFilter{
		filter: NewFeaturesFilter(stops...),
	}
}

// Match returns true if a filter matches given POS.
func (f POSFilter) Match(p POS) bool {
	return f.filter.Match(p)
}

// Drop drops a token if a filter matches token's POS.
func (f POSFilter) Drop(tokens *[]tokenizer.Token) {
	f.apply(tokens, true)
}

// Keep keeps a token if a filter matches token's POS.
func (f POSFilter) Keep(tokens *[]tokenizer.Token) {
	f.apply(tokens, false)
}

func (f POSFilter) apply(tokens *[]tokenizer.Token, drop bool) {
	if tokens == nil {
		return
	}
	tail := 0
	for i, v := range *tokens {
		if f.Match(v.POS()) {
			if drop {
				continue
			}
		} else if !drop {
			continue
		}
		if i != tail {
			(*tokens)[tail] = (*tokens)[i]
		}
		tail++
	}
	*tokens = (*tokens)[:tail]
}
