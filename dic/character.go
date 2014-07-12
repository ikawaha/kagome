package dic

import (
	"bytes"
	"encoding/gob"

	"github.com/ikawaha/kagome/data"
)

type CharacterClass byte

const (
	DEFAULT      CharacterClass = 0
	SPACE        CharacterClass = 1
	NUMERIC      CharacterClass = 4
	HIRAGANA     CharacterClass = 6
	KATAKANA     CharacterClass = 7
	KANJINUMERIC CharacterClass = 8
	GREEK        CharacterClass = 9
	KANJI        CharacterClass = 2
	SYMBOL       CharacterClass = 3
	ALPHA        CharacterClass = 5
	CYRILLIC     CharacterClass = 10
)

var (
	CharacterCategoryList []CharacterClass
	InvokeList            []bool
	GroupList             []bool
)

func init() {
	vec, err := data.Asset("data/char.dic")
	if err != nil {
		panic(err)
	}
	decorder := gob.NewDecoder(bytes.NewBuffer(vec))
	if err = decorder.Decode(&CharacterCategoryList); err != nil {
		panic(err)
	}
	if err = decorder.Decode(&InvokeList); err != nil {
		panic(err)
	}
	if err = decorder.Decode(&GroupList); err != nil {
		panic(err)
	}
}

func (this CharacterClass) String() string {
	switch this {
	case DEFAULT:
		return "DEFAULT"
	case SPACE:
		return "SPACE"
	case NUMERIC:
		return "NUMERIC"
	case HIRAGANA:
		return "HIRAGANA"
	case KATAKANA:
		return "KATAKANA"
	case KANJINUMERIC:
		return "KANJINUMERIC"
	case GREEK:
		return "GREEK"
	case KANJI:
		return "KANJI"
	case SYMBOL:
		return "SYMBOL"
	case ALPHA:
		return "ALPHA"
	case CYRILLIC:
		return "CYRILLIC"

	}
	return ""
}
