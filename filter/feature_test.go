package filter

import (
	"testing"
)

func TestFilter_New(t *testing.T) {
	f := NewFeaturesFilter(
		Features{"hello", "world"},
		Features{"hello", "goodbye"},
		Features{"hello", "goodbye"},
	)
	got := f.String()
	want := `hello
  world
  goodbye
`
	if got != want {
		t.Errorf("got:\n%swant:\n%s", got, want)
	}
}

func TestFilter_NewEmpty(t *testing.T) {
	f := NewFeaturesFilter()
	got := f.String()
	want := ""
	if got != want {
		t.Errorf("got:\n%swant:\n%s", got, want)
	}
}

func TestFilter_Add(t *testing.T) {
	f := FeaturesFilter{}
	f.add([]string{"hello", "world"})
	f.add([]string{"hello", "goodbye"})
	f.add([]string{"hello", "goodbye"})
	// only strong conditions are left.
	f.add([]string{"one", "two", "three"})
	f.add([]string{"one", "two"})         // strong condition
	f.add([]string{"one", "two", "piyo"}) // weak condition
	got := f.String()
	want := `hello
  world
  goodbye
one
  two
`
	if got != want {
		t.Errorf("got:\n%swant:\n%s", got, want)
	}
}

func TestFeaturesFilter_Pass(t *testing.T) {
	f := NewFeaturesFilter(
		Features{"hello", "world"},
		Features{"hello", "goodbye"},
		Features{"aloha"},
		Features{"one", "two", "three"},
		Features{"one", "two", "three", "four"},
		Features{Any, "foo", Any, "baa"},
	)
	testdata := []struct {
		features []string
		want     bool
	}{
		// Match
		{features: []string{"hello", "world"}, want: true},
		{features: []string{"hello", "goodbye"}, want: true},
		{features: []string{"aloha"}, want: true},
		{features: []string{"aloha", "aloha"}, want: true},
		{features: []string{"one", "two", "three"}, want: true},
		{features: []string{"one", "two", "three", "four"}, want: true},
		{features: []string{"one", "two", "three", "four", "five"}, want: true},
		{features: []string{"piyo", "foo", "hoge", "baa"}, want: true},
		{features: []string{"zzz", "foo", "zzz", "baa"}, want: true},
		// NG
		{features: []string{"hello"}, want: false},
		{features: []string{"world"}, want: false},
		{features: []string{"one", "world"}, want: false},
		{features: []string{"one"}, want: false},
		{features: []string{"one", "two"}, want: false},
		{features: []string{"foo", "zzz", "zzz", "baa"}, want: false},
	}
	for _, v := range testdata {
		if got := f.Match(v.features); got != v.want {
			t.Errorf("%+v: want %v, got %v", v.features, v.want, got)
		}
	}
}
