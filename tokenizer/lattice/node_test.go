package lattice

import (
	"testing"
)

func Test_NodeClassString(t *testing.T) {
	pairs := []struct {
		in  NodeClass
		out string
	}{
		{DUMMY, "DUMMY"},
		{KNOWN, "KNOWN"},
		{UNKNOWN, "UNKNOWN"},
		{USER, "USER"},
		{NodeClass(999), "UNDEF"},
	}

	for _, p := range pairs {
		if p.in.String() != p.out {
			t.Errorf("got %v, expected %v", p.in.String(), p.out)
		}
	}
}
