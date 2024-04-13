# Tokenizing Example with Kagome

## Analyzing a Japanese text into words and parts of speech with Kagome

This example demonstrates how to analyzes a sentence (tokenize) and get the part-of-speech (POS) of each word using Kagome.

- Target text data is as follows:

```text
すもももももももものうち
```

- Example output:

```shellsession
$ cd /path/to/kagome/_examples/tokenize
$ go run .
---tokenize---
すもも  名詞,一般,*,*,*,*,すもも,スモモ,スモモ
も      助詞,係助詞,*,*,*,*,も,モ,モ
もも    名詞,一般,*,*,*,*,もも,モモ,モモ
も      助詞,係助詞,*,*,*,*,も,モ,モ
もも    名詞,一般,*,*,*,*,もも,モモ,モモ
の      助詞,連体化,*,*,*,*,の,ノ,ノ
うち    名詞,非自立,副詞可能,*,*,*,うち,ウチ,ウチ
```

> __Note__ that tokenization varies depending on the dictionary used. In this example we use the IPA dictionary.
