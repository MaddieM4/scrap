package scrap

// A bucket allows for filtering requests according to previous
// (or parallel) requests. See CountBucket for an example of the
// kind of stuff Buckets can do.
type Bucket interface {
	Check(url string) bool
}

// Allows up to MaxHits requests for each unique URL. Use with
// MaxHits = 1 as a simple deduplicator.
type CountBucket struct {
	MaxHits int
	counts  map[string]int
}

func NewCountBucket(max_hits int) CountBucket {
	return CountBucket{
		MaxHits: max_hits,
		counts:  make(map[string]int),
	}
}

func (b CountBucket) Check(url string) bool {
	// if !ok, hits == 0
	hits, _ := b.counts[url]

	if hits < b.MaxHits {
		// Only increment on success
		b.counts[url] = hits + 1
		return true
	} else {
		return false
	}
}
