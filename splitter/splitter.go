package splitter

import (
	"unicode/utf8"
)

func isSpace(r rune) bool {
	if r <= '\u00FF' {
		// Obvious ASCII ones: \t through \r plus space. Plus two Latin-1 oddballs.
		switch r {
		case ' ', '\t', '\n', '\v', '\f', '\r':
			return true
		case '\u0085', '\u00A0':
			return true
		}
		return false
	}
	// High-valued ones.
	if '\u2000' <= r && r <= '\u200a' {
		return true
	}
	switch r {
	case '\u1680', '\u2028', '\u2029', '\u202f', '\u205f', '\u3000':
		return true
	}
	return false
}

type Spliter struct {
	Delim               []rune
	Follower            []rune
	SkipWhiteSpace      bool
	DoubleLineFeedSplit bool
	MaxRuneLen          int
}

var (
	spliter = &Spliter{
		Delim:               []rune{'。', '．'},
		Follower:            []rune{'」', '』'},
		SkipWhiteSpace:      true,
		DoubleLineFeedSplit: true,
		MaxRuneLen:          256,
	}
)

func ScanSentences(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return spliter.ScanSentences(data, atEOF)
}

func (s Spliter) isDelim(r rune) bool {
	for _, d := range s.Delim {
		if r == d {
			return true
		}
	}
	return false
}

func (s Spliter) isFollower(r rune) bool {
	for _, d := range s.Follower {
		if r == d {
			return true
		}
	}
	return false
}

func (s Spliter) ScanSentences(data []byte, atEOF bool) (advance int, token []byte, err error) {
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
		if s.SkipWhiteSpace && isSpace(r) {
			p += size
			if head {
				start, end = p, p
			}
			if s.DoubleLineFeedSplit && r == '\n' {
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
			if s.SkipWhiteSpace && isSpace(r) {
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
