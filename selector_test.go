package scrap

import "testing"

type generator func(string) StringTest
type selector_test struct {
	Pattern string
	Url     string
	Match   bool
}

func (st selector_test) Run(t *testing.T, gen generator) {
	selector := gen(st.Pattern)
	got := selector(st.Url)
	compare(t, st.Match, got)
}

type selector_test_slice []selector_test

func (st_slice selector_test_slice) Run(t *testing.T, gen generator) {
	for _, st := range st_slice {
		st.Run(t, gen)
	}
}

func TestSTExact(t *testing.T) {
	selector_test_slice{
		selector_test{
			Pattern: "a pattern for exact matches",
			Url:     "a pattern for exact matches",
			Match:   true,
		},
		selector_test{
			Pattern: "a pattern for exact matches",
			Url:     "a pattern",
			Match:   false,
		},
		selector_test{
			Pattern: "a pattern for exact matches",
			Url:     "exact matches",
			Match:   false,
		},
		selector_test{
			Pattern: "a pattern for exact matches",
			Url:     "a pattern for exact matches, with some at the end",
			Match:   false,
		},
		selector_test{
			Pattern: "a pattern for exact matches",
			Url:     "Before a pattern for exact matches",
			Match:   false,
		},
		selector_test{
			Pattern: "a pattern for exact matches",
			Url:     "A Pattern For Exact Matches",
			Match:   false,
		},
	}.Run(t, StringTestExact)
}

func TestSTPrefix(t *testing.T) {
	selector_test_slice{
		selector_test{
			Pattern: "a pattern for exact matches",
			Url:     "a pattern for exact matches",
			Match:   true,
		},
		selector_test{
			Pattern: "a pattern for exact matches",
			Url:     "a pattern",
			Match:   false,
		},
		selector_test{
			Pattern: "a pattern for exact matches",
			Url:     "exact matches",
			Match:   false,
		},
		selector_test{
			Pattern: "a pattern for exact matches",
			Url:     "a pattern for exact matches, with some at the end",
			Match:   true,
		},
		selector_test{
			Pattern: "a pattern for exact matches",
			Url:     "Before a pattern for exact matches",
			Match:   false,
		},
		selector_test{
			Pattern: "a pattern for exact matches",
			Url:     "A Pattern For Exact Matches",
			Match:   false,
		},
	}.Run(t, StringTestPrefix)
}
