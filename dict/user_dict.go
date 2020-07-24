package dict

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// UserDictColumnSize is the column size of the user dictionary.
const UserDictColumnSize = 4

// UserDictContent represents contents of a word in a user dictionary.
type UserDictContent struct {
	Tokens []string
	Yomi   []string
	Pos    string
}

// UserDict represents a user dictionary.
type UserDict struct {
	Index    IndexTable
	Contents []UserDictContent
}

// NewUserDict build a user dictionary from a file.
func NewUserDict(path string) (*UserDict, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, err := NewUserDicRecords(f)
	if err != nil {
		return nil, err
	}
	return r.NewUserDict()
}

// UserDicRecord represents a record of the user dictionary file format.
type UserDicRecord struct {
	Text   string   `json:"text"`
	Tokens []string `json:"tokens"`
	Yomi   []string `json:"yomi"`
	Pos    string   `json:"pos"`
}

// UserDictRecords represents user dictionary data.
type UserDictRecords []UserDicRecord

func (u UserDictRecords) Len() int           { return len(u) }
func (u UserDictRecords) Swap(i, j int)      { u[i], u[j] = u[j], u[i] }
func (u UserDictRecords) Less(i, j int) bool { return u[i].Text < u[j].Text }

// NewUserDicRecords loads user dictionary data from io.Reader.
func NewUserDicRecords(r io.Reader) (UserDictRecords, error) {
	var text []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		text = append(text, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var records UserDictRecords
	for _, line := range text {
		vec := strings.Split(line, ",")
		if len(vec) != UserDictColumnSize {
			return nil, fmt.Errorf("invalid format: %s", line)
		}
		tokens := strings.Split(vec[1], " ")
		yomi := strings.Split(vec[2], " ")
		if len(tokens) == 0 || len(tokens) != len(yomi) {
			return nil, fmt.Errorf("invalid format: %s", line)
		}
		r := UserDicRecord{
			Text:   vec[0],
			Tokens: tokens,
			Yomi:   yomi,
			Pos:    vec[3],
		}
		records = append(records, r)
	}
	return records, nil
}

// NewUserDict builds a user dictionary.
func (u UserDictRecords) NewUserDict() (*UserDict, error) {
	sort.Sort(u)
	prev := ""
	keys := make([]string, 0, len(u))
	var ret UserDict
	for _, r := range u {
		k := strings.TrimSpace(r.Text)
		if prev == k {
			return nil, fmt.Errorf("duplicated error, %+v", r)
		}
		prev = k
		keys = append(keys, k)
		if len(r.Tokens) == 0 || len(r.Tokens) != len(r.Yomi) {
			return nil, fmt.Errorf("invalid format, %+v", r)
		}
		c := UserDictContent{
			Tokens: r.Tokens,
			Yomi:   r.Yomi,
			Pos:    r.Pos,
		}
		ret.Contents = append(ret.Contents, c)
	}
	idx, err := BuildIndexTable(keys)
	ret.Index = idx
	return &ret, err
}
