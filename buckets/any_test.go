package buckets

import (
	"github.com/campadrenalin/scrap"
	"testing"
)

func TestAnyBucket_Empty(t *testing.T) {
	b := AnyBucket{}
	if b.Check("foo") {
		t.Fatal("Should always return false when there are no children")
	}
}

func TestAnyBucket_OneItem(t *testing.T) {
	b := AnyBucket{
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

func TestAnyBucket_MultipleItems(t *testing.T) {
	b := AnyBucket{
		Children: []scrap.Bucket{
			scrap.NewCountBucket(1),
			scrap.NewCountBucket(2),
		},
	}
	tests := bt_slice{
		bt{"foo", true, "First bucket says yes"},
		bt{"foo", true, "First bucket exhausted, second says yes"},
		bt{"foo", true, "Second bucket says yes for the final time"},
		bt{"foo", false, "Second bucket exhausted"},
	}
	tests.Run(t, b)

	// Adjust the lower MaxHits. Second CountBucket still exhausted.
	b.Children[0].(*scrap.CountBucket).SetMaxHits(3)
	tests = bt_slice{
		bt{"foo", true, "First bucket says yes"},
		bt{"foo", true, "First bucket says yes for the final time"},
		bt{"foo", false, "Both buckets exhausted"},
	}
	tests.Run(t, b)
}
