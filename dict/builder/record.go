package builder

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

// MorphRecordInfo represents a format of CSV records.
type MorphRecordInfo struct {
	ColSize                 int
	SurfaceIndex            int
	LeftIDIndex             int
	RightIDIndex            int
	WeightIndex             int
	POSStartIndex           int
	OtherContentsStartIndex int

	// extra info.
	Meta map[string]int8
}

// UnkRecordInfo represents a format of unk CSV records.
type UnkRecordInfo struct {
	ColSize                 int
	CategoryIndex           int
	LeftIDIndex             int
	RightIndex              int
	WeigthIndex             int
	POSStartIndex           int
	OtherContentsStartIndex int
}

// Record represents a record of CSV.
type Record []string

// Records represents records of CSV.
type Records []Record

func (p Records) Len() int           { return len(p) }
func (p Records) Less(i, j int) bool { return p[i][0] < p[j][0] }
func (p Records) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func parseCSVFiles(path string, enc encoding.Encoding, colSize int) (Records, error) {
	if dir, err := os.Stat(path); os.IsNotExist(err) || !dir.IsDir() {
		return nil, fmt.Errorf("no such directory, %q", path)
	}
	paths, err := filepath.Glob(path + "/*.csv")
	if err != nil {
		return nil, fmt.Errorf("path expansion error, %v", err)
	}
	var records Records
	for _, v := range paths {
		rec, err := parseCSVFile(v, enc, colSize)
		if err != nil {
			return nil, fmt.Errorf("read error, %q, %v", v, err)
		}
		records = append(records, rec...)
	}
	return records, nil
}

func parseCSVFile(path string, enc encoding.Encoding, colSize int) (Records, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var r io.Reader = f
	if enc != nil {
		r = transform.NewReader(r, enc.NewDecoder())
	}
	cr := csv.NewReader(r)
	cr.Comma = ','

	var ret Records
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if colSize > 0 && len(rec) != colSize {
			return nil, fmt.Errorf("invalid format csv: %v, want col size %d, got %d, %+v", path, colSize, len(ret), rec)
		}
		ret = append(ret, rec)
	}
	return ret, nil
}
