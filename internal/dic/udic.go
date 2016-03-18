// Copyright 2015 ikawaha
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// 	You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dic

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// UserDicContent represents contents of a word in a user dictionary.
type UserDicContent struct {
	Tokens []string
	Yomi   []string
	Pos    string
}

// UserDic represents a user dictionary.
type UserDic struct {
	Index    IndexTable
	Contents []UserDicContent
}

// NewUserDicFromReader build a user dictionary from io.Reader.
func NewUserDicFromReader(r io.Reader) (udic *UserDic, err error) {
	udic = new(UserDic)
	const (
		userDicColumnSize = 4
	)
	var text []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		text = append(text, line)
	}
	if err = scanner.Err(); err != nil {
		return
	}

	sort.Strings(text)

	prev := ""
	var keys []string
	for _, line := range text {
		record := strings.Split(line, ",")
		if len(record) != userDicColumnSize {
			err = fmt.Errorf("invalid format: %s", line)
			return
		}
		k := strings.TrimSpace(record[0])
		if prev == k {
			continue
		}
		prev = k
		keys = append(keys, k)
		tokens := strings.Split(record[1], " ")
		yomi := strings.Split(record[2], " ")
		if len(tokens) == 0 || len(tokens) != len(yomi) {
			err = fmt.Errorf("invalid format: %s", line)
			return
		}
		udic.Contents = append(udic.Contents, UserDicContent{tokens, yomi, record[3]})
	}
	udic.Index, err = BuildIndexTable(keys)
	return
}

// NewUserDic build a user dictionary from a file.
func NewUserDic(path string) (udic *UserDic, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	return NewUserDicFromReader(f)
}
