package dict

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
)

const (
	// MorphDictFileName is the default file name of a morph dict.
	MorphDictFileName = "morph.dict"
	// POSDictFileName is the default file name of a part of speech dict.
	POSDictFileName = "pos.dict"
	// ContentMetaFileName is the default file name of content meta.
	ContentMetaFileName = "content.meta"
	// ContentDictFileName is the default file name of a content dict.
	ContentDictFileName = "content.dict"
	// IndexDictFileName is the default filename of a dictionary index.
	IndexDictFileName = "index.dict"
	// ConnectionDictFileName is the default filename of a connection dict.
	ConnectionDictFileName = "connection.dict"
	// CharDefDictFileName is the default filename of a char def.
	CharDefDictFileName = "chardef.dict"
	// UnkDictFileName is the default filename of an unknown dict.
	UnkDictFileName = "unk.dict"
)

// Dict represents a dictionary of a tokenizer.
type Dict struct {
	Morphs       Morphs
	POSTable     POSTable
	ContentsMeta ContentsMeta
	Contents     Contents
	Connection   ConnectionTable
	Index        IndexTable
	CharClass    CharClass
	CharCategory CharCategory
	InvokeList   InvokeList
	GroupList    GroupList
	UnkDict      UnkDict
}

// CharacterCategory returns the category of a rune.
func (d Dict) CharacterCategory(r rune) byte {
	if int(r) < len(d.CharCategory) {
		return d.CharCategory[r]
	}
	return d.CharCategory[0] // default
}

func (d *Dict) loadMorphsDict(r io.Reader) error {
	m, err := ReadMorphs(r)
	if err != nil {
		return fmt.Errorf("dict initializer, Morphs: %v", err)
	}
	d.Morphs = m
	return nil
}

func (d *Dict) loadPOSDict(r io.Reader) error {
	p, err := ReadPOSTable(r)
	if err != nil {
		return fmt.Errorf("dict initializer, POSs: %v", err)
	}
	d.POSTable = p
	return nil
}

func (d *Dict) loadContentsMeta(r io.Reader) error {
	c, err := ReadContentsMeta(r)
	if err != nil {
		return fmt.Errorf("dict initializer, Contents meta: %v", err)
	}
	d.ContentsMeta = c
	return nil

}

func (d *Dict) loadContentsDict(r io.Reader) error {
	c, err := ReadContents(r)
	if err != nil {
		return fmt.Errorf("dict initializer, Contents: %v", err)
	}
	d.Contents = c
	return nil
}

func (d *Dict) loadIndexDict(r io.Reader) error {
	idx, err := ReadIndexTable(r)
	if err != nil {
		return fmt.Errorf("dict initializer, Index: %v", err)
	}
	d.Index = idx
	return nil
}

func (d *Dict) loadConnectionDict(r io.Reader) error {
	t, err := ReadConnectionTable(r)
	if err != nil {
		return fmt.Errorf("dict initializer, Connection: %v", err)
	}
	d.Connection = t
	return nil
}

func (d *Dict) loadCharDefDict(r io.Reader) error {
	def, err := ReadCharDef(r)
	if err != nil {
		return err
	}
	d.CharClass = def.CharClass
	d.CharCategory = def.CharCategory
	d.InvokeList = def.InvokeList
	d.GroupList = def.GroupList
	return nil
}

func (d *Dict) loadUnkDict(r io.Reader) error {
	unk, err := ReadUnkDic(r)
	if err != nil {
		return fmt.Errorf("dic initializer, UnkDict: %v", err)
	}
	d.UnkDict = unk
	return nil
}

// LoadDictFile loads a dictionary from a file.
func LoadDictFile(path string) (d *Dict, err error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return d, err
	}
	defer r.Close()
	return Load(&r.Reader, true)
}

// LoadShrink loads a dictionary from a file without contents.
func LoadShrink(path string) (d *Dict, err error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return d, err
	}
	defer r.Close()
	return Load(&r.Reader, false)
}

