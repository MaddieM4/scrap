package buckets

import "strings"

type RejectPrefixBucket string

func (b RejectPrefixBucket) Check(url string) bool {
	return !strings.HasPrefix(url, string(b))
}
