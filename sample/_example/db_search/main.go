/*
# TL; DR

This example provides a practical example of how to work with Japanese text data and perform efficient full-text search using Kagome and SQLite3.

# TS; WM

In this example, each line of text is inserted into a row of the SQLite3 database, and then the database is searched for the word "人魚" and "人".

Note that the string tokenized by Kagome, a.k.a. "Wakati", is recorded in a separate table for FTS (Full-Text-Search) at the same time as the original text.

This allows Unicode text data that is not separated by spaces, such as Japanese, to be searched by FTS.

Aim of this example:

This example can be useful in scenarios where you need to perform full-text searches on Japanese text. It demonstrates how to tokenize Japanese text using Kagome, which is a common requirement when working with text data in the Japanese language. By using SQLite with FTS4, it efficiently manages and searches through a large amount of text data, making it suitable for applications like:

1. **Search Engines:** You can use this code as a basis for building a search engine that indexes and searches Japanese text content.
2. **Document Management Systems:**	This code can be integrated into a document management system to enable full-text search capabilities for Japanese documents.
3. **Content Recommendation Systems:** When you have a large collection of Japanese content, you can use this code to implement content recommendation systems based on user queries.
4. **Chatbots and NLP:**  If you're building chatbots or natural language processing (NLP) systems for Japanese language, this code can assist in text analysis and search within the chatbot's knowledge base.

Acknowledgements:

This example is taken in part from the following book for reference.

- p.372, 9.2 "データーベース登録プログラム", "Go言語プログラミングエッセンス エンジニア選書"
  - Written by: Mattn
  - Published: 2023/3/9 (技術評論社)
  - ISBN: 4297134195 / 978-4297134198
  - ASIN: B0BVZCJQ4F / https://amazon.co.jp/dp/4297134195
*/
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Contents to be inserted into the database. Each element represents a line
	// of text and will be inserted into a row of the database.
	lines := []string{
		"人魚は、南の方の海にばかり棲んでいるのではありません。",
		"北の海にも棲んでいたのであります。",
		"北方の海の色は、青うございました。",
		"ある時、岩の上に、女の人魚があがって、",
		"あたりの景色を眺めながら休んでいました。",
		"小川未明 『赤い蝋燭と人魚』",
	}

	// Create a database. In-memory database is used for simplicity.
	db, err := sql.Open("sqlite3", ":memory:")
	PanicOnError(err)

	defer db.Close()

	// Create tables.
	// The first table "contents_fts" is for storing the original content, and
	// the second table "fts" is for storing the tokenized content.
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS contents_fts(docid INTEGER PRIMARY KEY AUTOINCREMENT, content TEXT);
		CREATE VIRTUAL TABLE IF NOT EXISTS fts USING fts4(words);
	`)
	PanicOnError(err)

	// Insert contents
	for _, line := range lines {
		rowID, err := insertContent(db, line)
		PanicOnError(err)

		err = insertSearchToken(db, rowID, line)
		PanicOnError(err)
	}

	// Search by word.
	// Note the difference between words and character search.
	for _, searchWord := range []string{
		// "人魚" should be found at line 1, 4, 6.
		"人魚",
		// "人" exists as a character but should not be found since it is not used
		// as a word.
		"人",
		// "北方" should be found at line 3.
		"北方",
		// "北" should be found at line 2.
		// The character "北" itself exists in line 3 as well, but it is not used
		// as a word. Therefore, line 3 should not match.
		"北",
	} {
		fmt.Println("Searching for:", searchWord)

		rowIDsFound, err := searchFTS4(db, searchWord)
		PanicOnError(err)

		if len(rowIDsFound) == 0 {
			fmt.Println("  No results found")
			continue
		}

		// Print search results
		for _, rowID := range rowIDsFound {
			cont, err := retrieveContent(db, rowID)
			PanicOnError(err)

			fmt.Printf("  Found content: %s at line: %v\n", cont, rowID)
		}
	}
	// Output:
	// Searching for: 人魚
	//   Found content: 人魚は、南の方の海にばかり棲んでいるのではありません。 at line: 1
	//   Found content: ある時、岩の上に、女の人魚があがって、 at line: 4
	//   Found content: 小川未明 『赤い蝋燭と人魚』 at line: 6
	// Searching for: 人
	//   No results found
	// Searching for: 北方
	//   Found content: 北方の海の色は、青うございました。 at line: 3
	// Searching for: 北
	//   Found content: 北の海にも棲んでいたのであります。 at line: 2
}

func insertContent(db *sql.DB, content string) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO contents_fts(content) VALUES(?)`,
		content,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func insertSearchToken(db *sql.DB, rowID int64, content string) error {
	// This example uses the IPA dictionary, but it may be more efficient to use
	// the 'Uni' dictionary if memory is available.
	tknzr, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	PanicOnError(err)

	seg := tknzr.Wakati(content)
	tokenizedContent := strings.Join(seg, " ")

	_, err = db.Exec(
		`INSERT INTO fts(docid, words) VALUES(?, ?)`,
		rowID,
		tokenizedContent,
	)

	return err
}

// PanicOnError exits the program with panic if the given error is not nil.
func PanicOnError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func retrieveContent(db *sql.DB, rowID int) (string, error) {
	rows, err := db.Query(
		`SELECT rowid, content FROM contents_fts WHERE rowid=?`,
		rowID,
	)
	if err != nil {
		return "", err
	}

	defer rows.Close()

	for rows.Next() {
		var foundID int
		var content string

		err := rows.Scan(&foundID, &content)
		if err != nil {
			return "", err
		}

		if foundID == rowID {
			return content, nil
		}
	}

	return "", errors.New("no content found")
}

func searchFTS4(db *sql.DB, searchWord string) ([]int, error) {
	rows, err := db.Query(`SELECT rowid, words FROM fts WHERE words MATCH ?`, searchWord)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var lineIDs []int

	for rows.Next() {
		var lineID int
		var words string

		if err := rows.Scan(&lineID, &words); err != nil {
			return nil, err
		}

		// Debug
		// fmt.Printf("- Table: fts, RowID: %d, Value: %s\n", lineID, words)

		lineIDs = append(lineIDs, lineID)
	}

	return lineIDs, nil
}
