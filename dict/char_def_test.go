package dict

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_CharDefWriteToReadCharDef(t *testing.T) {
	def := CharDef{
		CharClass:    CharClass{"class1", "class2", "class3"},
		CharCategory: CharCategory{'a', 'b', 'c'},
		InvokeList:   InvokeList{true, false, true},
		GroupList:    GroupList{false, true, false},
	}

	var b bytes.Buffer
	_, err := def.WriteTo(&b)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	got, err := ReadCharDef(&b)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	if !reflect.DeepEqual(&def, got) {
		t.Errorf("want %+v, got %+v", def, got)
	}
}
