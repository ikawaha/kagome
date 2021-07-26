[![GoDev](https://pkg.go.dev/badge/github.com/ikawaha/kagome/v2)](https://pkg.go.dev/github.com/ikawaha/kagome/v2)
[![Go](https://github.com/ikawaha/kagome/workflows/Go/badge.svg)](https://github.com/ikawaha/kagome/actions?query=workflow%3AGo)
[![Coverage Status](https://coveralls.io/repos/github/ikawaha/kagome/badge.svg?branch=v2)](https://coveralls.io/github/ikawaha/kagome?branch=v2)
[![Docker Images](https://github.com/ikawaha/kagome/workflows/Docker%20Images/badge.svg)](https://github.com/ikawaha/kagome/actions?query=workflow%3A%22Docker+Images%22)
[![Docker Pulls](https://img.shields.io/docker/pulls/ikawaha/kagome.svg?style)](https://hub.docker.com/r/ikawaha/kagome/)
[![deploy](https://img.shields.io/badge/heroku-deploy_to_heroku-blue.svg)](https://heroku.com/deploy?template=https://github.com/ikawaha/kagome/tree/v2)
[![demo](https://img.shields.io/badge/demo-heroku_deployed-blue.svg)](https://kagome.herokuapp.com/)

[Kagome](https://kagome.herokuapp.com/) v2
===

Kagome is an open source Japanese morphological analyzer written in pure golang.
The dictionary/statistical models such as MeCab-IPADIC, UniDic (unidic-mecab) and so on, are able to be embedded in binaries.

### Improvements from [v1](https://github.com/ikawaha/kagome/tree/master).

* Dictionaries are maintained in a separate repository, and only the dictionaries you need are embedded in the binary.
* Brushed up and added several APIs.

# Dictionaries

|dict| source | package |
|:---|:---|:---|
|MeCab IPADIC| mecab-ipadic-2.7.0-20070801 | [github.com/ikawaha/kagome-dict/ipa](https://github.com/ikawaha/kagome-dict/tree/master/ipa)|
|UniDIC| unidic-mecab-2.1.2_src | [github.com/ikawaha/kagome-dict/uni](https://github.com/ikawaha/kagome-dict/tree/master/uni) |

**Experimental Features**

|dict|source|package|
|:---|:---|:---|
|mecab-ipadic-NEologd|mecab-ipadic-neologd| [github.com/ikawaha/kagome-ipa-neologd](https://github.com/ikawaha/kagome-dict-ipa-neologd)|
|Korean MeCab|mecab-ko-dic-2.1.1-20180720 | [github.com/ikawaha/kagome-dict-ko](https://github.com/ikawaha/kagome-dict-ko)|

## Segmentation mode for search

Kagome has segmentation mode for search such as [Kuromoji](http://www.atilika.com/en/products/kuromoji.html).

* Normal: Regular segmentation
* Search: Use a heuristic to do additional segmentation useful for search
* Extended: Similar to search mode, but also uni-gram unknown words

|Untokenized|Normal|Search|Extended|
|:-------|:---------|:---------|:---------|
|関西国際空港|関西国際空港|関西　国際　空港|関西　国際　空港|
|日本経済新聞|日本経済新聞|日本　経済　新聞|日本　経済　新聞|
|シニアソフトウェアエンジニア|シニアソフトウェアエンジニア|シニア　ソフトウェア　エンジニア|シニア　ソフトウェア　エンジニア|
|デジカメを買った|デジカメ　を　買っ　た|デジカメ　を　買っ　た|デ　ジ　カ　メ　を　買っ　た|


# Programming example

```Go
package main

import (
	"fmt"
	"strings"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

func main() {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		panic(err)
	}
	// wakati
	fmt.Println("---wakati---")
	seg := t.Wakati("すもももももももものうち")
	fmt.Println(seg)

	// tokenize
	fmt.Println("---tokenize---")
	tokens := t.Tokenize("すもももももももものうち")
	for _, token := range tokens {
		features := strings.Join(token.Features(), ",")
		fmt.Printf("%s\t%v\n", token.Surface, features)
	}
}
```

output:

```shellsession
---wakati---
[すもも も もも も もも の うち]
---tokenize---
すもも	名詞,一般,*,*,*,*,すもも,スモモ,スモモ
も	助詞,係助詞,*,*,*,*,も,モ,モ
もも	名詞,一般,*,*,*,*,もも,モモ,モモ
も	助詞,係助詞,*,*,*,*,も,モ,モ
もも	名詞,一般,*,*,*,*,もも,モモ,モモ
の	助詞,連体化,*,*,*,*,の,ノ,ノ
うち	名詞,非自立,副詞可能,*,*,*,うち,ウチ,ウチ
```

## Reference

[![実践：形態素解析 kagome v2](https://user-images.githubusercontent.com/4232165/102152682-e281e400-3eb8-11eb-91f7-13e08a8977d9.png)](https://zenn.dev/ikawaha/books/kagome-v2-japanese-tokenizer)


# Commands

## Install

**Go**

Go 1.16 or later.
```shellsession
go install github.com/ikawaha/kagome/v2@latest
```

Or use `go get`.
```shellsession
env GO111MODULE=on go get -u github.com/ikawaha/kagome/v2
```

**Homebrew tap**
```shellsession
brew install ikawaha/kagome/kagome
```

## Usage

```shellsession
$ kagome -h
Japanese Morphological Analyzer -- github.com/ikawaha/kagome/v2
usage: kagome <command>
The commands are:
   [tokenize] - command line tokenize (*default)
   server - run tokenize server
   lattice - lattice viewer
   sentence - tiny sentence splitter
   version - show version

tokenize [-file input_file] [-dict dic_file] [-userdict userdic_file] [-sysdict (ipa|uni)] [-simple false] [-mode (normal|search|extended)] [-split] [-json]
  -dict string
        dict
  -file string
        input file
  -json
        outputs in JSON format
  -mode string
        tokenize mode (normal|search|extended) (default "normal")
  -simple
        display abbreviated dictionary contents
  -split
        use tiny sentence splitter
  -sysdict string
        system dict type (ipa|uni) (default "ipa")
  -udict string
        user dict
```

### Tokenize command

```shellsession
% # interactive mode
% kagome
すもももももももものうち
すもも	名詞,一般,*,*,*,*,すもも,スモモ,スモモ
も	助詞,係助詞,*,*,*,*,も,モ,モ
もも	名詞,一般,*,*,*,*,もも,モモ,モモ
も	助詞,係助詞,*,*,*,*,も,モ,モ
もも	名詞,一般,*,*,*,*,もも,モモ,モモ
の	助詞,連体化,*,*,*,*,の,ノ,ノ
うち	名詞,非自立,副詞可能,*,*,*,うち,ウチ,ウチ
EOS
```

```shellsession
% # piped standard input
echo "すもももももももものうち" | kagome
すもも  名詞,一般,*,*,*,*,すもも,スモモ,スモモ
も      助詞,係助詞,*,*,*,*,も,モ,モ
もも    名詞,一般,*,*,*,*,もも,モモ,モモ
も      助詞,係助詞,*,*,*,*,も,モ,モ
もも    名詞,一般,*,*,*,*,もも,モモ,モモ
の      助詞,連体化,*,*,*,*,の,ノ,ノ
うち    名詞,非自立,副詞可能,*,*,*,うち,ウチ,ウチ
EOS
```

```shellsession
% # JSON output
% echo "猫" | kagome -json | jq .
[
  {
    "id": 286994,
    "start": 0,
    "end": 1,
    "surface": "猫",
    "class": "KNOWN",
    "pos": [
      "名詞",
      "一般",
      "*",
      "*"
    ],
    "base_form": "猫",
    "reading": "ネコ",
    "pronunciation": "ネコ",
    "features": [
      "名詞",
      "一般",
      "*",
      "*",
      "*",
      "*",
      "猫",
      "ネコ",
      "ネコ"
    ]
  }
]
```

```shellsession
echo "私ははにわよわわわんわん" | kagome -json | jq -r '.[].pronunciation'
ワタシ
ワ
ハニワ
ヨ
ワ
ワ
ワンワン
```

### Server command

**API**

Start a server and try to access the "/tokenize" endpoint.

```shellsession
% kagome server &
% curl -XPUT localhost:6060/tokenize -d'{"sentence":"すもももももももものうち", "mode":"normal"}' | jq .
```

**Web App**

[![demo](https://img.shields.io/badge/demo-heroku_deployed-blue.svg)](https://kagome.herokuapp.com/)


Start a server and access `http://localhost:6060`.
(To draw a lattice, demo application uses graphviz . You need graphviz installed.)

```shellsession
% kagome server &
```

### Lattice command

A debug tool of tokenize process outputs a lattice in graphviz dot format.

```shellsession
% kagome lattice 私は鰻 | dot -Tpng -o lattice.png
```

![lattice](https://user-images.githubusercontent.com/4232165/89723585-74717000-da33-11ea-886a-baab85f7a06e.png)

# Docker
[![Docker](https://dockeri.co/image/ikawaha/kagome)](https://hub.docker.com/r/ikawaha/kagome)

[![](https://images.microbadger.com/badges/image/ikawaha/kagome.svg)](https://microbadger.com/images/ikawaha/kagome "View image info on microbadger.com")

# Building to WebAssembly

You can see how kagome wasm works in [demo site.](http://ikawaha.github.io/kagome/)
The source code can be found in `./sample/wasm`.

# Licence

MIT
