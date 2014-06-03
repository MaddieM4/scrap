package scrap

import "testing"

type generator func(string) StringTest
type st_test struct {
	Pattern string
	Url     string
	Match   bool
}

func (stt st_test) Run(t *testing.T, gen generator) {
	st := gen(stt.Pattern)
	got := st(stt.Url)
	compare(t, stt.Match, got)
}

type st_test_slice []st_test

func (st_slice st_test_slice) Run(t *testing.T, gen generator) {
	for _, stt := range st_slice {
		stt.Run(t, gen)
	}
}

func TestSTExact(t *testing.T) {
	st_test_slice{
		st_test{
			Pattern: "a pattern for exact matches",
			Url:     "a pattern for exact matches",
			Match:   true,
		},
		st_test{
			Pattern: "a pattern for exact matches",
			Url:     "a pattern",
			Match:   false,
		},
		st_test{
			Pattern: "a pattern for exact matches",
			Url:     "exact matches",
			Match:   false,
		},
		st_test{
			Pattern: "a pattern for exact matches",
			Url:     "a pattern for exact matches, with some at the end",
			Match:   false,
		},
		st_test{
			Pattern: "a pattern for exact matches",
			Url:     "Before a pattern for exact matches",
			Match:   false,
		},
		st_test{
			Pattern: "a pattern for exact matches",
			Url:     "A Pattern For Exact Matches",
			Match:   false,
		},
	}.Run(t, StringTestExact)
}

func TestSTPrefix(t *testing.T) {
	st_test_slice{
		st_test{
			Pattern: "a pattern for exact matches",
			Url:     "a pattern for exact matches",
			Match:   true,
		},
		st_test{
			Pattern: "a pattern for exact matches",
			Url:     "a pattern",
			Match:   false,
		},
		st_test{
			Pattern: "a pattern for exact matches",
			Url:     "exact matches",
			Match:   false,
		},
		st_test{
			Pattern: "a pattern for exact matches",
			Url:     "a pattern for exact matches, with some at the end",
			Match:   true,
		},
		st_test{
			Pattern: "a pattern for exact matches",
			Url:     "Before a pattern for exact matches",
			Match:   false,
		},
		st_test{
			Pattern: "a pattern for exact matches",
			Url:     "A Pattern For Exact Matches",
			Match:   false,
		},
	}.Run(t, StringTestPrefix)
}
