module kagome/examples/db_search

go 1.19

require (
	github.com/ikawaha/kagome-dict/ipa v1.0.10
	github.com/ikawaha/kagome/v2 v2.9.3
	github.com/mattn/go-sqlite3 v1.14.17
)

require github.com/ikawaha/kagome-dict v1.0.9 // indirect

replace github.com/ikawaha/kagome/v2 => ../../
