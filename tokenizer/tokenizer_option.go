package tokenizer

import (
	"errors"

	"github.com/ikawaha/kagome-dict/dict"
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

// OmitBosEos is a tokenizer option to omit BOS/EOS from output tokens.
func OmitBosEos() Option {
	return func(t *Tokenizer) error {
		t.omitBosEos = true
		return nil
	}
}
