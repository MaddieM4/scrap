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
	if !b.Check("foo") {
		t.Fatal("First should succeed")
	}
	if b.Check("foo") {
		t.Fatal("Second should fail")
	}
}

func TestAnyBucket_MultipleItems(t *testing.T) {
	b := AnyBucket{
		Children: []scrap.Bucket{
			scrap.NewCountBucket(1),
			scrap.NewCountBucket(2),
		},
	}
	expected_results := []bool{
		true, true, true, false, // Exhausts both buckets (2+1)
	}
	for attempt, expected := range expected_results {
		got := b.Check("foo")
		if got != expected {
			t.Fatalf("Part 1, attempt %d: Expected %v, got %v",
				attempt, expected, got)
		}
	}

	// Adjust the lower MaxHits. Second CountBucket still exhausted.
	b.Children[0].(*scrap.CountBucket).SetMaxHits(3)
	expected_results = []bool{
		true, true, false, // First bucket has two left now
	}
	for attempt, expected := range expected_results {
		got := b.Check("foo")
		if got != expected {
			t.Fatalf("Part 2, attempt %d: Expected %v, got %v",
				attempt, expected, got)
		}
	}
}
