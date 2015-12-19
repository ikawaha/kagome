// Copyright 2015 ikawaha
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// 	You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ikawaha/kagome/tokenizer"
)

// subcommand property
var (
	CommandName  = "server"
	Description  = `run tokenize server`
	usageMessage = "%s [-http=:6060] [-udic userdic_file] [-mode (normal|search|extended)]\n"
	ErrorWriter  = os.Stderr
)

// options
type option struct {
	http    string
	udic    string
	mode    string
	flagSet *flag.FlagSet
}

// ContinueOnError ErrorHandling // Return a descriptive error.
// ExitOnError                   // Call os.Exit(2).
// PanicOnError                  // Call panic with a descriptive error.flag.ContinueOnError
func newOption(w io.Writer, eh flag.ErrorHandling) (o *option) {
	o = &option{
		flagSet: flag.NewFlagSet(CommandName, eh),
	}
	// option settings
	o.flagSet.StringVar(&o.http, "http", ":6060", "HTTP service address (e.g., ':6060')")
	o.flagSet.StringVar(&o.udic, "udic", "", "user dictionary")
	o.flagSet.StringVar(&o.mode, "mode", "normal", "tokenize mode (normal|search|extended)")

	return
}

func (o *option) parse(args []string) (err error) {
	if err = o.flagSet.Parse(args); err != nil {
		return
	}
	// validations
	if nonFlag := o.flagSet.Args(); len(nonFlag) != 0 {
		return fmt.Errorf("invalid argument: %v", nonFlag)
	}
	if o.mode != "normal" && o.mode != "search" && o.mode != "extended" {
		return fmt.Errorf("unknown mode: %v", o.mode)
	}
	return
}

func OptionCheck(args []string) (err error) {
	opt := newOption(ioutil.Discard, flag.ContinueOnError)
	if e := opt.parse(args); e != nil {
		return fmt.Errorf("%v, %v", CommandName, e)
	}
	return nil
}

// command main
func command(opt *option) error {
	var udic tokenizer.UserDic
	if opt.udic != "" {
		var err error
		if udic, err = tokenizer.NewUserDic(opt.udic); err != nil {
			return err
		}
	}
	t := tokenizer.New()
	t.SetUserDic(udic)

	mux := http.NewServeMux()
	mux.Handle("/", &TokenizeDemoHandler{tokenizer: t})
	mux.Handle("/a", &TokenizeHandler{tokenizer: t})
	log.Fatal(http.ListenAndServe(opt.http, mux))

	return nil
}

type TokenizeHandler struct {
	tokenizer tokenizer.Tokenizer
}

func (h *TokenizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type record struct {
		ID       int      `json:"id"`
		Start    int      `json:"start"`
		End      int      `json:"end"`
		Surface  string   `json:"surface"`
		Class    string   `json:"class"`
		Features []string `json:"features"`
	}

	var body struct {
		Input string `json:"sentence"`
	}
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {
		fmt.Fprintf(w, "{\"status\":false,\"error\":\"%v\"}", e)
		return
	}
	if body.Input == "" {
		fmt.Fprint(w, "{\"status\":true,\"tokens\":[]}")
		return
	}
	tokens := h.tokenizer.Analyze(body.Input, tokenizer.Normal)
	var rsp []record
	for _, tok := range tokens {
		if tok.ID == tokenizer.BosEosID {
			continue
		}
		fs := tok.Features()
		m := record{
			ID:       tok.ID,
			Class:    fmt.Sprintf("%v", tok.Class),
			Start:    tok.Start,
			End:      tok.End,
			Surface:  tok.Surface,
			Features: fs,
		}
		rsp = append(rsp, m)
	}
	j, e := json.Marshal(struct {
		Status bool     `json:"status"`
		Tokens []record `json:"tokens"`
	}{Status: true, Tokens: rsp})
	if e != nil {
		fmt.Fprintf(w, "{\"status\":false,\"error\":\"%v\"}", e)
		return
	}
	if _, e := w.Write(j); e != nil {
		fmt.Fprintf(w, "{\"status\":false,\"error\":\"%v\"}", e)
		return
	}
}

