package builder

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/ikawaha/kagome/v2/dict"
	"golang.org/x/text/encoding"
)

func parseUnkDefFile(path string, enc encoding.Encoding, info *UnkRecordInfo, charClass dict.CharClass) (*dict.UnkDict, error) {
	records, err := parseCSVFile(path, enc, info.ColSize)
	if err != nil {
		return nil, err
	}
	ret := dict.UnkDict{
		Index:    map[int32]int32{},
		IndexDup: map[int32]int32{},
		ContentsMeta: dict.ContentsMeta{
			dict.POSStartIndex: int8(info.POSStartIndex - info.POSStartIndex),
			dict.POSHierarchy:  int8(info.OtherContentsStartIndex - info.POSStartIndex),
		},
	}
	sort.Sort(Records(records))
	for _, rec := range records {
		categoryID := int32(-1)
		for id, cat := range charClass {
			if cat == rec[info.CategoryIndex] {
				categoryID = int32(id)
				break
			}
		}
		if categoryID < 0 {
			return nil, fmt.Errorf("unknown unk category: %v", rec[info.CategoryIndex])
		}
		if _, ok := ret.Index[categoryID]; !ok {
			ret.Index[categoryID] = int32(len(ret.Contents))
		} else {
			ret.IndexDup[categoryID]++
		}
		l, err := strconv.Atoi(rec[info.LeftIDIndex])
		if err != nil {
			return nil, err
		}
		if l > MaxInt16 {
			return nil, fmt.Errorf("unk left ID %d > %d, record: %v", l, MaxInt16, rec)
		}
		r, err := strconv.Atoi(rec[info.RightIndex])
		if err != nil {
			return nil, err
		}
		if r > MaxInt16 {
			return nil, fmt.Errorf("unk right ID %d > %d, record: %v", r, MaxInt16, rec)
		}
		w, err := strconv.Atoi(rec[info.WeigthIndex])
		if err != nil {
			return nil, err
		}
		if w > MaxInt16 {
			return nil, fmt.Errorf("unk weight %d > %d, record: %v", w, MaxInt16, rec)
		}
		m := dict.Morph{LeftID: int16(l), RightID: int16(r), Weight: int16(w)}
		ret.Morphs = append(ret.Morphs, m)
		ret.Contents = append(ret.Contents, rec[info.POSStartIndex:])
	}
	return &ret, nil
}
