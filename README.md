Kagome v2
===

Kagome is an open source Japanese morphological analyzer written in pure golang.
The dictionary/statiscal models such as MeCab-IPADIC, UniDic (unidic-mecab), Korean MeCab and so on, be able to embedded in binaries.

### Improvements from [v1](https://github.com/ikawaha/kagome/tree/master).

* Dictionaries are maintained in a separate repository, and only the dictionaries you need are embedded in the binay.
* Brushed up and added several APIs.

# Usage

```Go
package main

import (
	"fmt"
	"strings"

	ipa "github.com/ikawaha/kagome-dict-ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

func main() {
	t, err := tokenizer.New(ipa.Dict())
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
		if token.Class == tokenizer.DUMMY {
			// BOS: Begin Of Sentence, EOS: End Of Sentence.
			fmt.Printf("%s\n", token.Surface)
			continue
		}
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
BOS
すもも	名詞,一般,*,*,*,*,すもも,スモモ,スモモ
も	助詞,係助詞,*,*,*,*,も,モ,モ
もも	名詞,一般,*,*,*,*,もも,モモ,モモ
も	助詞,係助詞,*,*,*,*,も,モ,モ
もも	名詞,一般,*,*,*,*,もも,モモ,モモ
の	助詞,連体化,*,*,*,*,の,ノ,ノ
うち	名詞,非自立,副詞可能,*,*,*,うち,ウチ,ウチ
EOS
```

# Dictionaris

|dict| base version | package |
|:---|:---|:---|
|MeCab IPADIC| mecab-ipadic-2.7.0-20070801 | [github.com/ikawaha/kagome-dict-ipa](https://github.com/ikawaha/kagome-dict-ipa)| 
|UniDIC| unidic-mecab-2.1.2_src | [github.com/ikawaha/kagome-dict-uni](https://github.com/ikawaha/kagome-dict-uni) |
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

# Commands

## Install

```shellsession
go get -u github.com/ikawaha/kagome/v2/...
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
   version - show version

tokenize [-file input_file] [-dict dic_file] [-userdict userdic_file] [-sysdict (ipa|uni|ko)] [-simple false] [-mode (normal|search|extended)]
  -dict string
    	dict
  -file string
    	input file
  -mode string
    	tokenize mode (normal|search|extended) (default "normal")
  -simple
    	display abbreviated dictionary contents
  -sysdict string
    	system dict type (ipa|uni|ko) (default "ipa")
  -udict string
    	user dict
```

### Tokenize command

```shellsession
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

### Server command

**API**

Start a server and try to access the "/tokenize" endpoint.

```shellsession
% kagome server &
% curl -XPUT localhost:6060/tokenize -d'{"sentence":"すもももももももものうち", "mode":"normal"}' | jq . 
```

**Web App**

Start a server and access `http://localhost:6060`. 
(To draw a lattice, demo application uses graphviz . You need graphviz installed.)

```shellsession
% kagome server &
```

![Demo](https://raw.githubusercontent.com/wiki/ikawaha/kagome/images/demoapp.gif)

### Lattice command

A debug tool of tokenize process outputs a lattice in graphviz dot format.

```shellsession
% kagome lattice すもももももももものうち | dot -Tpng -o lattice.png
```

![lattice](https://raw.githubusercontent.com/wiki/ikawaha/kagome/images/lattice.png)

# Licence

MIT