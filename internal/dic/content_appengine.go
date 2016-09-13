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

// +build appengine

package dic

import (
	"bytes"
)

// NewContents creates dictionary contents from byte slice
func NewContents(b []byte) [][]string {
	rows := bytes.Split(b, []byte(rowDelimiter))
	m := make([][]string, len(rows))
	for i, r := range rows {
		cols := bytes.Split(r, []byte(colDelimiter))
		m[i] = make([]string, len(cols))
		for j, c := range cols {
			m[i][j] = string(c)
		}
	}
	return m
}
