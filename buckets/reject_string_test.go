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

func TestRejectContainsBucket(t *testing.T) {
	b := RejectContainsBucket("abc")
	tests := bt_slice{
		bt{"abc", false, "URL is exact match"},
		bt{"abc123", false, "URL has prefix"},
		bt{"123abc", false, "URL has postfix"},
		bt{"123abc456", false, "URL has match in the middle"},
		bt{"123ab456", true, "Not quite a match"},
		bt{"bac", true, "All scrambled"},
	}
	tests.Run(t, b)
}

func TestRejectSuffixBucket(t *testing.T) {
	b := RejectSuffixBucket("abc")
	tests := bt_slice{
		bt{"abc", false, "URL is just suffix"},
		bt{"123abc", false, "Ends with suffix"},
		bt{"abc123", true, "URL has prefix"},
		bt{"ab", true, "Not all the suffix"},
		bt{"123abc456", true, "URL has match in the middle"},
	}
	tests.Run(t, b)
}
