# Full-text search with Kagome and SQLite3

This example provides a practical example of how to work with Japanese text data and **perform efficient [full-text search](https://en.wikipedia.org/wiki/Full-text_search) using Kagome and SQLite3**.

- Target text data is as follows:

```text
人魚は、南の方の海にばかり棲んでいるのではありません。
北の海にも棲んでいたのであります。
北方の海の色は、青うございました。
ある時、岩の上に、女の人魚があがって、
あたりの景色を眺めながら休んでいました。
小川未明 『赤い蝋燭と人魚』
```

- Example output:

```shellsession
$ cd /path/to/kagome/_examples/db_search
$ go run .
Searching for: 人魚
  Found content: 人魚は、南の方の海にばかり棲んでいるのではありません。 at line: 1
  Found content: ある時、岩の上に、女の人魚があがって、 at line: 4
  Found content: 小川未明 『赤い蝋燭と人魚』 at line: 6
Searching for: 人
  No results found
Searching for: 北方
  Found content: 北方の海の色は、青うございました。 at line: 3
Searching for: 北
  Found content: 北の海にも棲んでいたのであります。 at line: 2
```

- [View main.go](main.go)

## Details

In this example, each line of text is inserted into a row of the SQLite3 database, and then the database is searched for the word "人魚", "人", "北方" and "北".

When inserting text data into the database, Kagome is used to tokenize the text into words.

The string (or a line) tokenized by Kagome, a.k.a. "Wakati", is recorded in a separate table for [FTS4](https://www.sqlite.org/fts3.html) (Full-Text-Search) relative to the original text.

This allows Unicode text data that is not separated by spaces, such as Japanese, to be searched by FTS.

Note that it is searching by word and not by character. For example "人" doesn't match "人魚". Likewise, "北" doesn't match "北方".

This is due to the fact that the FTS4 module in SQLite3 is designed to search for words, not characters.

### Aim of this example

This example can be useful in scenarios where you need to perform full-text searches on Japanese text.

It demonstrates how to tokenize Japanese text using Kagome, which is a common requirement when working with text data in the Japanese language.

By using SQLite with FTS4, it efficiently manages and searches through a large amount of text data, making it suitable for applications like:

1. **Search Engines:** You can use this code as a basis for building a search engine that indexes and searches Japanese text content.
2. **Document Management Systems:** This code can be integrated into a document management system to enable full-text search capabilities for Japanese documents.
3. **Content Recommendation Systems:** When you have a large collection of Japanese content, you can use this code to implement content recommendation systems based on user queries.
4. **Chatbots and NLP:**  If you're building chatbots or natural language processing (NLP) systems for Japanese language, this code can assist in text analysis and search within the chatbot's knowledge base.

## Acknowledgements

This example is taken in part from the following book for reference.

- p.204, 9.2 "データーベース登録プログラム", "Go言語プログラミングエッセンス エンジニア選書"
  - Written by: [Mattn](https://github.com/mattn)
  - Published: 2023/3/9 (技術評論社)
  - ISBN: 4297134195 / 978-4297134198
  - ASIN: B0BVZCJQ4F / [https://amazon.co.jp/dp/4297134195](https://amazon.co.jp/dp/4297134195)
  - Original sample code: [https://github.com/mattn/aozora-search](https://github.com/mattn/aozora-search)
