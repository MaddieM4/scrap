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
	if !b.Check("foo") {
		t.Fatal("First should succeed")
	}
	if b.Check("foo") {
		t.Fatal("Second should fail")
	}
}

func TestAllBucket_MultipleItems(t *testing.T) {
	b := AllBucket{
		Children: []scrap.Bucket{
			scrap.NewCountBucket(1),
			scrap.NewCountBucket(2),
		},
	}
	if !b.Check("foo") {
		t.Fatal("First should succeed")
	}
	if b.Check("foo") {
		t.Fatal("Second should fail")
	}

	// Adjust the lower MaxHits. Second CountBucket should have one
	// hit left - then we get falses again.
	b.Children[0].(*scrap.CountBucket).SetMaxHits(4)
	if !b.Check("foo") {
		t.Fatal("Third should succeed")
	}
	if b.Check("foo") {
		t.Fatal("Fourth should fail")
	}
}
