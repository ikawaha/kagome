package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ikawaha/kagome/v2/tokenizer"
)

// TokenizeHandler represents the tokenizer API server struct
type TokenizeHandler struct {
	tokenizer *tokenizer.Tokenizer
}

func (h *TokenizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type record struct {
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

	var body struct {
		Input string `json:"sentence"`
		Mode  string `json:"mode,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("{\"status\":false,\"error\":\"%v\"}\n", err)
		if _, err := fmt.Fprint(w, msg); err != nil {
			log.Printf("write error, %s, %v", msg, err)
		}
		return
	}
	if body.Input == "" {
		msg := "{\"status\":true,\"tokens\":[]}\n"
		if _, err := fmt.Fprint(w, msg); err != nil {
			log.Fatalf("write error, %s, %v", msg, err)
		}
		return
	}
	mode := tokenizer.Normal
	switch body.Mode {
	case "Search":
		mode = tokenizer.Search
	case "Extended":
		mode = tokenizer.Extended
	}
	tokens := h.tokenizer.Analyze(body.Input, mode)
	var rsp []record
	for _, tok := range tokens {
		if tok.ID == tokenizer.BosEosID {
			continue
		}
		m := record{
			ID:       tok.ID,
			Start:    tok.Start,
			End:      tok.End,
			Surface:  tok.Surface,
			Class:    fmt.Sprintf("%v", tok.Class),
			POS:      tok.POS(),
			Features: tok.Features(),
		}
		m.BaseForm, _ = tok.BaseForm()
		m.Reading, _ = tok.Reading()
		m.Pronunciation, _ = tok.Pronunciation()
		rsp = append(rsp, m)
	}
	j, err := json.Marshal(struct {
		Status bool     `json:"status"`
		Tokens []record `json:"tokens"`
	}{
		Status: true,
		Tokens: rsp,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Sprintf("{\"status\":false,\"error\":\"%v\"}", err)
		if _, err := fmt.Fprint(w, msg); err != nil {
			log.Printf("write error, %s, %v", msg, err)
		}
		return
	}
	if _, err := w.Write(j); err != nil {
		log.Printf("write response json error, %v, %+v", err, body.Input)
		return
	}
}
