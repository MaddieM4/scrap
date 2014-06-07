package scrap

import (
	"bytes"
	"log"
	"reflect"
	"testing"
)

var sample_html = `
<html>
    <head>
        <title>Foo</title>
        <link rel="this" href="/blah/" rel="that"/>
    </head>
    <body>
        <p>First paragraph.</p>
        <p>Second paragraph with <em>emphasized text</em>.</p>
        <p>Third paragraph.</p>
        <hr/>
        <a href="/first">First link</a>
        <a href="/second">Second link</a>
        <a href="/third">Third link</a>
    </body>
</html>`

type TestRQ struct {
	Remarks *bytes.Buffer
	Debug   *bytes.Buffer
	Queued  []ScraperRequest
}

func NewTestRQ() *TestRQ {
	trq := new(TestRQ)
	trq.Remarks = new(bytes.Buffer)
	trq.Debug = new(bytes.Buffer)
	trq.Queued = make([]ScraperRequest, 0)
	return trq
}
func (trq *TestRQ) CreateRequest(url string) ScraperRequest {
	return ScraperRequest{
		Url:          url,
		RequestQueue: trq,
		Remarks:      log.New(trq.Remarks, url+": ", 0),
		Debug:        log.New(trq.Debug, url+": ", 0),
		Stats:        new(RequestStats),
	}
}
func (trq *TestRQ) DoRequest(req ScraperRequest) {
	trq.Queued = append(trq.Queued, req)
}

func testHtmlRetriever(req ScraperRequest) (Node, error) {
	data := new(bytes.Buffer)
	data.WriteString(sample_html)
	return parseReader(req, data)
}

func compare(t *testing.T, expected, got interface{}) {
	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("Expected %#v, got %#v", expected, got)
	}
}
