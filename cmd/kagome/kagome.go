package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ikawaha/kagome"
)

type KagomeHandler struct {
	tokenizer *kagome.Tokenizer
}

func (h *KagomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type morph struct {
		Id       int      `json:"id"`
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
		fmt.Fprintf(w, "{\"status\":false,\"message\":\"%v\"}", e)
		return
	}
	if body.Input == "" {
		fmt.Fprint(w, "{\"status\":true,\"tokens\":[]}")
		return
	}
	tokens := h.tokenizer.Tokenize(body.Input)
	var ans []morph
	for _, tok := range tokens {
		if tok.Id == kagome.BosEosId {
			continue
		}
		fs := tok.Features()
		m := morph{
			Id:       tok.Id,
			Class:    fmt.Sprintf("%v", tok.Class),
			Start:    tok.Start,
			End:      tok.End,
			Surface:  tok.Surface,
			Features: fs,
		}
		ans = append(ans, m)
	}
	j, e := json.Marshal(struct {
		Status bool    `json:"status"`
		Tokens []morph `json:"tokens"`
	}{Status: true, Tokens: ans})
	if e != nil {
		fmt.Fprintf(w, "{\"status\":false,\"message\":\"%v\"}", e)
		return
	}
	w.Write(j)
}

type KagomeDemoHandler struct {
	tokenizer *kagome.Tokenizer
}

func (h *KagomeDemoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type morph struct {
		Surface        string
		Pos            string
		Baseform       string
		Reading        string
		Pronounciation string
	}
	sen := r.FormValue("s")
	tokens := h.tokenizer.Tokenize(sen)
	var morphs []morph
	for _, tok := range tokens {
		if tok.Id == kagome.BosEosId {
			continue
		}
		m := morph{Surface: tok.Surface}
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
		morphs = append(morphs, m)
	}
	d := struct {
		Sentence string
		Tokens   []morph
	}{Sentence: sen, Tokens: morphs}
	t := template.Must(template.New("top").Parse(demo_html))
	t.Execute(w, d)
}

var usageMessage = "usage: kagome [-f input_file | --http addr] [-u userdic_file]"

func usage() {
	fmt.Fprintln(os.Stderr, usageMessage)
	flag.PrintDefaults()
	os.Exit(0)
}

var (
	fHttp        = flag.String("http", "", "HTTP service address (e.g., ':6060')")
	fInputFile   = flag.String("file", "", "input file")
	fUserDicFile = flag.String("udic", "", "user dic")
)

func Main() {
	if *fHttp != "" && *fInputFile != "" {
		usage()
	}

	var udic *kagome.UserDic
	if *fUserDicFile != "" {
		var err error
		udic, err = kagome.NewUserDic(*fUserDicFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if *fHttp != "" {
		t := kagome.NewThreadsafeTokenizer()
		if udic != nil {
			t.SetUserDic(udic)
		}
		hTok := &KagomeHandler{tokenizer: t}
		hDem := &KagomeDemoHandler{tokenizer: t}
		mux := http.NewServeMux()
		mux.Handle("/", hTok)
		mux.Handle("/_demo", hDem)
		log.Fatal(http.ListenAndServe(*fHttp, mux))
		os.Exit(0)
	}

	var inputFile = os.Stdin
	if *fInputFile != "" {
		var err error
		inputFile, err = os.Open(*fInputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer inputFile.Close()
	}

	t := kagome.NewTokenizer()
	if udic != nil {
		t.SetUserDic(udic)
	}
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := t.Tokenize(line)
		for i, size := 1, len(tokens); i < size; i++ {
			tok := tokens[i]
			c := tok.Features()
			if tok.Class == kagome.DUMMY {
				fmt.Printf("%s\n", tok.Surface)
			} else {
				fmt.Printf("%s\t%v\n", tok.Surface, strings.Join(c, ","))
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()
	Main()
}

var demo_html = `
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
         border:0;
         padding:10px;
         font-size:1.1em;
         font-family:Arial, sans-serif;
         border:solid 1px #ccc;
         margin:0 0 0px;
         width:70%;
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
        width: 50px;
        padding: 5px 0;
        margin: 0;
      }
  </style>
  <meta charset="UTF-8">
  <title>kagome demo - japanese morphological analyzer</title>
  <!-- for IE6-8 support of HTML elements -->
  <!--[if lt IE 9]>
  <script src="http://html5shim.googlecode.com/svn/trunk/html5.js"></script>
  <![endif]-->
  <body>
  <div id="center">
  <h1>Kagome Demo</h1>
  <form class="frm" action="/_demo" method="POST">
    <p><textarea class="txar" rows="2" name="s" placeholder="日本語の文章(UTF-8)を入力してから実行をクリックしてください">{{.Sentence}}</textarea></p>
   <p><input class="btn" type="submit" value="実行"/></p>
  </form>
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
