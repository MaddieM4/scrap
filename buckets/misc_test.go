package buckets

import (
	"github.com/campadrenalin/scrap"
	"testing"
)

// Bucket test. Tiny name for mass array construction.
type bt struct {
	Url      string
	Expected bool
	Info     string
}

type bt_slice []bt

func (bts bt_slice) Run(t *testing.T, b scrap.Bucket) {
	for _, test := range bts {
		got := b.Check(test.Url)
		if got != test.Expected {
			t.Fatalf("%s\n(expected %v, got %v)",
				test.Info, test.Expected, got)
		}
	}
}
