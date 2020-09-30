package filter

import (
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// Drop drops a token given the provided match function.
func Drop(tokens *[]tokenizer.Token, match func(t tokenizer.Token) bool) {
	applyFilter(match, tokens, true)
}

// PickUp picks up a token given the provided match function.
func PickUp(tokens *[]tokenizer.Token, match func(t tokenizer.Token) bool) {
	applyFilter(match, tokens, false)
}

func applyFilter(match func(t tokenizer.Token) bool, tokens *[]tokenizer.Token, drop bool) {
	if tokens == nil {
		return
	}
	tail := 0
	for i, v := range *tokens {
		if match(v) {
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
