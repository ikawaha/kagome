package dict

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
)

// POSTable represents a table for managing part of speeches.
type POSTable struct {
	POSs     []POS
	NameList []string
}

// POSID represents a ID of part of speech.
type POSID uint16

// POS represents a vector of part of speech.
type POS []POSID

const maxPOSID = 1<<16 - 1

// WriteTo saves a POS table.
func (p POSTable) WriteTo(w io.Writer) (int64, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(p.POSs); err != nil {
		return 0, err
	}
	if err := enc.Encode(p.NameList); err != nil {
		return 0, err
	}
	return b.WriteTo(w)
}

// ReadPOSTable loads a POS table.
func ReadPOSTable(r io.Reader) (POSTable, error) {
	ret := POSTable{}
	dec := gob.NewDecoder(r)
	if err := dec.Decode(&ret.POSs); err != nil {
		return ret, fmt.Errorf("POSs read error, %v", err)
	}
	if err := dec.Decode(&ret.NameList); err != nil {
		return ret, fmt.Errorf("name list read error, %v", err)
	}
	return ret, nil
}

// POSMap represents a part of speech control table.
type POSMap map[string]POSID

// Add adds part of speech item to the POS control table and returns it's id.
func (p POSMap) Add(pos []string) POS {
	ret := make(POS, 0, len(pos))
	for _, name := range pos {
		id := p.add(name)
		ret = append(ret, id)
	}
	return ret
}

func (p POSMap) add(pos string) POSID {
	id, ok := p[pos]
	if !ok {
		if len(p) > maxPOSID {
			panic(fmt.Errorf("new POSID overflowed %q %d > %d", pos, len(p), maxPOSID))
		}
		id = POSID(len(p)) + 1
		p[pos] = id
	}
	return id
}

// List returns a list whose index is POS ID and value is its name.
func (p POSMap) List() []string {
	ret := make([]string, len(p)+1)
	for k, v := range p {
		ret[v] = k
	}
	return ret
}
