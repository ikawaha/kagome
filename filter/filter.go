package filter

import (
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// Drop drops a token given the provided match function.
func Drop(tokens *[]tokenizer.Token, match func(t tokenizer.Token) bool) {
	applyFilter(match, tokens, true)
}

// Keep keeps a token given the provided match function.
func Keep(tokens *[]tokenizer.Token, match func(t tokenizer.Token) bool) {
	applyFilter(match, tokens, false)
}

func applyFilter(match func(t tokenizer.Token) bool, tokens *[]tokenizer.Token, drop bool) {
	if tokens == nil {
		return
	}
	tail := 0
	for i, v := range *tokens {
		if match(v) == drop {
			continue
		}
		if i != tail {
			(*tokens)[tail] = v
		}
		tail++
	}
	*tokens = (*tokens)[:tail]
}
