[![Build Status](https://travis-ci.org/ikawaha/kagome.svg?branch=master)](https://travis-ci.org/ikawaha/kagome) [![Coverage Status](https://coveralls.io/repos/ikawaha/kagome/badge.png?branch=master)](https://coveralls.io/r/ikawaha/kagome?branch=master) [![GoDoc](https://godoc.org/github.com/ikawaha/kagome?status.svg)](https://godoc.org/github.com/ikawaha/kagome)

Kagome Japanese Morphological Analyzer
===

Kagome(籠目)は Pure Go な日本語形態素解析器です．辞書をソースにエンコードして同梱しているので，バイナリだけで動作します．辞書データとして，MeCab-IPADICを利用しています．

```
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

Install
---

### Source

```
% go get github.com/ikawaha/kagome/...
```

Usage
---
### 形態素解析
```
% kagome -h
usage: kagome [-f input_file] [-u userdic_file]
  -f="": input file
  -u="": user dic
```

ユーザ辞書の形式は kuromoji 形式です．`_sample`にサンプルがあります．
```
第68代横綱朝青龍
第	接頭詞,数接続,*,*,*,*,第,ダイ,ダイ
68	名詞,数,*,*,*,*,*
代	名詞,接尾,助数詞,*,*,*,代,ダイ,ダイ
横綱	名詞,一般,*,*,*,*,横綱,ヨコヅナ,ヨコズナ
朝青龍	カスタム人名,朝青龍,アサショウリュウ
EOS
```

### 解析状況確認
lattice ツールを利用すると，解析状況を [graphviz](http://www.graphviz.org/) の dot形式で出力することができます．グラフ化には graphviz のインストールが別途必要です．
```
$ lattice -v すもももももももものうち  |dot -Tpng -o lattice.png
すもも	名詞,一般,*,*,*,*,すもも,スモモ,スモモ
も	助詞,係助詞,*,*,*,*,も,モ,モ
もも	名詞,一般,*,*,*,*,もも,モモ,モモ
も	助詞,係助詞,*,*,*,*,も,モ,モ
もも	名詞,一般,*,*,*,*,もも,モモ,モモ
の	助詞,連体化,*,*,*,*,の,ノ,ノ
うち	名詞,非自立,副詞可能,*,*,*,うち,ウチ,ウチ
EOS
```
![lattice](https://raw.githubusercontent.com/wiki/ikawaha/kagome/images/lattice.png)

License
---
Kagome is licensed under the Apache License v2.0 and uses the MeCab-IPADIC dictionary/statistical model. See NOTICE.txt for license details. 

TODO
---
* 検索用モードの実装
* API 整備
