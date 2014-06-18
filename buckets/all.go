package buckets

import "github.com/campadrenalin/scrap"

// Bucket that contains buckets - only returns true if all sub-buckets
// return true.
//
// Buckets are tested in order, using short-circuiting behavior. This
// means that if an earlier bucket returns false, the ones after it
// will not be tested.
type AllBucket struct {
	Children []scrap.Bucket
}

func (ab AllBucket) Check(url string) bool {
	for _, child := range ab.Children {
		if !child.Check(url) {
			return false
		}
	}
	return true
}
