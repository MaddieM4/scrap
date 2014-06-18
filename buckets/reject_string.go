package buckets

import "strings"

type RejectPrefixBucket string

func (b RejectPrefixBucket) Check(url string) bool {
	return !strings.HasPrefix(url, string(b))
}

type RejectExactBucket string

func (b RejectExactBucket) Check(url string) bool {
	return url != string(b)
}
