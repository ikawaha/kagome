package tokenize

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ikawaha/kagome/v2/tokenizer"
)

// tokenedJSON is a struct to output the tokens as JSON format.
type tokenedJSON struct {
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

// Variable for dependency injection and/or mocking for testing
var (
	JSONMarshal = json.Marshal
	FmtPrintF   = fmt.Printf
)

func fmtPrintF(format string, a ...interface{}) {
	_, _ = FmtPrintF(format, a...)
}

// parseTokenToJSON parses the token to JSON in the same format as the server mode response does.
func parseTokenToJSON(tok tokenizer.Token) ([]byte, error) {
	j := tokenedJSON{
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

	return JSONMarshal(j)
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
				fmtPrintF("%s\n", tok.Surface)
			} else {
				fmtPrintF("%s\t%v\n", tok.Surface, strings.Join(c, ","))
			}
		}
	}

	return s.Err()
}

// printTokensInJSON prints the tokenized text in JSON array format.
func printTokensInJSON(s *bufio.Scanner, t *tokenizer.Tokenizer, mode tokenizer.TokenizeMode) (err error) {
	var buff []byte

	for s.Scan() {
		fmtPrintF("[\n") // Begin array bracket

		sen := s.Text()
		tokens := t.Analyze(sen, mode)

		for _, tok := range tokens {
			if tok.ID == tokenizer.BosEosID {
				continue
			}

			if len(buff) > 0 {
				fmtPrintF("%s,\n", buff) // Print array element (JSON with comma)
			}

			if buff, err = parseTokenToJSON(tok); err != nil {
				return err
			}
		}

		fmtPrintF("%s\n", buff) // Spit out the last buffer without comma to close the array
		fmtPrintF("]\n")        // End array bracket
	}

	return s.Err()
}

// PrintScannedTokens scans and analyzes to tokenize the input and print out.
func PrintScannedTokens(s *bufio.Scanner, t *tokenizer.Tokenizer, mode tokenizer.TokenizeMode, opt *option) error {
	if opt.json {
		return printTokensInJSON(s, t, mode)
	}

	return printTokensAsDefault(s, t, mode)
}
