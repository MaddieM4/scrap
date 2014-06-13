package scrap

import (
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
