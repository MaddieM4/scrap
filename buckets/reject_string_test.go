package buckets

import "testing"

func TestRejectPrefixBucket(t *testing.T) {
	b := RejectPrefixBucket("abc")
	tests := bt_slice{
		bt{"abc", false, "URL is just prefix"},
		bt{"abc123", false, "URL has prefix"},
		bt{"ab", true, "Not all the prefix"},
		bt{"123abc", true, "Does not start with prefix"},
	}
	tests.Run(t, b)
}

func TestRejectExactBucket(t *testing.T) {
	b := RejectExactBucket("abc")
	tests := bt_slice{
		bt{"abc", false, "URL is exact match"},
		bt{"abc123", true, "URL has prefix"},
		bt{"123abc", true, "URL has postfix"},
	}
	tests.Run(t, b)
}
