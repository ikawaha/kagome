package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

const testDictPath = "../../testdata/ipa.dict"

func TestTokenizerAPI(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected dict loading error, %v", err)
	}
	tnz, err := tokenizer.New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	p, err := json.Marshal(TokenizerRequestBody{
		Input: "ねこですねこはいます",
	})
	if err != nil {
		t.Fatalf("unexpected json marshal error, %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/tokenizer", bytes.NewReader(p))
	w := httptest.NewRecorder()
	(&TokenizeHandler{tokenizer: tnz}).ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("http status code got %d(%s), want %d", got, resp.Status, want)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unexpected response body error, %v", err)
	}
	var body TokenizerResponseBody
	if err := json.Unmarshal(b, &body); err != nil {
		t.Fatalf("unexpected json unmarshal error, %v, %s", err, string(b))
	}
	if got, want := body.Status, true; got != want {
		t.Errorf("status: got %t, want %t", got, want)
	}
	if got, want := len(body.Tokens), 6; got != want {
		t.Fatalf("len: got %d, want %d, %+v", got, want, body.Tokens)
	}
	var buf bytes.Buffer
	for _, v := range body.Tokens {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("unexpected json marshal error, %v", err)
		}
		buf.Write(b)
		buf.WriteString("\n")
	}
	want := `{"id":54873,"start":0,"end":2,"surface":"ねこ","class":"KNOWN","pos":["名詞","一般","*","*"],"base_form":"ねこ","reading":"ネコ","pronunciation":"ネコ","features":["名詞","一般","*","*","*","*","ねこ","ネコ","ネコ"]}
{"id":47492,"start":2,"end":4,"surface":"です","class":"KNOWN","pos":["助動詞","*","*","*"],"base_form":"です","reading":"デス","pronunciation":"デス","features":["助動詞","*","*","*","特殊・デス","基本形","です","デス","デス"]}
{"id":54873,"start":4,"end":6,"surface":"ねこ","class":"KNOWN","pos":["名詞","一般","*","*"],"base_form":"ねこ","reading":"ネコ","pronunciation":"ネコ","features":["名詞","一般","*","*","*","*","ねこ","ネコ","ネコ"]}
{"id":57061,"start":6,"end":7,"surface":"は","class":"KNOWN","pos":["助詞","係助詞","*","*"],"base_form":"は","reading":"ハ","pronunciation":"ワ","features":["助詞","係助詞","*","*","*","*","は","ハ","ワ"]}
{"id":3664,"start":7,"end":8,"surface":"い","class":"KNOWN","pos":["動詞","自立","*","*"],"base_form":"いる","reading":"イ","pronunciation":"イ","features":["動詞","自立","*","*","一段","連用形","いる","イ","イ"]}
{"id":68729,"start":8,"end":10,"surface":"ます","class":"KNOWN","pos":["助動詞","*","*","*"],"base_form":"ます","reading":"マス","pronunciation":"マス","features":["助動詞","*","*","*","特殊・マス","基本形","ます","マス","マス"]}
`
	if got := buf.String(); got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
