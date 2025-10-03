[![GoDev](https://pkg.go.dev/badge/github.com/ikawaha/kagome/v2)](https://pkg.go.dev/github.com/ikawaha/kagome/v2)
[![Go](https://github.com/ikawaha/kagome/workflows/Go/badge.svg)](https://github.com/ikawaha/kagome/actions?query=workflow%3AGo)
[![Release](https://github.com/ikawaha/kagome/actions/workflows/release.yml/badge.svg?branch=)](https://github.com/ikawaha/kagome/actions/workflows/release.yml)
[![Coverage Status](https://coveralls.io/repos/github/ikawaha/kagome/badge.svg?branch=v2)](https://coveralls.io/github/ikawaha/kagome?branch=v2)
[![Docker Pulls](https://img.shields.io/docker/pulls/ikawaha/kagome.svg?style)](https://hub.docker.com/r/ikawaha/kagome/)

# Kagome v2

Kagome is an open source Japanese morphological analyzer written in pure Go. It can tokenize Japanese text into words and analyze parts of speech, with dictionaries embedded in the binary for easy deployment.

> [!NOTE]
> **Key features** (Improvements from [v1](https://github.com/ikawaha/kagome/tree/master)):
>
> * Self-contained binaries with embedded dictionaries (MeCab-IPADIC, UniDic)
> * Multiple segmentation modes for different use cases
> * RESTful API server mode for production use
> * WebAssembly support for browser environments

## Index

* [Basic Usage](#basic-usage)
  * [Command line](#command-line)
  * [As a Go library](#as-a-go-library)
* [Install](#install)
* [Commands](#commands)
  * [Tokenize command](#tokenize-command)
  * [Server command](#server-command)
    * [RESTful API](#restful-api)
    * [Web App](#web-app)
  * [Lattice command](#lattice-command)
  * [Sentence command](#sentence-command)
* [Dictionaries](#dictionaries)
* [Segmentation modes](#segmentation-modes)
* [Docker](#docker)
* [WebAssembly](#webassembly)
* [Reference](#reference)
* [License](#license)

## Basic Usage

### Command line

```shellsession
% kagome -h
Japanese Morphological Analyzer -- github.com/ikawaha/kagome/v2
usage: kagome <command>
The commands are:
   [tokenize] - command line tokenize (*default)
   server - run tokenize server
   lattice - lattice viewer
   sentence - tiny sentence splitter
   version - show version

tokenize [-file input_file] [-dict dic_file] [-userdict user_dic_file] [-sysdict (ipa|uni)] [-simple false] [-mode (normal|search|extended)] [-split] [-json]
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

```shellsession
% # piped standard input
% echo "すもももももももものうち" | kagome
すもも	名詞,一般,*,*,*,*,すもも,スモモ,スモモ
も	助詞,係助詞,*,*,*,*,も,モ,モ
もも	名詞,一般,*,*,*,*,もも,モモ,モモ
も	助詞,係助詞,*,*,*,*,も,モ,モ
もも	名詞,一般,*,*,*,*,もも,モモ,モモ
の	助詞,連体化,*,*,*,*,の,ノ,ノ
うち	名詞,非自立,副詞可能,*,*,*,うち,ウチ,ウチ
EOS
```

* For more details, see the [Commands section](#commands).

### As a Go library

```sh
# Install Kagome module
go get github.com/ikawaha/kagome/v2
```

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
  // wakati (simple word splitting/segmentation)
  fmt.Println("---wakati---")
  seg := t.Wakati("すもももももももものうち")
  fmt.Println(seg)

  // tokenize w/ morphological analysis
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

* For more examples, see:
  * [examples directory](https://github.com/ikawaha/kagome/tree/v2/_examples)
  * [GoDoc](https://pkg.go.dev/github.com/ikawaha/kagome/v2)

## Install

To **get the `kagome` command line tool**, choose your preferred installation method below:

* **Go (recommended)**

  ```shellsession
  go install github.com/ikawaha/kagome/v2@latest
  ```

* **Homebrew**

  ```shellsession
  # macOS and Linux (for both AMD64 and Arm64)
  brew install ikawaha/kagome/kagome
  ```

* **Manual Install**

  * For manual installation, download and extract the appropriate archived file for your OS and architecture from the [releases page](https://github.com/ikawaha/kagome/releases/latest).
  * Note that the extracted binary must be placed in an accessible directory with execution permission.

* **Docker/Docker Compose**

  * See the [Docker section](#docker) below

## Commands

Major sub-commands of `kagome` command line tool.

### Tokenize command

```shellsession
% # interactive/REPL mode
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
% # (For jq command see https://jqlang.org/)
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
% # word splitting/segmentation only (equivalent to "wakati" functionality)
% echo "すもももももももものうち" | kagome -json | jq -r '[.[].surface] | join(" ")'
すもも も もも も もも の うち
```

```shellsession
% # Extract only pronunciations using jq (for Text-to-Speech purposes, etc.)
% echo "私ははにわよわわわんわん" | kagome -json | jq -r '.[].pronunciation'
ワタシ
ワ
ハニワ
ヨ
ワ
ワ
ワンワン
```

### Server command

For continuous usage, `kagome` provides a server mode to decouple the startup time of the tokenizer.

#### RESTful API

Start a server and try to access the "/tokenize" endpoint.

```shellsession
% kagome server &
% curl -XPUT localhost:6060/tokenize -d'{"sentence":"すもももももももものうち", "mode":"normal"}' | jq .
```

#### Web App

Start a server and access `http://localhost:6060` in your browser.

```shellsession
% kagome server &
```

![webapp](https://raw.githubusercontent.com/wiki/ikawaha/kagome/images/demoapp.gif)

> [!IMPORTANT]
> The demo web application uses [graphviz](https://graphviz.org/) to draw a lattice. You need graphviz to be installed on your system.
>
> [!TIP]
> Kagome can be compiled to WebAssembly (wasm) and run locally in a web browser as well. For details, see the [WebAssembly section](#webassembly).
>
> * Wasm Demo: [https://ikawaha.github.io/kagome/](https://ikawaha.github.io/kagome/)

### Lattice command

A debug tool of tokenize process outputs a lattice in graphviz dot format.

```shellsession
% kagome lattice 私は鰻 | dot -Tpng -o lattice.png
```

![lattice](https://user-images.githubusercontent.com/4232165/89723585-74717000-da33-11ea-886a-baab85f7a06e.png)

### Sentence command

Split long text into sentences:

```shellsession
% echo "吾輩は猫である。名前はまだ無い。" | kagome sentence
吾輩は猫である。
名前はまだ無い。
```

This command is useful if a single line of data is too lengthy, and you want to avoid errors such as `bufio.Scanner: token too long`.

```shellsession
% echo "吾輩は猫である。名前はまだ無い。" | kagome -json | jq -r '[.[].surface] | join(" ")'
吾輩 は 猫 で ある 。 名前 は まだ 無い 。

% echo "吾輩は猫である。名前はまだ無い。" | kagome sentence | kagome -json | jq -r '[.[].surface] | join(" ")'
吾輩 は 猫 で ある 。
名前 は まだ 無い 。
```

This command is equivalent to the `-split` option of the `tokenize` command.

```shellsession
% echo "吾輩は猫である。名前はまだ無い。" | kagome -split -json | jq -r '[.[].surface] | join(" ")'
吾輩 は 猫 で ある 。
名前 は まだ 無い 。
```

## Dictionaries

* Currently supported dictionaries by default.

  |dict| source | package |
  |:---|:---|:---|
  |MeCab IPADIC| mecab-ipadic-2.7.0-20070801 | [github.com/ikawaha/kagome-dict/ipa](https://github.com/ikawaha/kagome-dict/tree/master/ipa)|
  |UniDIC| unidic-mecab-2.1.2_src | [github.com/ikawaha/kagome-dict/uni](https://github.com/ikawaha/kagome-dict/tree/master/uni) |

* Experimental Features

  |dict|source|package|
  |:---|:---|:---|
  |mecab-ipadic-NEologd|mecab-ipadic-neologd| [github.com/ikawaha/kagome-ipa-neologd](https://github.com/ikawaha/kagome-dict-ipa-neologd)|
  |Korean MeCab|mecab-ko-dic-2.1.1-20180720 | [github.com/ikawaha/kagome-dict-ko](https://github.com/ikawaha/kagome-dict-ko)|

> [!NOTE]
> For more details and differences between the dictionaries, see the [wiki](https://github.com/ikawaha/kagome/wiki/About-the-dictionary).

## Segmentation modes

Similar to [Kuromoji](https://www.atilika.org/), Kagome also supports various **segmentation modes** (splitting strategies) to tokenize the input text.

* **Normal:** Regular segmentation
* **Search:** Use a heuristic to perform additional segmentation that is **useful for search** purposes
* **Extended:** Similar to search mode, but also unknown words with [uni-grams](https://en.wikipedia.org/wiki/N-gram)

|Untokenized|Normal|Search|Extended|
|:-------|:---------|:---------|:---------|
|関西国際空港|関西国際空港|関西　国際　空港|関西　国際　空港|
|日本経済新聞|日本経済新聞|日本　経済　新聞|日本　経済　新聞|
|シニアソフトウェアエンジニア|シニアソフトウェアエンジニア|シニア　ソフトウェア　エンジニア|シニア　ソフトウェア　エンジニア|
|デジカメを買った|デジカメ　を　買っ　た|デジカメ　を　買っ　た|デ　ジ　カ　メ　を　買っ　た|

> [!NOTE]
>If your purpose is for search, try changing the mode before switching to another dictionary.

## Docker

[![Docker](https://dockerico.blankenship.io/image/ikawaha/kagome)](https://hub.docker.com/r/ikawaha/kagome)

We provide `scratch`-based Docker images that simply run the `kagome` command line tool on various architectures: AMD64, Arm64, Arm32 (Arm v5, v6 and v7)

* Pull the image

  ```sh
  docker pull ikawaha/kagome:latest
  ```

  ```sh
  # Alternatively, you can pull from GitHub Container Registry
  docker pull ghcr.io/ikawaha/kagome:latest
  ```

* Run the command via Docker

  ```sh
  # Interactive/REPL mode
  docker run --rm -it ikawaha/kagome:latest
  ```

  ```sh
  # If pulling from GitHub Container Registry
  docker run --rm -it ghcr.io/ikawaha/kagome:latest
  ```

* Run the server via Docker

  ```sh
  # Server mode (http://localhost:6060)
  docker run --rm -p 6060:6060 ikawaha/kagome:latest server
  ```

  ```sh
  # If pulling from GitHub Container Registry
  docker run --rm -p 6060:6060 ghcr.io/ikawaha/kagome:latest server
  ```

* `docker-compose.yml` example

  ```yaml
  services:
    kagome:
      image: ikawaha/kagome:latest
      ports: ["6060:6060"]
      command: server
      restart: unless-stopped
  ```

> **Note:** Base image doesn't include Graphviz. For lattice visualization, see [examples](./_examples/server_docker_graphviz/).

## WebAssembly

Kagome compiles to WebAssembly for browser use.

* **Live demo:** [https://ikawaha.github.io/kagome/](https://ikawaha.github.io/kagome/)
* **Source code:** [./_examples/wasm](./_examples/wasm)

## Reference

* Detailed Reference Manual in Japanese:

  [![実践：形態素解析 kagome v2](https://user-images.githubusercontent.com/4232165/102152682-e281e400-3eb8-11eb-91f7-13e08a8977d9.png)](https://zenn.dev/ikawaha/books/kagome-v2-japanese-tokenizer)

* Community Wiki in English:
  * [https://github.com/ikawaha/kagome/wiki](https://github.com/ikawaha/kagome/wiki)

## License

* MIT
