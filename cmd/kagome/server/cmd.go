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
	usageMessage = "%s [-http=:6060] [-udic userdic_file]\n"
	ErrorWriter  = os.Stderr
)

// options
type option struct {
	http    string
	udic    string
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
	o.flagSet.StringVar(&o.http, "http", ":6060", "HTTP service address")
	o.flagSet.StringVar(&o.udic, "udic", "", "user dictionary")

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
	w.Header().Set("Content-Type", "application/json")
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
		Mode  string `json:"mode,omitempty"`
	}
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"status\":false,\"error\":\"%v\"}\n", e)
		return
	}
	if body.Input == "" {
		fmt.Fprint(w, "{\"status\":true,\"tokens\":[]}\n")
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
		fs := tok.Features()
		m := record{
			ID:       tok.ID,
			Start:    tok.Start,
			End:      tok.End,
			Surface:  tok.Surface,
			Class:    fmt.Sprintf("%v", tok.Class),
			Features: fs,
		}
		rsp = append(rsp, m)
	}
	j, e := json.Marshal(struct {
		Status bool     `json:"status"`
		Tokens []record `json:"tokens"`
	}{
		Status: true,
		Tokens: rsp,
	})
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"status\":false,\"error\":\"%v\"}", e)
		return
	}
	if _, e := w.Write(j); e != nil {
		log.Printf("write response json error, %v, %+v", e, body.Input)
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
	mode := r.FormValue("r")
	lattice := r.FormValue("lattice")

	if lattice == "" {
		d := struct {
			Sentence string
			RadioOpt string
		}{Sentence: sen, RadioOpt: mode}
		t := template.Must(template.New("top").Parse(demoHTML))
		if e := t.Execute(w, d); e != nil {
			http.Error(w, e.Error(), http.StatusInternalServerError)
		}
		return
	}

	const (
		graphvizCmd = "circo" // "dot"
		cmdTimeout  = 15 * time.Second
	)
	var (
		records []record
		tokens  []tokenizer.Token
		svg     string
		cmdErr  string
	)

	m := tokenizer.Normal
	switch mode {
	case "search":
		m = tokenizer.Search
	case "extended":
		m = tokenizer.Extended
	}
	if _, e := exec.LookPath(graphvizCmd); e != nil {
		cmdErr = "Error: graphviz is not in your furure"
		log.Print("graphviz is not in your future\n")
	} else {
		var buf bytes.Buffer
		cmd := exec.Command("dot", "-Tsvg")
		r0, w0 := io.Pipe()
		cmd.Stdin = r0
		cmd.Stdout = &buf
		cmd.Stderr = ErrorWriter
		if err := cmd.Start(); err != nil {
			cmdErr = "Error"
			log.Printf("process done with error = %v", err)
		}
		tokens = h.tokenizer.AnalyzeGraph(sen, m, w0)
		w0.Close()

		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()
		select {
		case <-time.After(cmdTimeout):
			if err := cmd.Process.Kill(); err != nil {
				log.Fatal("failed to kill: ", err)
			}
			cmdErr = "Error: Graphviz time out"
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
	}
	d := struct {
		Sentence string
		Tokens   []record
		CmdErr   string
		GraphSvg template.HTML
		Mode     string
	}{Sentence: sen, Tokens: records, CmdErr: cmdErr, GraphSvg: template.HTML(svg), Mode: mode}
	t := template.Must(template.New("top").Parse(graphHTML))
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

var graphHTML = `
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
  </style>
  <meta charset="UTF-8">
  <title>Kagome demo - Japanese morphological analyzer</title>
  <!-- for IE6-8 support of HTML elements -->
  <!--[if lt IE 9]>
  <script src="http://html5shim.googlecode.com/svn/trunk/html5.js"></script>
  <![endif]-->
</head>
<body>
<div id="center">
  <table class="tbl">
    <tr><th>Input</th><td>{{.Sentence}}</td></tr>
    <tr><th>Mode</th><td>{{.Mode}}</td></tr>
  </table>

  <table class="tbl">
    <thread><tr>
      <th>Surface</th>
      <th>Part-of-Speech</th>
      <th>Base Form</th>
      <th>Reading</th>
      <th>Pronounciation</th>
    </tr></thread>
    <tbody id="morphs">
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
  <div id="graph">
  {{if .CmdErr}}
    <strong>{{.CmdErr}}</strong>
  {{end}}
  {{if .GraphSvg}}
    {{.GraphSvg}}
  {{end}}
  </div>
</div>
</body>
</html>
`
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
  <title>Kagome demo - Japanese morphological analyzer</title>
  <!-- for IE6-8 support of HTML elements -->
  <!--[if lt IE 9]>
  <script src="http://html5shim.googlecode.com/svn/trunk/html5.js"></script>
  <![endif]-->
  <script type="text/javascript" src="https://ajax.googleapis.com/ajax/libs/jquery/1.6.0/jquery.min.js"></script>
</head>
<body>
<div id="center">
  <h1>Kagome demo</h1>
  <form class="frm" action="/_demo" method="POST" oninput="tokenize()" target="_blank">
    <div id="box">
    <textarea id="inp" class="txar" rows="3" name="s"
       placeholder="Enter Japanese text below.">{{.Sentence}}</textarea>
    <div id="rbox">
      <div><label><input type="radio" name="r" value="Normal" checked>Normal</label></div>
      <div><label><input type="radio" name="r" value="Search" {{if eq .RadioOpt "search"}}checked{{end}}>Search</label></div>
      <div><label><input type="radio" name="r" value="Extended" {{if eq .RadioOpt "extended"}}checked{{end}}>Extended</label></div>
    </div>
    <p><input class="btn" type="submit" name="lattice" value="Lattice"/></p>
    </div>
  </form>

  <table class="tbl">
    <thread><tr>
      <th>Surface</th>
      <th>Part-of-Speech</th>
      <th>Base Form</th>
      <th>Reading</th>
      <th>Pronounciation</th>
    </tr></thread>
    <tbody id="morphs">
    </tbody>
  </table>
</div>

<script>
function cb(data, status) {
      console.log(data);
      console.log(status);
      if(status == "success" && Array.isArray(data.tokens)){
        $("#morphs").empty();
        $.each(data.tokens, function(i, val) {
          var pos = "*", base = "*", reading = "*", pronoun = "*";
          var len = 0;
          if (Array.isArray(val.features)) {
            len = val.features.length;
          }
          switch (len) {
          case 9:
            pos = val.features.slice(0,5).join(",")
            base = val.features[6];
            reading = val.features[7];
            pronoun = val.features[8];
            break;
          case 7:
            pos = val.features.slice(0,5).join(",")
            base = val.features[6];
            break;
          case 3:
            pos = val.features[0];
            base = val.features[1];
            reading = val.features[2];
            break;
          }
          $("#morphs").append(
          "<tr>"+"<td>" + val.surface + "</td>" +
                 "<td>" + pos + "</td>"+
                 "<td>" + base + "</td>"+
                 "<td>" + reading + "</td>"+
                 "<td>" + pronoun + "</td>"+
          "</tr>"
          );
        });
      }
}

function tokenize() {
  var s = document.getElementById("inp").value;
  var m = $('input[name="r"]').filter(':checked').val();
  var o = {"sentence" : s, "mode" : m};
  $.post('./a', JSON.stringify(o), cb, 'json');
}

$('input[name="r"]:radio').change( function() {
  var s = document.getElementById("inp").value;
  var m = $('input[name="r"]').filter(':checked').val();
  var o = {"sentence" : s, "mode" : m};
  $.post('./a', JSON.stringify(o), cb, 'json');
})
</script>

</body>
</html>
`
