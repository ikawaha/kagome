Kagome Japanese Morphological Analyzer
===

Kagome(籠目)は Pure Go な日本語形態素解析器のプロトタイプです．辞書をソースにエンコードして同梱しているので，バイナリだけで動作します．
辞書データとして，MeCab-IPADICを利用しています．


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
~~ソースの容量が大きく，コンパイル時間もかなりかかりますのでご注意ください．~~ 改善しました

```
% go get github.com/ikawaha/kagome/...
```

Usage
---

```
% kagome -h
usage: kagome [-f input_file] [-u userdic_file]
  -f="": input file
  -u="": input file
```

ユーザ辞書の形式は kuromoji 形式です．`_sample`にサンプルがあります．
```
% kagome -u _sample/userdic.txt
第68代横綱朝青龍
第	接頭詞,数接続,*,*,*,*,第,ダイ,ダイ
68	名詞,数,*,*,*,*,*,*,*
代	名詞,接尾,助数詞,*,*,*,代,ダイ,ダイ
横綱	名詞,一般,*,*,*,*,横綱,ヨコヅナ,ヨコズナ
朝青龍	カスタム人名,*,*,*,*,*,*,アサショウリュウ,*
EOS
```

License
---
Kagome is licensed under the Apache License v2.0 and uses the MeCab-IPADIC dictionary/statistical model. See NOTICE.txt for license details. 

TODO
---
* Kuromoji like search mode
