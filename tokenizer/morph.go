package tokenizer

import (
	"github.com/ikawaha/kagome/dic"

	"fmt"
)

type Morph struct {
	id         int
	class      NodeType
	Start, End int
	Surface    string
}

func (this Morph) Content() (content dic.Content, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Morph.Content(): %v, %v", e.(error), this)
		}
	}()
	if this.id == BOSEOS {
		return
	}
	switch this.class {
	case KNOWN:
		content = dic.Contents[this.id]
	case UNKNOWN:
		content = dic.UnkContents[this.id]
	}
	return
}

func (this Morph) String() string {
	return fmt.Sprintf("%v(%v, %v)", this.Surface, this.Start, this.End)
}
