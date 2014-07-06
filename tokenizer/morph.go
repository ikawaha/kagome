package tokenizer

import (
	"fmt"
	"strings"

	"github.com/ikawaha/kagome/dic"
)

type Morph struct {
	Id         int
	Class      NodeClass
	Start, End int
	Surface    string
	ex         *dic.UserDicContent
}

func (this Morph) Content() (content dic.Content, err error) {
	switch this.Class {
	case DUMMY:
		return
	case KNOWN:
		content = dic.Contents[this.Id]
	case UNKNOWN:
		content = dic.UnkContents[this.Id]
	case USER:
		content.Pos = this.ex.Pos
		content.Yomi = strings.Join(this.ex.Yomi, "")
	}
	return
}

func (this Morph) String() string {
	return fmt.Sprintf("%v(%v, %v)", this.Surface, this.Start, this.End)
}
