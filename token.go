//  Copyright (c) 2014 ikawaha.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package kagome

import (
	"fmt"
	"strings"
)

// Token represents a morph of a sentence.
type Token struct {
	Id      int
	Class   NodeClass
	Start   int
	End     int
	Surface string
	dic     *Dic
	udic    *UserDic
}

// Features returns contents of a token.
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

// String returns a string representation of a token.
func (t Token) String() string {
	return fmt.Sprintf("%v(%v, %v)%v[%v]", t.Surface, t.Start, t.End, t.Class, t.Id)
}
