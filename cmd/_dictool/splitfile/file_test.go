package splitfile

import (
	"testing"
)

func TestSplitFileWrite(t *testing.T) {
	var (
		name = "zzz"
		body = []byte("1234567890abcあいう")
	)
	f, err := Open(name, 3)
	if err != nil {
		t.Fatal(err)
	}
	f.Write(body)
	f.Close()
}
