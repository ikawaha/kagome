Kagome v2
===

Kagome is an open source Japanese morphological analyzer written in pure golang.
The dictionary/statiscal models such as MeCab-IPADIC, UniDic (unidic-mecab), Korean MeCab and so on, be able to embedded in binaries.

### Improvements from v1.

* Dictionaries are now maintained in a separate repository and only the dictionaries you need can be embedded in binaries.
* Brushed up and added several APIs.

If you like kagome, please star the repository! :)

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
|MeCab IPADIC| mecab-ipadic-2.7.0-20070801 | github.com/ikawaha/kagome-dict-ipa| 
|UniDIC| unidic-mecab-2.1.2_src | github.com/ikawaha/kagome-dict-uni |
|Korean MeCab|mecab-ko-dic-2.1.1-20180720 | github.com/ikawaha/kagome-dict-ko|

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

### Command line mode

```shellsession
$ kagome tokenize -h
Usage of tokenize:
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

### Server mode

```shellsession
$ kagome server -h
Usage of server:
  -dict string
    	system dict type (ipa|uni|ko) (default "ipa")
  -http string
    	HTTP service address (default ":6060")
  -userdict string
    	user dict
```

### Lattice mode

```shellsession
$ kagome lattice -h
Usage of lattice:
  -dict string
    	dict type (ipa|uni|ko) (default "ipa")
  -mode string
    	tokenize mode (normal|search|extended) (default "normal")
  -output string
    	output file
  -userDict string
    	user dict
  -v	verbose mode
```

# Licence

MIT