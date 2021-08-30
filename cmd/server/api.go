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

type TokenizerRequestBody struct {
	Input string `json:"sentence"`
	Mode  string `json:"mode,omitempty"`
}

type TokenizerResponseBody struct {
	Status bool                  `json:"status"`
	Tokens []tokenizer.TokenData `json:"tokens"`
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

	var req TokenizerRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("{\"status\":false,\"error\":\"%v\"}\n", err)
		if _, err := fmt.Fprint(w, msg); err != nil {
			log.Printf("write error, %s, %v", msg, err)
		}
		return
	}
	if req.Input == "" {
		msg := "{\"status\":true,\"tokens\":[]}\n"
		if _, err := fmt.Fprint(w, msg); err != nil {
			log.Fatalf("write error, %s, %v", msg, err)
		}
		return
	}
	mode := tokenizer.Normal
	switch req.Mode {
	case "Search":
		mode = tokenizer.Search
	case "Extended":
		mode = tokenizer.Extended
	}
	tokens := h.tokenizer.Analyze(req.Input, mode)
	var tokenData []tokenizer.TokenData
	for _, v := range tokens {
		if v.ID == tokenizer.BosEosID {
			continue
		}
		tokenData = append(tokenData, tokenizer.NewTokenData(v))
	}
	resp, err := json.Marshal(TokenizerResponseBody{
		Status: true,
		Tokens: tokenData,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Sprintf("{\"status\":false,\"error\":\"%v\"}", err)
		if _, err := fmt.Fprint(w, msg); err != nil {
			log.Printf("write error, %s, %v", msg, err)
		}
		return
	}
	if _, err := w.Write(resp); err != nil {
		log.Printf("write response json error, %v, %+v", err, req.Input)
		return
	}
}
