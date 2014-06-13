package scrap

import (
	"bytes"
	"net/http"
	"testing"
)

func dummyRetriever(ScraperRequest) (*http.Response, error) {
	return nil, nil
}

type sc_valid_test struct {
	Config   ScraperConfig
	ErrorMsg string
}

func (svt sc_valid_test) Run(t *testing.T) {
	err := svt.Config.Validate()
	if svt.ErrorMsg == "" {
		// Should not result in an error
		if err != nil {
			t.Fatal(err)
		}
	} else {
		// Should result in specific error
		if err == nil {
			t.Fatalf("Should have failed (%s), but didn't", svt.ErrorMsg)
		}
		compare(t, svt.ErrorMsg, err.Error())
	}
}

func TestScraperConfig_Validate(t *testing.T) {
	var remarks, debug bytes.Buffer
	tests := []sc_valid_test{
		sc_valid_test{
			ScraperConfig{
				Retriever: dummyRetriever,
				Bucket:    NewCountBucket(1),
				Remarks:   &remarks,
				Debug:     &debug,
			},
			"",
		},
		sc_valid_test{
			ScraperConfig{
				Retriever: nil,
				Bucket:    NewCountBucket(1),
				Remarks:   &remarks,
				Debug:     &debug,
			},
			"ScraperConfig not valid if Retriever == nil",
		},
		sc_valid_test{
			ScraperConfig{
				Retriever: dummyRetriever,
				Bucket:    nil,
				Remarks:   &remarks,
				Debug:     &debug,
			},
			"ScraperConfig not valid if Bucket == nil",
		},
		sc_valid_test{
			ScraperConfig{
				Retriever: dummyRetriever,
				Bucket:    NewCountBucket(1),
				Remarks:   nil,
				Debug:     &debug,
			},
			"ScraperConfig not valid if Remarks == nil",
		},
		sc_valid_test{
			ScraperConfig{
				Retriever: dummyRetriever,
				Bucket:    NewCountBucket(1),
				Remarks:   &remarks,
				Debug:     nil,
			},
			"ScraperConfig not valid if Debug == nil",
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestScraper_BadConfig(t *testing.T) {
	config := ScraperConfig{}
	_, err := NewScraper(config)
	if err == nil {
		t.Fatal("Should have failed, didn't")
	}
}

func TestScraper_CreateRequest(t *testing.T) {
	var remarks, debug bytes.Buffer
	config := ScraperConfig{
		Retriever: dummyRetriever,
		Bucket:    NewCountBucket(1),
		Remarks:   &remarks,
		Debug:     &debug,
	}
	s, err := NewScraper(config)
	if err != nil {
		t.Fatal(err)
	}
	req := s.CreateRequest("/")

	if req.RequestQueue.(*Scraper) != &s {
		t.Fatal("Should have set req.RequestQueue")
	}
	if req.Url != "/" {
		t.Fatal("Should have set req.Url")
	}

	req.Remarks.Println("This is a remark")
	req.Debug.Println("This is a debug note")

	compare(t, "/: This is a remark\n", remarks.String())
	compare(t, "/: This is a debug note\n", debug.String())
}

func TestScraper_DoRequest_NoRoute(t *testing.T) {
	var remarks, debug bytes.Buffer
	config := ScraperConfig{
		Retriever: testHtmlRetriever,
		Bucket:    NewCountBucket(1),
		Remarks:   &remarks,
		Debug:     &debug,
	}
	s, err := NewScraper(config)
	if err != nil {
		t.Fatal(err)
	}

	req := s.CreateRequest("/")
	s.DoRequest(req)

	compare(t, "", remarks.String())
	compare(t, "/: No route found\n", debug.String())
}

func TestScraper_DoRequest_Seen(t *testing.T) {
	var remarks, debug bytes.Buffer
	config := ScraperConfig{
		Retriever: testHtmlRetriever,
		Bucket:    NewCountBucket(1),
		Remarks:   &remarks,
		Debug:     &debug,
	}
	s, err := NewScraper(config)
	if err != nil {
		t.Fatal(err)
	}
	s.Routes.AppendPrefix("/", func(req ScraperRequest, resp ServerResponse) {
		root, _ := resp.Parse()
		req.Remarks.Printf("%d <a> elements\n", len(root.Find("a")))
	})

	req := s.CreateRequest("/")
	s.DoRequest(req)
	s.DoRequest(req) // Same request again
	s.Wait()

	// Only one instance of this element being run
	compare(t, "/: 3 <a> elements\n", remarks.String())
	compare(t, "/: Found a route\n", debug.String())

	// Different request
	req = s.CreateRequest("/different")
	s.DoRequest(req)
	s.Wait()

	// Each unique request was run exactly once
	compare(t,
		"/: 3 <a> elements\n/different: 3 <a> elements\n",
		remarks.String(),
	)
	compare(t,
		"/: Found a route\n/different: Found a route\n",
		debug.String(),
	)
}

// The "Seen" test covers basic usage. Let's make sure we're actually
// using the bucket, though, instead of an internal counter.
func TestScraper_UsesBucket(t *testing.T) {
	var remarks, debug bytes.Buffer
	max_hits := 5
	config := ScraperConfig{
		Retriever: testHtmlRetriever,
		Bucket:    NewCountBucket(max_hits),
		Remarks:   &remarks,
		Debug:     &debug,
	}
	s, err := NewScraper(config)
	if err != nil {
		t.Fatal(err)
	}
	s.Routes.AppendExact("/", func(req ScraperRequest, resp ServerResponse) {
		req.QueueAnother("/") // Keeps queueing itself
	})

	var expected_debug string
	each_line := "/: Found a route\n"
	for i := 0; i < max_hits; i++ {
		expected_debug = expected_debug + each_line
	}

	s.Scrape("/")
	s.Wait()

	compare(t, expected_debug, debug.String())
}

func TestScraper_Scrape(t *testing.T) {
	var remarks, debug bytes.Buffer
	config := ScraperConfig{
		Retriever: testHtmlRetriever,
		Bucket:    NewCountBucket(1),
		Remarks:   &remarks,
		Debug:     &debug,
	}
	s, err := NewScraper(config)
	if err != nil {
		t.Fatal(err)
	}
	s.Routes.AppendPrefix("/", func(req ScraperRequest, resp ServerResponse) {
		root, _ := resp.Parse()
		req.Remarks.Printf("%d <a> elements\n", len(root.Find("a")))
		req.QueueAnother("/second")
	})

	s.Scrape("/")
	s.Wait()

	compare(t,
		"/: 3 <a> elements\n/second: 3 <a> elements\n",
		remarks.String(),
	)
}
