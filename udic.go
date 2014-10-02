package kagome

import (
	"bufio"
	"fmt"
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
	Index    Trie
	Contents []UserDicContent
}

// NewUserDic build a user dictionary from a file.
func NewUserDic(path string) (udic *UserDic, err error) {
	const userDicColumnSize = 4
	var file *os.File
	file, err = os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	var text []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		text = append(text, line)
	}
	if e := scanner.Err(); e != nil {
		err = e
		return
	}

	sort.Strings(text)

	udic = new(UserDic)
	prev := ""
	var keys []string
	var ids []int
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
		ids = append(ids, len(keys))
		keys = append(keys, k)
		tokens := strings.Split(record[1], " ")
		yomi := strings.Split(record[2], " ")
		if len(tokens) == 0 || len(tokens) != len(yomi) {
			err = fmt.Errorf("invalid format: %s", line)
			return
		}
		udic.Contents = append(udic.Contents, UserDicContent{tokens, yomi, record[3]})
	}
	da := new(DoubleArray)
	da.BuildWithIds(keys, ids)
	udic.Index = da
	return
}
