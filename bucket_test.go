package scrap

import "testing"

func TestCountBucket_Check(t *testing.T) {
	url := "/foo"
	b := NewCountBucket(0)
	comparem(t, false, b.Check(url), "Should fail with MaxHits = 0")

	b.MaxHits = 1
	comparem(t, true, b.Check(url), "Failed check should not count")
	comparem(t, false, b.Check(url), "Already had 1 hit")

	// Ensure that previous hits were not forgotten
	b.MaxHits = 3
	comparem(t, true, b.Check(url), "Should get 2 more successes")
	comparem(t, true, b.Check(url), "Should get 1 more success")
	comparem(t, false, b.Check(url), "Already had 3 hits")

	// Rewind to 1, should still fail
	b.MaxHits = 1
	comparem(t, false, b.Check(url), "Has more hits than MaxHits")

	url = "/bar"
	comparem(t, true, b.Check(url), "Other URL has its own count")
	comparem(t, false, b.Check(url), "Other URL maxes out")
}
