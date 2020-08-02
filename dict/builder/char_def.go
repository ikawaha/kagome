package builder

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	reCharClass         = regexp.MustCompile(`^(\w+)\s+(\d+)\s+(\d+)\s+(\d+)`)
	reCharCategory      = regexp.MustCompile(`^(0x[0-9A-F]+)(?:\s+([^#\s]+))(?:\s+([^#\s]+))?`)
	reCharCategoryRange = regexp.MustCompile(`^(0x[0-9A-F]+)..(0x[0-9A-F]+)(?:\s+([^#\s]+))(?:\s+([^#\s]+))?`)
)

// CharClassDef represents char.def.
type CharClassDef struct {
	charClass    []string
	charCategory []byte
	invokeMap    []bool
	groupMap     []bool
}

func parseCharClassDefFile(path string) (*CharClassDef, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret := CharClassDef{
		charClass:    make([]string, 0),
		charCategory: make([]byte, 65536),
	}

	scanner := bufio.NewScanner(file)
	cc2id := make(map[string]byte)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		line = strings.TrimSpace(line)

		if matches := reCharClass.FindStringSubmatch(line); len(matches) > 0 {
			ret.invokeMap = append(ret.invokeMap, matches[2] == "1")
			ret.groupMap = append(ret.groupMap, matches[3] == "1")
			cc2id[matches[1]] = byte(len(ret.charClass))
			ret.charClass = append(ret.charClass, matches[1])
		} else if matches := reCharCategory.FindStringSubmatch(line); len(matches) > 0 {
			ch, _ := strconv.ParseInt(matches[1], 0, 0)
			ret.charCategory[ch] = cc2id[matches[2]]
		} else if matches := reCharCategoryRange.FindStringSubmatch(line); len(matches) > 0 {
			start, _ := strconv.ParseInt(matches[1], 0, 0)
			end, _ := strconv.ParseInt(matches[2], 0, 0)
			for x := start; x <= end; x++ {
				ret.charCategory[x] = cc2id[matches[3]]
			}
		} else {
			return nil, fmt.Errorf("invalid format error: %v", line)
		}
	}
	return &ret, scanner.Err()
}
