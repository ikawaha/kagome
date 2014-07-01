package tokenizer

import (
	"fmt"

	"github.com/ikawaha/kagome/dic"
)

type Morph struct {
	Id         int
	Class      NodeClass
	Start, End int
	Surface    string
}

func (this Morph) Content() (content dic.Content, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Morph.Content(): %v, %v", e.(error), this)
		}
	}()
	switch this.Class {
	case DUMMY:
		return
	case KNOWN:
		content = dic.Contents[this.Id]
	case UNKNOWN:
		content = dic.UnkContents[this.Id]
	}
	return
}

func (this Morph) String() string {
	return fmt.Sprintf("%v(%v, %v)", this.Surface, this.Start, this.End)
}
