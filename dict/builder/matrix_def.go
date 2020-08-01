package builder

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MatrixDef represents matrix.def.
type MatrixDef struct {
	rowSize int64
	colSize int64
	vec     []int16
}

func parseMatrixDefFile(path string) (*MatrixDef, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	dim := strings.Split(line, " ")
	if len(dim) != 2 {
		return nil, fmt.Errorf("invalid format: %s", line)
	}
	rowSize, err := strconv.ParseInt(dim[0], 10, 0)
	if err != nil {
		return nil, fmt.Errorf("invalid format: %s, %s", err, line)
	}
	colSize, err := strconv.ParseInt(dim[1], 10, 0)
	if err != nil {
		return nil, fmt.Errorf("invalid format: %s, %s", err, line)
	}
	vec := make([]int16, rowSize*colSize)
	for scanner.Scan() {
		line := scanner.Text()
		ary := strings.Split(line, " ")
		if len(ary) != 3 {
			return nil, fmt.Errorf("invalid format: %s", line)
		}
		row, err := strconv.ParseInt(ary[0], 10, 0)
		if err != nil {
			return nil, fmt.Errorf("invalid format: %s, %s", err, line)
		}
		col, err := strconv.ParseInt(ary[1], 10, 0)
		if err != nil {
			return nil, fmt.Errorf("invalid format: %s, %s", err, line)
		}
		val, err := strconv.Atoi(ary[2])
		if err != nil {
			return nil, fmt.Errorf("invalid format: %s, %s", err, line)
		}
		vec[row*colSize+col] = int16(val)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("invalid format: %s, %s", err, line)
	}
	return &MatrixDef{
		rowSize: rowSize,
		colSize: colSize,
		vec:     vec,
	}, nil
}
