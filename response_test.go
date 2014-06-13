package scrap

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestServerResponse_GetResponse(t *testing.T) {
	url := "http://no.such.host/"
	trq := NewTestRQ()
	req := trq.CreateRequest(url)
	ret := testHtmlRetriever

	resp, err := GetResponse(req, ret)
	if err != nil {
		t.Fatal(err)
	}

	comparem(t, req, resp.Request, "resp.Request not set to req")
	var empty_time time.Time
	var empty_duration time.Duration
	if resp.Stats.Start == empty_time {
		t.Fatal("resp.Stats.Start not set")
	}
	if resp.Stats.Duration == empty_duration {
		t.Fatal("resp.Stats.Duration not set")
	}
}

func TestServerResponse_Parse(t *testing.T) {
	url := "http://no.such.host/"
	trq := NewTestRQ()
	req := trq.CreateRequest(url)
	ret := testHtmlRetriever

	resp, err := GetResponse(req, ret)
	if err != nil {
		t.Fatal(err)
	}
	node, err := resp.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(node.Find("a")) != 3 {
		t.Fatal("Parsed HTML is wonky or blank")
	}
}

type BrokenReader struct{}

func (br BrokenReader) Read([]byte) (int, error) {
	return 0, errors.New("BrokenReader says hello")
}
func (br BrokenReader) Close() error {
	return nil
}
func BrokenRetriever(req ScraperRequest) (*http.Response, error) {
	var resp http.Response
	resp.Body = BrokenReader{}
	return &resp, nil
}

func TestServerResponse_Parse_BrokenReader(t *testing.T) {
	url := "http://no.such.host/"
	trq := NewTestRQ()
	req := trq.CreateRequest(url)
	ret := BrokenRetriever

	resp, err := GetResponse(req, ret)
	if err != nil {
		t.Fatal(err)
	}
	_, err = resp.Parse()
	if err == nil {
		t.Fatal("Should have failed - broken reader")
	}
}

func NilRetriever(req ScraperRequest) (*http.Response, error) {
	return nil, nil
}
func TestServerResponse_Parse_NilResponse(t *testing.T) {
	url := "http://no.such.host/"
	trq := NewTestRQ()
	req := trq.CreateRequest(url)
	ret := NilRetriever

	resp, err := GetResponse(req, ret)
	if err != nil {
		t.Fatal(err)
	}
	_, err = resp.Parse()
	if err == nil {
		t.Fatal("Should have failed - nil response")
	}
}

func ErrorRetriever(req ScraperRequest) (*http.Response, error) {
	return nil, errors.New("Awww fudge")
}
func TestServerResponse_Parse_ErrorRetriever(t *testing.T) {
	url := "http://no.such.host/"
	trq := NewTestRQ()
	req := trq.CreateRequest(url)
	ret := ErrorRetriever

	_, err := GetResponse(req, ret)
	if err == nil {
		t.Fatal("Should have failed - ret returns error")
	}
}
