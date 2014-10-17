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

```
$ kagome -h
usage: kagome [-f input_file | --http addr] [-u userdic_file]
  -file="": input file
  -http="": HTTP service address (e.g., ':6060')
  -udic="": user dic
```

#### 標準入力，もしくはファイルを指定しての解析
入力ファイルを指定した場合，1行1文として解析します．
ファイルのエンコードは utf8 である必要があります．
ファイルを指定しない場合，標準入力から1行1文として解析します．
```
$ kagome
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

#### 検索用の分割モード

![kuromoji](https://github.com/atilika/kuromoji) の検索用分割モード相当の分割が出来るようになっています．

* 標準　標準の分割
* 検索　ヒューリスティックの適用によって検索に役立つよう細分割
* 拡張　検索モードに加えて未知語を unigram に分割します

|入力内容|標準モード|検索モード|拡張モード|
|:-------|:---------|:---------|:---------|
|関西国際空港|関西国際空港|関西　国際　空港|関西　国際　空港|
|日本経済新聞|日本経済新聞|日本　経済　新聞|日本　経済　新聞|
|シニアソフトウェアエンジニア|シニアソフトウェアエンジニア|シニア　ソフトウェア　エンジニア|シニア　ソフトウェア　エンジニア|
|デジカメを買った|デジカメ　を　買っ　た|デジカメ　を　買っ　た|デ　ジ　カ　メ　を　買っ　た|

#### HTTP service
サーバとして動作させると，以下の2つの機能が利用できます．

##### Web API
`-http `オプションを指定するとWebサーバが立ち上がります．
`localhost` にポート`8080`でサーバを立ち上げた場合，'http://localhost:8080/' に REST でアクセスできます．

```
$ kagome -http=":8080" &
$ curl -XPUT localhost:8080 -d'{"sentence":"すもももももももものうち"}'
{"status":true,"tokens":[{"id":36163,"start":0,"end":3,"surface":"すもも","class":"KNOWN","features":["名詞","一般","*","*","*","*","すもも","スモモ","スモモ"]},{"id":73244,"start":3,"end":4,"surface":"も","class":"KNOWN","features":["助詞","係助詞","*","*","*","*","も","モ","モ"]},{"id":74989,"start":4,"end":6,"surface":"もも","class":"KNOWN","features":["名詞","一般","*","*","*","*","もも","モモ","モモ"]},{"id":73244,"start":6,"end":7,"surface":"も","class":"KNOWN","features":["助詞","係助詞","*","*","*","*","も","モ","モ"]},{"id":74989,"start":7,"end":9,"surface":"もも","class":"KNOWN","features":["名詞","一般","*","*","*","*","もも","モモ","モモ"]},{"id":55829,"start":9,"end":10,"surface":"の","class":"KNOWN","features":["助詞","連体化","*","*","*","*","の","ノ","ノ"]},{"id":8024,"start":10,"end":12,"surface":"うち","class":"KNOWN","features":["名詞","非自立","副詞可能","*","*","*","うち","ウチ","ウチ"]}]}
```

##### 形態素解析デモ
Web サーバを立ち上げた状態で，ブラウザで `/_demo` にアクセスすると，形態素解析のデモ利用できます．
`-http=:8080` を指定した場合，`http://localhost:8080/_demo` になります．Lattice の表示には [graphviz](http://www.graphviz.org/) が必要です．

![lattice](https://raw.githubusercontent.com/wiki/ikawaha/kagome/images/demoapp.png)

#### ユーザー辞書について
ユーザ辞書の形式は kuromoji 形式です．`_sample`にサンプルがあります．
```
% kagome
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
