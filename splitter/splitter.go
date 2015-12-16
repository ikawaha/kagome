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

package splitter

import (
	"unicode"
	"unicode/utf8"
)

// SentenceSplitter is a tiny sentence splitter for japanese texts.
type SentenceSplitter struct {
	Delim               []rune // delimiter set. ex. {'。','．'}
	Follower            []rune // allow following after delimiters. ex. {'」','』'}
	SkipWhiteSpace      bool   // eliminate white space or not
	DoubleLineFeedSplit bool   // splite at '\n\n' or not
	MaxRuneLen          int    // max sentence length
}

var (
	// default sentence splitter
	defaultSpliter = &SentenceSplitter{
		Delim:               []rune{'。', '．', '！', '!', '？', '?'},
		Follower:            []rune{'.', '｣', '」', '』', ')', '）', '｝', '}', '〉', '》'},
		SkipWhiteSpace:      true,
		DoubleLineFeedSplit: true,
		MaxRuneLen:          256,
	}
)

// ScanSentences is a split function for a bufio.Scanner that returns each sentece of text.
func ScanSentences(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return defaultSpliter.ScanSentences(data, atEOF)
}

func (s SentenceSplitter) isDelim(r rune) bool {
	for _, d := range s.Delim {
		if r == d {
			return true
		}
	}
	return false
}

func (s SentenceSplitter) isFollower(r rune) bool {
	for _, d := range s.Follower {
		if r == d {
			return true
		}
	}
	return false
}

// ScanSentences is a split function for a Scanner that returns each sentece of text.
func (s SentenceSplitter) ScanSentences(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	var (
		start, end, rcount int
		head, nn           bool
	)
	head = true
	for p := 0; p < len(data); {
		r, size := utf8.DecodeRune(data[p:])
		if s.SkipWhiteSpace && unicode.IsSpace(r) {
			p += size
			if head {
				start, end = p, p
			} else if s.isDelim(r) {
				return p, data[start:end], nil
			} else if s.DoubleLineFeedSplit && r == '\n' {
				if nn {
					return p, data[start:end], nil
				}
				nn = true
			}
			continue
		}
		head, nn = false, false // clear flags
		if end != p {
			for i := 0; i < size; i++ {
				data[end+i] = data[p+i]
			}
		}
		p += size
		end += size
		rcount++
		if !s.isDelim(r) && rcount < s.MaxRuneLen {
			continue
		}
		// split
		nn = false
		for p < len(data) {
			r, size := utf8.DecodeRune(data[p:])
			if s.SkipWhiteSpace && unicode.IsSpace(r) {
				p += size
				if s.DoubleLineFeedSplit && r == '\n' {
					if nn {
						return p, data[start:end], nil
					}
					nn = true
				}
			} else if s.isDelim(r) || s.isFollower(r) {
				if end != p {
					for i := 0; i < size; i++ {
						data[end+i] = data[p+i]
					}
				}
				p += size
				end += size
			} else {
				break
			}
		}
		return p, data[start:end], nil
	}
	if !atEOF {
		// Request more data
		for i := end; i < len(data); i++ {
			data[i] = ' '
		}
		return start, nil, nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	return len(data), data[start:end], nil

}
