package tokenize

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ikawaha/kagome/v2/tokenizer"
)

type TokenedJSON struct {
	ID            int      `json:"id"`
	Start         int      `json:"start"`
	End           int      `json:"end"`
	Surface       string   `json:"surface"`
	Class         string   `json:"class"`
	POS           []string `json:"pos"`
	BaseForm      string   `json:"base_form"`
	Reading       string   `json:"reading"`
	Pronunciation string   `json:"pronunciation"`
	Features      []string `json:"features"`
}

// ParseTokenToJSON parses the token to JSON in the same format as the server mode response does.
func ParseTokenToJSON(tok tokenizer.Token) ([]byte, error) {
	j := TokenedJSON{
		ID:       tok.ID,
		Start:    tok.Start,
		End:      tok.End,
		Surface:  tok.Surface,
		Class:    fmt.Sprintf("%v", tok.Class),
		POS:      tok.POS(),
		Features: tok.Features(),
	}

	j.BaseForm, _ = tok.BaseForm()
	j.Reading, _ = tok.Reading()
	j.Pronunciation, _ = tok.Pronunciation()

	return json.Marshal(j)
}

func PrintTokensAsDefault(s *bufio.Scanner, t *tokenizer.Tokenizer, mode tokenizer.TokenizeMode) error {
	for s.Scan() {
		sen := s.Text()
		tokens := t.Analyze(sen, mode)

		for i, size := 1, len(tokens); i < size; i++ {
			tok := tokens[i]
			c := tok.Features()
			if tok.Class == tokenizer.DUMMY {
				fmt.Printf("%s\n", tok.Surface)
			} else {
				fmt.Printf("%s\t%v\n", tok.Surface, strings.Join(c, ","))
			}
		}
	}

	return s.Err()
}

func PrintTokensInJSON(s *bufio.Scanner, t *tokenizer.Tokenizer, mode tokenizer.TokenizeMode) (err error) {
	var buff []byte

	fmt.Println("[") // Begin bracket

	for s.Scan() {
		sen := s.Text()
		tokens := t.Analyze(sen, mode)

		for _, tok := range tokens {
			if tok.ID == tokenizer.BosEosID {
				continue
			}

			if len(buff) > 0 {
				fmt.Printf("%s,\n", buff) // Print with comma
			}

			if buff, err = ParseTokenToJSON(tok); err != nil {
				return err
			}
		}
	}

	if s.Err() == nil {
		fmt.Printf("%s\n", buff) // Spit the last buffer w/no comma
		fmt.Println("]")         // End bracket
	}

	return s.Err()

}

func PrintScannedTokens(s *bufio.Scanner, t *tokenizer.Tokenizer, mode tokenizer.TokenizeMode, jsonOut bool) error {
	if !jsonOut {
		return PrintTokensAsDefault(s, t, mode)
	}

	return PrintTokensInJSON(s, t, mode)
}
