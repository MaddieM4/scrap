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

type RejectContainsBucket string

func (b RejectContainsBucket) Check(url string) bool {
	return strings.Index(url, string(b)) == -1 // No match
}

type RejectSuffixBucket string

func (b RejectSuffixBucket) Check(url string) bool {
	return !strings.HasSuffix(url, string(b))
}
