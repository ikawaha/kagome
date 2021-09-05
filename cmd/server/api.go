package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ikawaha/kagome/v2/tokenizer"
)

// TokenizeHandler represents the tokenizer API server struct
type TokenizeHandler struct {
	tokenizer *tokenizer.Tokenizer
}

// TokenizerRequestBody is the type of the "tokenize" endpoint HTTP request body.
type TokenizerRequestBody struct {
	Input string `json:"sentence"`
	Mode  string `json:"mode,omitempty"`
}

// TokenizerResponseBody is the response type of the "tokenize" endpoint.
type TokenizerResponseBody struct {
	Status bool                  `json:"status"`
	Tokens []tokenizer.TokenData `json:"tokens"`
}

func (h *TokenizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req TokenizerRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("{\"status\":false,\"error\":\"%v\"}", err), http.StatusBadRequest)
		return
	}
	if req.Input == "" {
		w.Write([]byte(`{"status":true,"tokens":[]}`))
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
		http.Error(w, fmt.Sprintf("{\"status\":false,\"error\":\"%v\"}", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
