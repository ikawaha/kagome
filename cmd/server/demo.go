package server

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/ikawaha/kagome/v2/tokenizer"
)

// assets
var (
	//go:embed asset/graph.html
	graphHTML string

	//go:embed asset/demo.html
	demoHTML string
)

const (
	graphvizCmd = "circo" // "dot"
	cmdTimeout  = 25 * time.Second
)

// TokenizeDemoHandler represents the tokenizer demo server struct.
type TokenizeDemoHandler struct {
	tokenizer *tokenizer.Tokenizer
}

type record struct {
	Surface       string
	POS           string
	Baseform      string
	Reading       string
	Pronunciation string
}

func (h *TokenizeDemoHandler) analyze(sen string, mode tokenizer.TokenizeMode) (rec []record, svg string, err error) {
	if _, err := exec.LookPath(graphvizCmd); err != nil {
		return nil, "", errors.New("circo/graphviz is not installed in your $PATH")
	}
	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	var b bytes.Buffer
	cmd := exec.CommandContext(ctx, "dot", "-Tsvg")
	r0, w0 := io.Pipe()
	cmd.Stdin = r0
	cmd.Stdout = &b
	cmd.Stderr = ErrorWriter
	if err := cmd.Start(); err != nil {
		return nil, "", fmt.Errorf("process done with error, %w", err)
	}
	tokens := h.tokenizer.AnalyzeGraph(w0, sen, mode)
	if err := w0.Close(); err != nil {
		return nil, "", fmt.Errorf("pipe close error, %w", err)
	}
	if err := cmd.Wait(); err != nil {
		return nil, "", fmt.Errorf("process done with error, %w", err)
	}
	svg = b.String()
	if pos := strings.Index(svg, "<svg"); pos > 0 {
		svg = svg[pos:]
	}
	records := make([]record, 0, len(tokens))
	for _, tok := range tokens {
		if tok.ID == tokenizer.BosEosID {
			continue
		}
		m := record{
			Surface: tok.Surface,
		}
		if m.POS = strings.Join(tok.POS(), ","); m.POS == "" {
			m.POS = "*"
		}
		var ok bool
		if m.Baseform, ok = tok.BaseForm(); !ok {
			m.Baseform = "*"
		}
		if m.Reading, ok = tok.Reading(); !ok {
			m.Reading = "*"
		}
		if m.Pronunciation, ok = tok.Pronunciation(); !ok {
			m.Pronunciation = "*"
		}
		records = append(records, m)
	}
	return records, svg, nil
}

// ServeHTTP serves a tokenize demo server.
func (h *TokenizeDemoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sen := r.FormValue("s")
	mode := r.FormValue("r")
	lattice := r.FormValue("lattice")
	if lattice == "" {
		d := struct {
			Sentence string
			RadioOpt string
		}{
			Sentence: sen,
			RadioOpt: mode,
		}
		t := template.Must(template.New("top").Parse(demoHTML))
		if err := t.Execute(w, d); err != nil {
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
	records, svg, err := h.analyze(sen, m)
	if err != nil {
		cmdErr = "Error: " + err.Error()
		if errors.Is(err, context.DeadlineExceeded) {
			cmdErr = "Error: graphviz time out"
		}
	}
	d := struct {
		Sentence string
		Tokens   []record
		CmdErr   string
		GraphSvg template.HTML
		Mode     string
	}{
		Sentence: sen,
		Tokens:   records,
		CmdErr:   cmdErr,
		GraphSvg: template.HTML(svg),
		Mode:     mode,
	}
	t := template.Must(template.New("top").Parse(graphHTML))
	if err := t.Execute(w, d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
