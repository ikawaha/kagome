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

package fst

// Pair implements a pair of input and output.
type Pair struct {
	In  string
	Out int32
}

// PairSlice implements a slice of input and output pairs.
type PairSlice []Pair

func (ps PairSlice) Len() int      { return len(ps) }
func (ps PairSlice) Swap(i, j int) { ps[i], ps[j] = ps[j], ps[i] }
func (ps PairSlice) Less(i, j int) bool {
	if ps[i].In == ps[j].In {
		return ps[i].Out < ps[j].Out
	}
	return ps[i].In < ps[j].In
}

func (ps PairSlice) maxInputWordLen() (max int) {
	for _, pair := range ps {
		if size := len(pair.In); size > max {
			max = size
		}
	}
	return
}
