package scrap

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"log"
	"reflect"
	"strings"
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
	}
}
func (trq *TestRQ) DoRequest(req ScraperRequest) {
	trq.Queued = append(trq.Queued, req)
}

func setupNode(t *testing.T) Node {
	trq := NewTestRQ()
	request := trq.CreateRequest("foo")

	raw_node, err := html.Parse(strings.NewReader(sample_html))
	if err != nil {
		t.Fatal(err)
	}

	return Node{raw_node, &request}
}

func compare(t *testing.T, expected, got interface{}) {
	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("Expected %#v, got %#v", expected, got)
	}
}

func TestNode_Find_NoSuchElement(t *testing.T) {
	n := setupNode(t)
	found := n.Find("blah")
	if len(found) != 0 {
		t.Fatalf("Slice should have 0 elements, has %d!", len(found))
	}
}

func TestNode_Find_OneElement(t *testing.T) {
	n := setupNode(t)
	tagname := "body"
	found := n.Find(tagname)
	if len(found) != 1 {
		t.Fatalf("Slice should have 1 element, has %d!", len(found))
	}
	found_tagname := found[0].Node.Data
	if found_tagname != tagname {
		t.Fatalf("Expected %s, got %s", tagname, found_tagname)
	}
	if found[0].req != n.req {
		t.Fatalf(
			"Found node has different req! %r vs %r",
			found[0].req,
			n.req,
		)
	}
}

func TestNode_Find_MultipleElements(t *testing.T) {
	n := setupNode(t)
	tagname := "p"
	found := n.Find(tagname)
	if len(found) != 3 {
		t.Fatalf("Slice should have 1 element, has %d!", len(found))
	}
	for _, f := range found {
		found_tagname := f.Node.Data
		if found_tagname != tagname {
			t.Fatalf("Expected %s, got %s", tagname, found_tagname)
		}
		if f.req != n.req {
			t.Fatalf("Found node has different req! %r vs %r", f.req, n.req)
		}
	}
}

func TestNode_Find_BadSelector(t *testing.T) {
	n := setupNode(t)
	rq := n.req.RequestQueue
	selector := "*&{"

	found := n.Find(selector)
	if len(found) != 0 {
		t.Fatalf("Slice should have 0 elements, has %d!", len(found))
	}

	expected_remarks := ""
	expected_debug := "foo: End of Selector\n"
	got_remarks := rq.(*TestRQ).Remarks.String()
	got_debug := rq.(*TestRQ).Debug.String()

	compare(t, expected_remarks, got_remarks)
	compare(t, expected_debug, got_debug)
}

func TestNode_Attr_NoElements(t *testing.T) {
	n := setupNode(t)
	attrs := n.Find("foo").Attr("bar")
	if len(attrs) != 0 {
		t.Fatalf("Should have 0 results, got %d", len(attrs))
	}
}

func TestNode_Attr_ElementsButNoAttr(t *testing.T) {
	n := setupNode(t)
	attrs := n.Find("p").Attr("bar")
	if len(attrs) != 0 {
		t.Fatalf("Should have 0 results, got %d", len(attrs))
	}
}

// Should only ever return one result for each element
func TestNode_Attr_OneElementMultipleAttrs(t *testing.T) {
	n := setupNode(t)
	attrs := n.Find("link").Attr("rel")
	if len(attrs) != 1 {
		t.Fatalf("Should have 1 results, got %d", len(attrs))
	}

	// Order should be deterministic
	expected := "that"
	got := attrs[0]
	if got != expected {
		t.Fatalf("Expected %s, got %s", expected, got)
	}
}

func TestNode_Attr_OneElementMultipleAttrsFiltered(t *testing.T) {
	n := setupNode(t)
	attrs := n.Find("link").Attr("href")
	if len(attrs) != 1 {
		t.Fatalf("Should have 1 results, got %d", len(attrs))
	}

	expected := "/blah/"
	got := attrs[0]
	if got != expected {
		t.Fatalf("Expected %s, got %s", expected, got)
	}
}

func TestNode_Attr_MultipleElements(t *testing.T) {
	n := setupNode(t)
	attrs := n.Find("a").Attr("href")
	expected := []string{
		"/first",
		"/second",
		"/third",
	}
	if !reflect.DeepEqual(attrs, expected) {
		t.Fatalf("Expected %s, got %s", expected, attrs)
	}
}

func TestNode_Queue_NoHref(t *testing.T) {
	n := setupNode(t)
	n.Find("body").Queue()
}
