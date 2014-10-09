package kagome

import (
	"fmt"
	"strings"
)

type Token struct {
	Id      int
	Class   NodeClass
	Start   int
	End     int
	Surface string
	dic     *Dic
	udic    *UserDic
}

func (t Token) Features() (features []string) {
	switch t.Class {
	case DUMMY:
		return
	case KNOWN:
		features = t.dic.Contents[t.Id]
	case UNKNOWN:
		features = sysDic.UnkContents[t.Id]
	case USER:
		// XXX
		pos := t.udic.Contents[t.Id].Pos
		tokens := strings.Join(t.udic.Contents[t.Id].Tokens, "/")
		yomi := strings.Join(t.udic.Contents[t.Id].Yomi, "/")
		features = append(features, pos, tokens, yomi)
	}
	return
}

func (t Token) String() string {
	return fmt.Sprintf("%v(%v, %v)%v", t.Surface, t.Start, t.End, t.Class)
}