type dictionaryPartLoader func(*Dict, io.Reader) error

var loaders = map[string]dictionaryPartLoader{
	MorphDictFileName:      (*Dict).loadMorphsDict,
	POSDictFileName:        (*Dict).loadPOSDict,
	ContentMetaFileName:    (*Dict).loadContentsMeta,
	ContentDictFileName:    (*Dict).loadContentsDict,
	IndexDictFileName:      (*Dict).loadIndexDict,
	ConnectionDictFileName: (*Dict).loadConnectionDict,
	CharDefDictFileName:    (*Dict).loadCharDefDict,
	UnkDictFileName:        (*Dict).loadUnkDict,
}

// Load loads a dictionary from a zipped reader.
func Load(r *zip.Reader, full bool) (*Dict, error) {
	var d Dict
	for _, f := range r.File {
		if !full && f.Name == ContentDictFileName {
			continue
		}
		if err := loadZippedDictPart(f, &d); err != nil {
			return nil, fmt.Errorf("%q, %v", f.Name, err)
		}
	}
	return &d, nil
}

func loadZippedDictPart(f *zip.File, d *Dict) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	loader, ok := loaders[f.Name]
	if !ok {
		return errors.New("unknown file")
	}
	return loader(d, rc)
}

var dictionaryPartFiles = []string{
	MorphDictFileName,
	POSDictFileName,
	ContentMetaFileName,
	ContentDictFileName,
	IndexDictFileName,
	ConnectionDictFileName,
	CharDefDictFileName,
	UnkDictFileName,
}

type dictionaryPartSaver func(Dict, io.Writer) error

var savers = map[string]dictionaryPartSaver{
	MorphDictFileName:      Dict.saveMorphsDict,
	POSDictFileName:        Dict.savePOSTableDict,
	ContentMetaFileName:    Dict.saveContentsMeta,
	ContentDictFileName:    Dict.saveContentsDict,
	IndexDictFileName:      Dict.saveIndexDict,
	ConnectionDictFileName: Dict.saveConnectionDict,
	CharDefDictFileName:    Dict.saveCharDefDict,
	UnkDictFileName:        Dict.saveUnkDict,
}

// Save saves a dictionary in a zipped format.
func (d Dict) Save(zw *zip.Writer) error {
	for _, f := range dictionaryPartFiles {
		saver, ok := savers[f]
		if !ok {
			return fmt.Errorf("unknown file, %q", f)
		}
		w, err := zw.Create(f)
		if err != nil {
			return fmt.Errorf("create file error, %q, %v", f, err)
		}
		if err := saver(d, w); err != nil {
			return fmt.Errorf("write error, %q, %v", f, err)
		}
	}
	return nil
}

func (d Dict) saveMorphsDict(w io.Writer) error {
	_, err := d.Morphs.WriteTo(w)
	return err
}

func (d Dict) savePOSTableDict(w io.Writer) error {
	_, err := d.POSTable.WriteTo(w)
	return err
}

func (d Dict) saveContentsMeta(w io.Writer) error {
	_, err := d.ContentsMeta.WriteTo(w)
	return err
}

func (d Dict) saveContentsDict(w io.Writer) error {
	_, err := d.Contents.WriteTo(w)
	return err
}

func (d Dict) saveIndexDict(w io.Writer) error {
	_, err := d.Index.WriteTo(w)
	return err
}

func (d Dict) saveConnectionDict(w io.Writer) error {
	_, err := d.Connection.WriteTo(w)
	return err
}

func (d Dict) saveCharDefDict(w io.Writer) error {
	def := CharDef{
		CharClass:    d.CharClass,
		CharCategory: d.CharCategory,
		InvokeList:   d.InvokeList,
		GroupList:    d.GroupList,
	}
	if _, err := def.WriteTo(w); err != nil {
		return fmt.Errorf("save char def error, %v", err)
	}
	return nil
}

func (d Dict) saveUnkDict(w io.Writer) error {
	_, err := d.UnkDict.WriteTo(w)
	return err
}
