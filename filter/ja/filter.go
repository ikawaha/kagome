package ja

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/ikawaha/kagome/v2/filter"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// Filter represents a japanese token filter.
type Filter struct {
	baserForm *filter.POSFilter
	stopTags  *filter.POSFilter
	stopWords *filter.WordFilter
}

// FilterOption represents an option of the japanese token filter.
type FilterOption func(*Filter)

// BaseFormFilterOption returns a base form filter option.
func BaseFormFilterOption(p []filter.POS) FilterOption {
	return func(f *Filter) {
		f.baserForm = filter.NewPOSFilter(p...)
	}
}

// StopTagsFilterOption returns a stop tags filter option.
func StopTagsFilterOption(p []filter.POS) FilterOption {
	return func(f *Filter) {
		f.stopTags = filter.NewPOSFilter(p...)
	}
}

// StopWordsFilterOption returns a stop words filter option.
func StopWordsFilterOption(p []string) FilterOption {
	return func(f *Filter) {
		f.stopWords = filter.NewWordFilter(p)
	}
}

// NewFilter returns a filter with the settings commonly used in lucene.
// To customize, set the options and overwrite the filter.
func NewFilter(opts ...FilterOption) (*Filter, error) {
	ret, err := newDefaultLuceneFilter()
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		opt(ret)
	}
	return ret, err
}

//go:embed asset/stop_tags.txt
var stopTags []byte

//go:embed asset/stop_words.txt
var stropWords []byte

const (
	posHierarchy      = 4
	defaultPOSFeature = "*"
)

func newDefaultLuceneFilter() (*Filter, error) {
	ta, err := newDefaultLuceneStopTagPOSFilter()
	if err != nil {
		return nil, fmt.Errorf("failed to load stop tags: %w", err)
	}
	wo, err := newDefaultLuceneStopWordFilter()
	if err != nil {
		return nil, fmt.Errorf("failed to load stop words: %w", err)
	}
	return &Filter{
		baserForm: filter.NewPOSFilter(filter.POS{"動詞"}, filter.POS{"形容詞"}, filter.POS{"形容動詞"}),
		stopTags:  ta,
		stopWords: wo,
	}, nil
}

func newDefaultLuceneStopTagPOSFilter() (*filter.POSFilter, error) {
	t, err := loadConfig(stopTags)
	if err != nil {
		return nil, fmt.Errorf("failed to load stop tags: %w", err)
	}
	ps := make([]filter.POS, 0, len(t))
	for _, v := range t {
		pos := strings.Split(v, "-")
		for i := len(pos); i < posHierarchy; i++ {
			pos = append(pos, defaultPOSFeature)
		}
		ps = append(ps, pos)
	}
	return filter.NewPOSFilter(ps...), nil
}

func newDefaultLuceneStopWordFilter() (*filter.WordFilter, error) {
	t, err := loadConfig(stropWords)
	if err != nil {
		return nil, fmt.Errorf("failed to load stop words: %w", err)
	}
	return filter.NewWordFilter(t), nil
}

// Yield returns a filtered word sequence from a token sequence.
func (f Filter) Yield(tokens []tokenizer.Token) []string {
	var ret []string
	for _, v := range tokens {
		if f.stopTags.Match(v.POS()) {
			continue
		}
		if f.stopWords.Match(v.Surface) {
			continue
		}
		if f.baserForm.Match(v.POS()) {
			if b, ok := v.BaseForm(); ok {
				ret = append(ret, b)
			}
			continue
		}
		ret = append(ret, v.Surface)
	}
	return ret
}

// Drop drops a token given the provided match function (stop-tags and stop-words).
func (f Filter) Drop(tokens *[]tokenizer.Token) {
	f.stopTags.Drop(tokens)
	f.stopWords.Drop(tokens)
}
