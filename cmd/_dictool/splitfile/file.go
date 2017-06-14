package splitfile

import (
	"fmt"
	"os"
)

const SplitFileFormat = "%s.%03x"

type SplitFile struct {
	name  string
	limit int
	count int
	id    int
	file  *os.File
}

func Open(name string, limit int) (*SplitFile, error) {
	if limit <= 0 {
		fmt.Errorf("invalid limit, %d > 0", limit)
	}
	s := &SplitFile{
		name:  name,
		limit: limit,
	}
	err := s.newFile()
	return s, err
}

func (s *SplitFile) newFile() error {
	if s.file != nil {
		if err := s.file.Close(); err != nil {
			return err
		}
	}
	name := fmt.Sprintf(SplitFileFormat, s.name, s.id)
	var err error
	if s.file, err = os.Create(name); err != nil {
		return err
	}
	s.id++
	return nil
}

func min(lhs, rhs int) int {
	if lhs < rhs {
		return lhs
	}
	return rhs
}

func (s *SplitFile) Write(b []byte) (n int, err error) {
	for p := 0; p < len(b); {
		if s.file == nil {
			if err = s.newFile(); err != nil {
				return
			}
		}
		cap := s.limit - s.count
		z := min(cap, len(b)-p)
		nn, err := s.file.Write(b[p : p+z])
		if err != nil {
			return n + nn, err
		}
		n += nn
		p = p + z
		s.count += z
		if s.count == s.limit {
			s.file.Close()
			s.file = nil
			s.count = 0
		}
	}
	return
}

func (s *SplitFile) Close() error {
	return s.file.Close()
}
