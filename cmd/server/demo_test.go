package server

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/ikawaha/kagome/v2/tokenizer"
)

func TestTokenizeDemoHandler_ServeHTTP(t *testing.T) {
	tnz, err := tokenizer.New(loadTestDict(t))
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	t.Run("no path params", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		w := httptest.NewRecorder()
		(&TokenizeDemoHandler{tokenizer: tnz}).ServeHTTP(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("http status: got %d, want %d", got, want)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("response body read error, %v", err)
		}
		if !bytes.Contains(body, []byte(`Kagome demo - Japanese morphological analyzer`)) {
			t.Errorf("demo title not be found")
		}
	})
	t.Run("w/o lattice", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, `/?s=ねこです&r=Extended`, nil)
		w := httptest.NewRecorder()
		(&TokenizeDemoHandler{tokenizer: tnz}).ServeHTTP(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("http status: got %d, want %d", got, want)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("response body read error, %v", err)
		}
		if !bytes.Contains(body, []byte(`Kagome demo - Japanese morphological analyzer`)) {
			t.Errorf("demo title not be found")
		}
		if bytes.Contains(body, []byte(`<svg width=`)) {
			t.Errorf("unexpected svg found")
		}
	})

	t.Run("w/ lattice", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, `/?s=ねこです&r=Search&lattice=true`, nil)
		w := httptest.NewRecorder()
		(&TokenizeDemoHandler{tokenizer: tnz}).ServeHTTP(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("http status: got %d, want %d", got, want)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("response body read error, %v", err)
		}
		if !bytes.Contains(body, []byte(`Kagome demo - Japanese morphological analyzer`)) {
			t.Errorf("demo title not be found")
		}
		if !bytes.Contains(body, []byte(`<svg width=`)) {
			t.Errorf("svg not found")
		}
	})
}

func TestTokenizeDemoHandler_analyzeGraph(t *testing.T) {
	if _, err := exec.LookPath(graphvizCmd); err != nil {
		t.Skipf("graphviz command not found, %v", err)
	}
	tnz, err := tokenizer.New(loadTestDict(t))
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	handler := TokenizeDemoHandler{tokenizer: tnz}
	records, svg, err := handler.analyzeGraph(context.TODO(), "ねこです", tokenizer.Normal)
	if err != nil {
		t.Fatalf("unexpected error, analyzeGraph() failed, %v", err)
	}
	if got, want := records, []record{
		{
			Surface:       "ねこ",
			POS:           "名詞,一般,*,*",
			BaseForm:      "ねこ",
			Reading:       "ネコ",
			Pronunciation: "ネコ",
		},
		{
			Surface:       "です",
			POS:           "助動詞,*,*,*",
			BaseForm:      "です",
			Reading:       "デス",
			Pronunciation: "デス",
		},
	}; !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %v", got, want)
	}
	if !strings.HasPrefix(svg, "<svg width=") {
		if len(svg) > 50 {
			svg = svg[:50]
		}
		t.Errorf("broken svg, %s", svg)
	}
}

func Test_newRecord(t *testing.T) {
	dict := loadTestDict(t)
	tnz, err := tokenizer.New(dict)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tokens := tnz.Tokenize("ねこ")
	if len(tokens) != 3 {
		t.Fatalf("unexpected tokenize error, got %+v", tokens)
	}
	tests := []struct {
		name  string
		token tokenizer.Token
		want  record
	}{
		{
			name: "surface only",
			token: tokenizer.Token{
				Surface: "piyo",
			},
			want: record{
				Surface:       "piyo",
				POS:           "*",
				BaseForm:      "*",
				Reading:       "*",
				Pronunciation: "*",
			},
		},
		{
			name:  "tokens[0]:BOS",
			token: tokens[0],
			want: record{
				Surface:       "BOS",
				POS:           "*",
				BaseForm:      "*",
				Reading:       "*",
				Pronunciation: "*",
			},
		},
		{
			name:  "tokens[1]:ねこ",
			token: tokens[1],
			want: record{
				Surface:       "ねこ",
				POS:           strings.Join(tokens[1].POS(), ","),
				BaseForm:      "ねこ",
				Reading:       "ネコ",
				Pronunciation: "ネコ",
			},
		},
		{
			name:  "tokens[2]:EOS",
			token: tokens[2],
			want: record{
				Surface:       "EOS",
				POS:           "*",
				BaseForm:      "*",
				Reading:       "*",
				Pronunciation: "*",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newRecord(tt.token); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toRecords(t *testing.T) {
	dict := loadTestDict(t)
	tnz, err := tokenizer.New(dict)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tokens := tnz.Tokenize("ねこ")
	if len(tokens) != 3 {
		t.Fatalf("unexpected tokenize error, got %+v", tokens)
	}
	tests := []struct {
		name   string
		tokens []tokenizer.Token
		want   []record
	}{
		{
			name:   "empty",
			tokens: []tokenizer.Token{},
			want:   []record{},
		},
		{
			name:   "sentence",
			tokens: tokens,
			want: []record{
				{
					Surface:       "ねこ",
					POS:           strings.Join(tokens[1].POS(), ","),
					BaseForm:      "ねこ",
					Reading:       "ネコ",
					Pronunciation: "ネコ",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toRecords(tt.tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toRecords() = %v, want %v", got, tt.want)
			}
		})
	}
}
