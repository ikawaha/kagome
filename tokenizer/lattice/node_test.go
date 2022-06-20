package lattice

import (
	"testing"
)

func Test_NodeClassString(t *testing.T) {
	testdata := []struct {
		class NodeClass
		want  string
	}{
		{class: DUMMY, want: "DUMMY"},
		{class: KNOWN, want: "KNOWN"},
		{class: UNKNOWN, want: "UNKNOWN"},
		{class: USER, want: "USER"},
		{class: NodeClass(999), want: "UNDEF"},
	}
	for _, p := range testdata {
		if got := p.class.String(); got != p.want {
			t.Errorf("got %v, expected %v", got, p.want)
		}
	}
}
