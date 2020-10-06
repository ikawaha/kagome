package server

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/ikawaha/kagome/v2/tokenizer"
)

// TokenizeDemoHandler represents the tokenizer demo server struct
type TokenizeDemoHandler struct {
	tokenizer *tokenizer.Tokenizer
}

func (h *TokenizeDemoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type record struct {
		Surface       string
		POS           string
		Baseform      string
		Reading       string
		Pronunciation string
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
		if err := t.Execute(w, d); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	const (
		graphvizCmd = "circo" // "dot"
		cmdTimeout  = 25 * time.Second
	)
	var (
		records []record
		tokens  []tokenizer.Token
		svg     string
		cmdErr  string
	)

	m := tokenizer.Normal
	switch mode {
	case "Search", "Extended": // Extended uses search mode
		m = tokenizer.Search
	}
	if _, err := exec.LookPath(graphvizCmd); err != nil {
		cmdErr = "Error: circo/graphviz is not installed in your $PATH"
		log.Print("Error: circo/graphviz is not installed in your $PATH\n")
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
		defer cancel()
		var buf bytes.Buffer
		cmd := exec.CommandContext(ctx, "dot", "-Tsvg")
		r0, w0 := io.Pipe()
		cmd.Stdin = r0
		cmd.Stdout = &buf
		cmd.Stderr = ErrorWriter
		if err := cmd.Start(); err != nil {
			cmdErr = "Error"
			log.Printf("process done with error = %v", err)
		}
		tokens = h.tokenizer.AnalyzeGraph(w0, sen, m)
		if err := w0.Close(); err != nil {
			log.Printf("pipe close error, %v", err)
		}

		if err := cmd.Wait(); err != nil {
			switch err {
			case context.DeadlineExceeded:
				cmdErr = "Error: Graphviz time out"
			default:
				cmdErr = fmt.Sprintf("Error: process done with error, %v", err)
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
	}
	d := struct {
		Sentence string
		Tokens   []record
		CmdErr   string
		GraphSvg template.HTML
		Mode     string
	}{Sentence: sen, Tokens: records, CmdErr: cmdErr, GraphSvg: template.HTML(svg), Mode: mode}
	t := template.Must(template.New("top").Parse(graphHTML))
	if err := t.Execute(w, d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
      <th>Pronunciation</th>
    </tr></thread>
    <tbody id="morphs">
    {{range .Tokens}}
      <tr>
      <td>{{.Surface}}</td>
      <td>{{.POS}}</td>
      <td>{{.Baseform}}</td>
      <td>{{.Reading}}</td>
      <td>{{.Pronunciation}}</td>
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
      <th>Pronunciation</th>
    </tr></thread>
    <tbody id="morphs">
    </tbody>
  </table>
</div>

<script>
function cb(data, status) {
      //console.log(data);
      //console.log(status);
      if(status == "success" && Array.isArray(data.tokens)){
        $("#morphs").empty();
        $.each(data.tokens, function(i, val) {
          pos = (val.pos == null) ? "*" : val.pos;
          base = val.base_form != "" ? val.base_form : "*";
          reading = val.reading != "" ? val.reading : "*";
          pronoun = val.pronunciation!= "" ? val.pronunciation : "*";
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
  $.post('./tokenize', JSON.stringify(o), cb, 'json');
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
