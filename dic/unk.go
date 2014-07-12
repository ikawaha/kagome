package dic

import (
	"bytes"
	"encoding/gob"

	"github.com/ikawaha/kagome/data"
)

var (
	UnkContents []Content
	UnkCosts    []Cost
	UnkIndex    map[CharacterClass][2]int
)

func init() {
	vec, err := data.Asset("data/unk.dic")
	if err != nil {
		panic(err)
	}
	decorder := gob.NewDecoder(bytes.NewBuffer(vec))
	if err = decorder.Decode(&UnkContents); err != nil {
		panic(err)
	}
	if err = decorder.Decode(&UnkCosts); err != nil {
		panic(err)
	}
	if err = decorder.Decode(&UnkIndex); err != nil {
		panic(err)
	}
}

func ReferUnkContent(a_id int) (content Content, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	return UnkContents[a_id], nil
}

func GetUnkCost(a_id int) (cost Cost, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	cost = UnkCosts[a_id]
	return
}
