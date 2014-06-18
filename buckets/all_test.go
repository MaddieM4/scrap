package buckets

import (
	"github.com/campadrenalin/scrap"
	"testing"
)

func TestAllBucket_Empty(t *testing.T) {
	b := AllBucket{}
	if !b.Check("foo") {
		t.Fatal("Should always return true when there are no children")
	}
}

func TestAllBucket_OneItem(t *testing.T) {
	b := AllBucket{
		Children: []scrap.Bucket{
			scrap.NewCountBucket(1),
		},
	}
	tests := bt_slice{
		bt{"foo", true, "First should succeed"},
		bt{"foo", false, "Second should fail"},
	}
	tests.Run(t, b)
}

func TestAllBucket_MultipleItems(t *testing.T) {
	b := AllBucket{
		Children: []scrap.Bucket{
			scrap.NewCountBucket(1),
			scrap.NewCountBucket(2),
		},
	}
	tests := bt_slice{
		bt{"foo", true, "First should succeed"},
		bt{"foo", false, "Second should fail"},
	}
	tests.Run(t, b)

	// Adjust the lower MaxHits. Second CountBucket should have one
	// hit left - then we get falses again.
	b.Children[0].(*scrap.CountBucket).SetMaxHits(4)
	tests = bt_slice{
		bt{"foo", true, "Third should succeed"},
		bt{"foo", false, "Fourth should fail"},
	}
	tests.Run(t, b)
}
