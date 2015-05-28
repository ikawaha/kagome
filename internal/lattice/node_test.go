// Copyright 2015 ikawaha
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// 	You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lattice

import (
	"testing"
)

func TestNodeClassString(t *testing.T) {

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