type TokenizeDemoHandler struct {
	tokenizer tokenizer.Tokenizer
}

func (h *TokenizeDemoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type record struct {
		Surface        string
		Pos            string
		Baseform       string
		Reading        string
		Pronounciation string
	}
	sen := r.FormValue("s")
	opt := r.FormValue("r")

	const cmdTimeout = 15 * time.Second
	var (
		records []record
		tokens  []tokenizer.Token
		svg     string
		cmdErr  string
	)

	switch opt {
	case "normal":
		tokens = h.tokenizer.Analyze(sen, tokenizer.Normal)
	case "search":
		tokens = h.tokenizer.Analyze(sen, tokenizer.Search)
	case "extended":
		tokens = h.tokenizer.Analyze(sen, tokenizer.Extended)
	case "lattice":
		if _, e := exec.LookPath("dot"); e != nil {
			log.Print("graphviz is not in your future\n")
			break
		}
		var buf bytes.Buffer
		cmd := exec.Command("dot", "-Tsvg")
		r, w := io.Pipe()
		cmd.Stdin = r
		cmd.Stdout = &buf
		cmd.Stderr = ErrorWriter
		if err := cmd.Start(); err != nil {
			cmdErr = "Error"
			log.Printf("process done with error = %v", err)
		}
		h.tokenizer.Dot(sen, w)
		w.Close()

		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()
		select {
		case <-time.After(cmdTimeout):
			if err := cmd.Process.Kill(); err != nil {
				log.Fatal("failed to kill: ", err)
			}
			cmdErr = "Time out"
			<-done
		case err := <-done:
			if err != nil {
				cmdErr = "Error"
				log.Printf("process done with error = %v", err)
			}
		}
		svg = buf.String()
		if pos := strings.Index(svg, "<svg"); pos > 0 {
			svg = svg[pos:]
		}
	}
	for _, tok := range tokens {
		if tok.ID == tokenizer.BosEosID {
			continue
		}
		m := record{Surface: tok.Surface}
		fs := tok.Features()
		switch len(fs) {
		case 9:
			m.Pos = strings.Join(fs[0:5], ",")
			m.Baseform = fs[6]
			m.Reading = fs[7]
			m.Pronounciation = fs[8]
		case 7:
			m.Pos = strings.Join(fs[0:5], ",")
			m.Baseform = fs[6]
			m.Reading = "*"
			m.Pronounciation = "*"
		case 3:
			m.Pos = fs[0]
			m.Baseform = fs[1]
			m.Reading = fs[2]
			m.Pronounciation = "*"
		}
		records = append(records, m)
	}
	d := struct {
		Sentence string
		Tokens   []record
		CmdErr   string
		GraphSvg template.HTML
		RadioOpt string
	}{Sentence: sen, Tokens: records, CmdErr: cmdErr, GraphSvg: template.HTML(svg), RadioOpt: opt}
	t := template.Must(template.New("top").Parse(demoHTML))
	if e := t.Execute(w, d); e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

func Run(args []string) error {
	opt := newOption(ErrorWriter, flag.ExitOnError)
	if e := opt.parse(args); e != nil {
		Usage()
		PrintDefaults(flag.ExitOnError)
		return fmt.Errorf("%v, %v", CommandName, e)
	}
	return command(opt)
}

func Usage() {
	fmt.Fprintf(os.Stderr, usageMessage, CommandName)
}

func PrintDefaults(eh flag.ErrorHandling) {
	o := newOption(ErrorWriter, eh)
	o.flagSet.PrintDefaults()
}

var demoHTML = `
<!DOCTYPE html>
<html lang="ja">
<head>
    <style type="text/css">
      body {
        text-align: center;
      }
      div#center{
        width: 800px;
        margin: 0 auto;
        text-align: left;
      }
      .tbl{
        width: 100%;
        border-collapse: separate;
      }
      .tbl th{
        width: 20%;
        padding: 6px;
        text-align: left;
        vertical-align: top;
        color: #333;
        background-color: #eee;
        border: 1px solid #b9b9b9;
      }
      .tbl td{
        padding: 6px;
        background-color: #fff;
        border: 1px solid #b9b9b9;
      }
      .frm {
        min-height: 10px;
        padding: 0 10px 0;
        margin-bottom: 20px;
        background-color: #f5f5f5;
        border: 1px solid #e3e3e3;
        -webkit-border-radius: 4px;
        -moz-border-radius: 4px;
        border-radius: 4px;
        -webkit-box-shadow: inset 0 1px 1px rgba(0,0,0,0.05);
        -moz-box-shadow: inset 0 1px 1px rgba(0,0,0,0.05);
        box-shadow: inset 0 1px 1px rgba(0,0,0,0.05);
      }
      .txar {
         border:10px;
         padding:10px;
         font-size:1.1em;
         font-family:Arial, sans-serif;
         border:solid 1px #ccc;
         margin:0;
         width:80%;
         -webkit-border-radius: 3px;
         -moz-border-radius: 3px;
         border-radius: 3px;
         -moz-box-shadow: inset 0 0 4px rgba(0,0,0,0.2);
         -webkit-box-shadow: inset 0 0 4px rgba(0, 0, 0, 0.2);
         box-shadow: inner 0 0 4px rgba(0, 0, 0, 0.2);
      }
      .btn {
        background: -moz-linear-gradient(top,#FFF 0%,#EEE);
        background: -webkit-gradient(linear, left top, left bottom, from(#FFF), to(#EEE));
        border: 1px solid #DDD;
        border-radius: 3px;
        color:#111;
        width: 100px;
        padding: 5px 0;
        margin: 0;
      }
      #box {
        width:100%;
        margin:10px;
        auto;
      }
      #rbox {
        width:15%;
        float:right;
      }
  </style>
  <meta charset="UTF-8">
  <title>Kagome demo - japanese morphological analyzer</title>
  <!-- for IE6-8 support of HTML elements -->
  <!--[if lt IE 9]>
  <script src="http://html5shim.googlecode.com/svn/trunk/html5.js"></script>
  <![endif]-->
  <body>
  <div id="center">
  <h1>Kagome</h1>
    Kagome is an open source Japanese morphological analyzer written in Golang
    <h2>Feature summary</h2>
    <ul>
      <li><strong>Word segmentation.</strong> Segmenting text into words (or morphemes)</li>
      <li><strong>Part-of-speech tagging.</strong> Assign word-categories (nouns, verbs, particles, adjectives, etc.)</li>
      <li><strong>Lemmatization.</strong> Get dictionary forms for inflected verbs and adjectives</li>
      <li><strong>Readings.</strong> Extract readings for kanji.</li>
    </ul>
  <form class="frm" action="/_demo" method="POST">
    <div id="box">
    <textarea class="txar" rows="3" name="s" placeholder="Enter Japanese text below in UTF-8 and click tokenize.">{{.Sentence}}</textarea>
    <div id="rbox">
      <div><input type="radio" name="r" value="normal" checked>Normal</div>
      <div><input type="radio" name="r" value="search" {{if eq .RadioOpt "search"}}checked{{end}}>Search</div>
      <div><input type="radio" name="r" value="extended" {{if eq .RadioOpt "extended"}}checked{{end}}>Extended</div>
      <div><input type="radio" name="r" value="lattice" {{if eq .RadioOpt "lattice"}}checked{{end}}>Lattice</div>
    </div>
     <p><input class="btn" type="submit" value="Tokenize"/></p>
    </div>
  </form>
  {{if .CmdErr}}
    <strong>{{.CmdErr}}</strong>
  {{end}}
  {{if .GraphSvg}}
    {{.GraphSvg}}
  {{end}}
  {{if .Tokens}}
  <table class="tbl">
    <thread><tr>
      <th>Surface</th>
      <th>Part-of-Speech</th>
      <th>Base Form</th>
      <th>Reading</th>
      <th>Pronounciation</th>
    </tr></thread>
    <tbody>
    {{range .Tokens}}
      <tr>
      <td>{{.Surface}}</td>
      <td>{{.Pos}}</td>
      <td>{{.Baseform}}</td>
      <td>{{.Reading}}</td>
      <td>{{.Pronounciation}}</td>
      </tr>
    {{end}}
    </tbody>
  </table>
  {{end}}
  </div>
  </body>
</html>
`
