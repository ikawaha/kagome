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

	ipa "github.com/ikawaha/kagome-dict-ipa"
	ko "github.com/ikawaha/kagome-dict-ko"
	uni "github.com/ikawaha/kagome-dict-uni"
	"github.com/ikawaha/kagome/v2/dict"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// subcommand property
var (
	CommandName  = "server"
	Description  = `run tokenize server`
	usageMessage = "%s [-http=:6060] [-userdict userdic_file] [-dict (ipa|uni|ko)]\n"
	ErrorWriter  = os.Stderr
)

// options
type option struct {
	http    string
	udict   string
	dict    string
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
	o.flagSet.StringVar(&o.udict, "userdict", "", "user dict")
	o.flagSet.StringVar(&o.dict, "dict", "ipa", "system dict type (ipa|uni|ko)")
	return
}

func (o *option) parse(args []string) error {
	if err := o.flagSet.Parse(args); err != nil {
		return err
	}
	// validations
	if nonFlag := o.flagSet.Args(); len(nonFlag) != 0 {
		return fmt.Errorf("invalid argument: %v", nonFlag)
	}
	if o.dict != "" && o.dict != "ipa" && o.dict != "uni" && o.dict != "ko" {
		return fmt.Errorf("invalid argument: -dict %v", o.dict)
	}
	return nil
}

//OptionCheck receives a slice of args and returns an error if it was not successfully parsed
func OptionCheck(args []string) error {
	opt := newOption(ioutil.Discard, flag.ContinueOnError)
	if err := opt.parse(args); err != nil {
		return fmt.Errorf("%v, %v", CommandName, err)
	}
	return nil
}

func selectDict(name string) (*dict.Dict, error) {
	switch name {
	case "ipa":
		return ipa.Dict(), nil
	case "uni":
		return uni.Dict(), nil
	case "ko":
		return ko.Dict(), nil
	}
	return nil, fmt.Errorf("unknown name type, %v", name)
}

// command main
func command(opt *option) error {
	d, err := selectDict(opt.dict)
	if err != nil {
		return err
	}
	t := tokenizer.New(d)
	if opt.udict != "" {
		udict, err := dict.NewUserDict(opt.udict)
		if err != nil {
			return err
		}
		t.SetUserDict(udict)
	}

	mux := http.NewServeMux()
	mux.Handle("/", &TokenizeDemoHandler{tokenizer: t})
	mux.Handle("/tokenize", &TokenizeHandler{tokenizer: t})
	log.Fatal(http.ListenAndServe(opt.http, mux))
	return nil
}

// TokenizeHandler represents the tokenizer API server struct
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
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"status\":false,\"error\":\"%v\"}\n", err)
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
	j, err := json.Marshal(struct {
		Status bool     `json:"status"`
		Tokens []record `json:"tokens"`
	}{
		Status: true,
		Tokens: rsp,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"status\":false,\"error\":\"%v\"}", err)
		return
	}
	if _, err := w.Write(j); err != nil {
		log.Printf("write response json error, %v, %+v", err, body.Input)
		return
	}
}

//TokenizeDemoHandler represents the tokenizer demo server struct
type TokenizeDemoHandler struct {
	tokenizer tokenizer.Tokenizer
}

func (h *TokenizeDemoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type record struct {
		Surface       string
		Pos           string
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
	case "Search":
		m = tokenizer.Search
	case "Extended":
		// use Normal mode
	}
	if _, err := exec.LookPath(graphvizCmd); err != nil {
		cmdErr = "Error: circo/graphviz is not installed in your $PATH"
		log.Print("Error: circo/graphviz is not installed in your $PATH\n")
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
		tokens = h.tokenizer.AnalyzeGraph(w0, sen, m)
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
			case 17: // unidic
				m.Pos = strings.Join(fs[0:5], ",")
				m.Baseform = fs[10]
				m.Reading = fs[6]
				m.Pronunciation = fs[9]
			case 9:
				m.Pos = strings.Join(fs[0:5], ",")
				m.Baseform = fs[6]
				m.Reading = fs[7]
				m.Pronunciation = fs[8]
			case 7:
				m.Pos = strings.Join(fs[0:5], ",")
				m.Baseform = fs[6]
				m.Reading = "*"
				m.Pronunciation = "*"
			case 6: // unidic
				m.Pos = strings.Join(fs[0:5], ",")
				m.Baseform = "*"
				m.Reading = "*"
				m.Pronunciation = "*"
			case 3:
				m.Pos = fs[0]
				m.Baseform = fs[1]
				m.Reading = fs[2]
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

// Run receives the slice of args and executes the server
func Run(args []string) error {
	opt := newOption(ErrorWriter, flag.ExitOnError)
	if err := opt.parse(args); err != nil {
		Usage()
		PrintDefaults(flag.ExitOnError)
		return fmt.Errorf("%v, %v", CommandName, err)
	}
	return command(opt)
}

// Usage provides information on the use of the server
func Usage() {
	fmt.Fprintf(os.Stderr, usageMessage, CommandName)
}

// PrintDefaults prints out the default flags
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
      <th>Pronunciation</th>
    </tr></thread>
    <tbody id="morphs">
    {{range .Tokens}}
      <tr>
      <td>{{.Surface}}</td>
      <td>{{.Pos}}</td>
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
          var pos = "*", base = "*", reading = "*", pronoun = "*";
          var len = 0;
          if (Array.isArray(val.features)) {
            len = val.features.length;
          }
          switch (len) {
          case 17: // unidic
            pos = val.features.slice(0,5).join(",")
            base = val.features[10];
            reading = val.features[6];
            pronoun = val.features[9];
            break;
          case 9: // ipa
            pos = val.features.slice(0,5).join(",")
            base = val.features[6];
            reading = val.features[7];
            pronoun = val.features[8];
            break;
          case 7: // ipa
            pos = val.features.slice(0,5).join(",")
            base = val.features[6];
            break;
          case 6: // unidic
            pos = val.features.slice(0,5).join(",")
            break;
          case 3: // ipa
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
