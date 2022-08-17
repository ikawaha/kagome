package filter

import (
	"fmt"
	"io"
	"strings"
)

type (
	// Feature represents a feature.
	Feature = string
	// Features represents a vector of features.
	Features = []string
)

// Any represents an arbitrary feature.
const Any Feature = "\x00"

type filterEdge struct {
	val Feature
	to  *filterNode
}

type filterNode struct {
	fanout []*filterEdge
}

func (n filterNode) isLeaf() bool {
	return n.fanout == nil
}

func (n filterNode) has(s Feature) int {
	for i, v := range n.fanout {
		if v.val == Any || v.val == s {
			return i
		}
	}
	return -1
}

// FeaturesFilter represents a filter that filters a vector of features.
type FeaturesFilter struct {
	root filterNode
}

// NewFeaturesFilter returns a features filter.
func NewFeaturesFilter(fs ...Features) *FeaturesFilter {
	var ret FeaturesFilter
	for _, v := range fs {
		ret.add(v)
	}
	return &ret
}

func (f *FeaturesFilter) add(fs Features) {
	n := &f.root
	for i, v := range fs {
		tail := i == len(fs)-1
		if k := n.has(v); k >= 0 {
			n = n.fanout[k].to
			if n.isLeaf() {
				break // a stronger condition exists.
			} else if tail {
				n.fanout = nil // clear week conditions.
				break
			}
			continue
		}
		edge := filterEdge{
			val: v,
		}
		n.fanout = append(n.fanout, &edge)
		edge.to = &filterNode{}
		n = edge.to
	}
}

// Match returns true if a filter matches given features.
func (f *FeaturesFilter) Match(fs Features) bool {
	n := &f.root
	for _, v := range fs {
		if k := n.has(v); k >= 0 {
			n = n.fanout[k].to
			if n.isLeaf() {
				return true
			}
			continue
		}
		return false
	}
	return false
}

// String implements string interface.
func (f *FeaturesFilter) String() string {
	var buf strings.Builder
	filterString(&buf, 0, &f.root)
	return buf.String()
}

func filterString(w io.Writer, indent int, n *filterNode) {
	const space = `  `
	for _, edge := range n.fanout {
		_, _ = fmt.Fprintf(w, "%s%s\n", strings.Repeat(space, indent), edge.val)
		if edge.to != nil {
			filterString(w, indent+1, edge.to)
		}
	}
}
