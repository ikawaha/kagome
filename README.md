Kagome Japanese Morphological Analyzer
===

Pure Go で日本語形態素解析器のプロトタイプです．辞書をソースにエンコードして同梱しているので，バイナリだけで動作します．
辞書データとして，MeCab-IPADICを利用しています．

`_sample`ディレクトリに実行用のサンプルが置いてあります．
ソースの容量が大きく，コンパイル時間もかなりかかりますのでご注意ください．

```
すもももももももものうち
  0, BOSEOS(0, 0)       , , , , , , , ,
  1, すもも(0, 3)       名詞, 一般, *, *, *, *, すもも, スモモ, スモモ
  2, も(3, 4)   助詞, 係助詞, *, *, *, *, も, モ, モ
  3, もも(4, 6) 名詞, 一般, *, *, *, *, もも, モモ, モモ
  4, も(6, 7)   助詞, 係助詞, *, *, *, *, も, モ, モ
  5, もも(7, 9) 名詞, 一般, *, *, *, *, もも, モモ, モモ
  6, の(9, 10)  助詞, 連体化, *, *, *, *, の, ノ, ノ
  7, うち(10, 12)       名詞, 非自立, 副詞可能, *, *, *, うち, ウチ, ウチ
  8, BOSEOS(12, 12)     , , , , , , , ,
```
License
---
Kagome is licensed under the Apache License v2.0 and uses the MeCab-IPADIC dictionary/statistical model. See NOTICE.txt for license details. 
