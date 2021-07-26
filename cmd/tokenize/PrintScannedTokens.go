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

// Variables for dependency injection or mocking for testing
var (
	JsonMarshal = json.Marshal
	FmtPrintF   = fmt.Printf
)

// parseTokenToJSON parses the token to JSON in the same format as the server mode response does.
func parseTokenToJSON(tok tokenizer.Token) ([]byte, error) {
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

	return JsonMarshal(j)
}

// printTokensAsDefault prints the tokenized text in the default format.
// The default format is: [Surface]\t[Features in CSV]\n
func printTokensAsDefault(s *bufio.Scanner, t *tokenizer.Tokenizer, mode tokenizer.TokenizeMode) error {
	for s.Scan() {
		sen := s.Text()
		tokens := t.Analyze(sen, mode)

		for i, size := 1, len(tokens); i < size; i++ {
			tok := tokens[i]
			c := tok.Features()
			if tok.Class == tokenizer.DUMMY {
				FmtPrintF("%s\n", tok.Surface)
			} else {
				FmtPrintF("%s\t%v\n", tok.Surface, strings.Join(c, ","))
			}
		}
	}

	return s.Err()
}

// printTokensInJSON prints the tokenized text in JSON format.
func printTokensInJSON(s *bufio.Scanner, t *tokenizer.Tokenizer, mode tokenizer.TokenizeMode) (err error) {
	var buff []byte

	FmtPrintF("[\n") // Begin array bracket

	for s.Scan() {
		sen := s.Text()
		tokens := t.Analyze(sen, mode)

		for _, tok := range tokens {
			if tok.ID == tokenizer.BosEosID {
				continue
			}

			if len(buff) > 0 {
				FmtPrintF("%s,\n", buff) // Print array element (JSON with comma)
			}

			if buff, err = parseTokenToJSON(tok); err != nil {
				return err
			}
		}
	}

	if s.Err() == nil {
		FmtPrintF("%s\n", buff) // Spit out the last buffer without comma to close the array
		FmtPrintF("]\n")        // End array bracket
	}

	return s.Err()
}

func PrintScannedTokens(s *bufio.Scanner, t *tokenizer.Tokenizer, mode tokenizer.TokenizeMode, opt *option) error {
	if opt.json {
		return printTokensInJSON(s, t, mode)
	}

	return printTokensAsDefault(s, t, mode)
}
