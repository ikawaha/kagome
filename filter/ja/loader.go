package ja

import (
	"bufio"
	"bytes"
	"strings"
)

func loadConfig(b []byte) ([]string, error) {
	s := bufio.NewScanner(bytes.NewReader(b))
	var ret []string
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		if i := strings.Index(line, "#"); i > 0 {
			line = line[:i]
		}
		if line == "" {
			continue
		}
		ret = append(ret, line)
	}
	return ret, s.Err()
}
