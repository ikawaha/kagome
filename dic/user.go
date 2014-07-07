package dic

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/ikawaha/kagome/trie/da"
)

const (
	_USER_DIC_COLUMN_SIZE = 4
)

type UserDicContent struct {
	Surface, Yomi []string
	Pos           string
}

type UserDic struct {
	Index    *da.DoubleArray
	Contents []UserDicContent
}

func NewUserDic(a_file *os.File) (dic *UserDic, err error) {
	text := make([]string, 0)
	scanner := bufio.NewScanner(a_file)
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
	keys := make([]string, 0)
	dic = new(UserDic)
	dic.Contents = make([]UserDicContent, 0)
	prev := ""
	for _, line := range text {
		record := strings.Split(line, ",")
		if len(record) != _USER_DIC_COLUMN_SIZE {
			err = fmt.Errorf("invalid format: %s", line)
			return
		}
		k := strings.TrimSpace(record[0])
		if prev == k {
			continue
		}
		prev = k
		keys = append(keys, k)
		surface := strings.Split(record[1], " ")
		yomi := strings.Split(record[2], " ")
		if len(surface) == 0 || len(surface) != len(yomi) {
			err = fmt.Errorf("invalid format: %s", line)
			return
		}
		dic.Contents = append(dic.Contents, UserDicContent{surface, yomi, record[3]})
	}
	dic.Index = da.NewDoubleArray()
	dic.Index.Build(keys)
	return
}
