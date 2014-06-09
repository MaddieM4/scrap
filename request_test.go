package scrap

import "testing"

type ContextualizeUrlTest struct {
	SR     string
	Input  string
	Output string
	Error  bool
}

func TestSR_ContextualizeUrl(t *testing.T) {
	tests := []ContextualizeUrlTest{
		ContextualizeUrlTest{
			"/foo/",
			"/bar/",
			"/bar/",
			false,
		},
		ContextualizeUrlTest{
			"http://hostname/foo/",
			"/bar/",
			"http://hostname/bar/",
			false,
		},
		ContextualizeUrlTest{
			"http://hostname/foo/",
			"https://other.hostname/bar/",
			"https://other.hostname/bar/",
			false,
		},
		ContextualizeUrlTest{
			"",
			"/bar/",
			"/bar/",
			false,
		},
		ContextualizeUrlTest{
			"",
			"https://other.hostname/bar/",
			"https://other.hostname/bar/",
			false,
		},
		ContextualizeUrlTest{
			"https://hostname/foo/",
			"",
			"https://hostname/foo/",
			false,
		},
		ContextualizeUrlTest{
			"/foo/",
			"",
			"/foo/",
			false,
		},
		ContextualizeUrlTest{
			"",
			"",
			"",
			false,
		},
		ContextualizeUrlTest{
			"%",
			"/bar/",
			"parse %: invalid URL escape \"%\"",
			true,
		},
		ContextualizeUrlTest{
			"/foo/",
			"%",
			"parse %: invalid URL escape \"%\"",
			true,
		},
	}
	for _, test := range tests {
		var req ScraperRequest
		req.Url = test.SR

		got, err := req.ContextualizeUrl(test.Input)
		if test.Error && err == nil {
			t.Fatal("Should have failed: %#v", test.Input)
		} else if !test.Error && err != nil {
			t.Fatal(err)
		}

		if test.Error {
			got = err.Error()
		}
		compare(t, test.Output, got)
	}
}

func TestSR_QueueAnother(t *testing.T) {
	trq := NewTestRQ()
	req := trq.CreateRequest("http://hostname/pasta")
	req.QueueAnother("/sauce")

	if len(trq.Queued) != 1 {
		t.Fatalf("Should have queued one item - %d queued", len(trq.Queued))
	}

	queued := trq.Queued[0]
	compare(t, "http://hostname/sauce", queued.Url)
	compare(t, "http://hostname/pasta", queued.Context.Referer)
	compare(t, "", trq.Debug.String())
}

func TestSR_QueueAnother_BadUrl(t *testing.T) {
	trq := NewTestRQ()
	req := trq.CreateRequest("http://hostname/pasta")
	req.QueueAnother("%")

	if len(trq.Queued) != 0 {
		t.Fatalf("Should not have queued any items - %d queued", len(trq.Queued))
	}
	compare(t,
		"http://hostname/pasta: parse %: invalid URL escape \"%\"\n",
		trq.Debug.String(),
	)
}
