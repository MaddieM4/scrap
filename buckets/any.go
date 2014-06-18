package buckets

import "github.com/campadrenalin/scrap"

// Bucket that contains buckets - returns true if any sub-buckets
// return true.
//
// Buckets are tested in order, using short-circuiting behavior. This
// means that if an earlier bucket returns true, the ones after it
// will not be tested.
type AnyBucket struct {
	Children []scrap.Bucket
}

func (ab AnyBucket) Check(url string) bool {
	for _, child := range ab.Children {
		if child.Check(url) {
			return true
		}
	}
	return false
}
