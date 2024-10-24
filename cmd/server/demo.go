package server

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// assets
var (
	//go:embed asset/graph.html
	graphHTML string
	graphT    = template.Must(template.New("graph").Parse(graphHTML))

	//go:embed asset/demo.html
	demoHTML string
	demoT    = template.Must(template.New("demo").Parse(demoHTML))
)

// TokenizeDemoHandler represents the tokenizer demo server struct.
type TokenizeDemoHandler struct {
	tokenizer *tokenizer.Tokenizer
}

type record struct {
	Surface       string
	POS           string
	BaseForm      string
	Reading       string
	Pronunciation string
}

func newRecord(t tokenizer.Token) record {
	ret := record{
		Surface:       t.Surface,
		POS:           "*",
		BaseForm:      "*",
		Reading:       "*",
		Pronunciation: "*",
	}
	if v := strings.Join(t.POS(), ","); v != "" {
		ret.POS = v
	}
	if v, ok := t.BaseForm(); ok {
		ret.BaseForm = v
	}
	if v, ok := t.Reading(); ok {
		ret.Reading = v
	}
	if v, ok := t.Pronunciation(); ok {
		ret.Pronunciation = v
	}
	return ret
}

func toRecords(tokens []tokenizer.Token) []record {
	ret := make([]record, 0, len(tokens))
	for _, t := range tokens {
		if t.ID == tokenizer.BosEosID {
			continue
		}
		ret = append(ret, newRecord(t))
	}
	return ret
}

func (h *TokenizeDemoHandler) analyzeGraph(ctx context.Context, sen string, mode tokenizer.TokenizeMode) (records []record, svg string, err error) {
	var b bytes.Buffer
	tokens := h.tokenizer.AnalyzeGraph(&b, sen, mode)
	graph, err := graphviz.ParseBytes(b.Bytes())
	if err != nil {
		return nil, "", err
	}
	g, err := graphviz.New(ctx)
	if err != nil {
		return nil, "", err
	}
	b.Reset()
	if err := g.Render(ctx, graph, graphviz.SVG, &b); err != nil {
		return nil, "", fmt.Errorf("render error: %w", err)
	}
	svg = b.String()
	records = toRecords(tokens)
	return records, svg, nil
}

// ServeHTTP serves a tokenize demo server.
func (h *TokenizeDemoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sen := r.FormValue("s")
	mode := r.FormValue("r")
	lattice := r.FormValue("lattice")
	if lattice == "" {
		if err := demoT.Execute(w, struct {
			Sentence string
			RadioOpt string
		}{
			Sentence: sen,
			RadioOpt: mode,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	m := tokenizer.Normal
	switch mode {
	case "Search", "Extended": // Extended uses search mode
		m = tokenizer.Search
	}
	var cmdErr string
	records, svg, err := h.analyzeGraph(r.Context(), sen, m)
	if err != nil {
		cmdErr = "Error: " + err.Error()
		if errors.Is(err, context.DeadlineExceeded) {
			cmdErr = "Error: graphviz time out"
		}
	}
	if err := graphT.Execute(w, struct {
		Sentence string
		Tokens   []record
		CmdErr   string
		GraphSVG template.HTML
		Mode     string
	}{
		Sentence: sen,
		Tokens:   records,
		CmdErr:   cmdErr,
		GraphSVG: template.HTML(svg),
		Mode:     mode,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
