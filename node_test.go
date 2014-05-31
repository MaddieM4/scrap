package scrap

import (
	"code.google.com/p/go.net/html"
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

func setup(t *testing.T) Node {
	scraper := NewScraper()
	request := ScraperRequest{
		Url:     "foo",
		scraper: &scraper,
	}
	raw_node, err := html.Parse(strings.NewReader(sample_html))
	if err != nil {
		t.Fatal(err)
	}

	return Node{raw_node, &request}
}

func TestNode_Find_NoSuchElement(t *testing.T) {
	n := setup(t)
	found := n.Find("blah")
	if len(found) != 0 {
		t.Fatalf("Slice should have 0 elements, has %d!", len(found))
	}
}

func TestNode_Find_OneElement(t *testing.T) {
	n := setup(t)
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
	n := setup(t)
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
	n := setup(t)
	scraper := n.req.scraper
	scraper.Debug = false
	selector := "*&{"

	found := n.Find(selector)
	if len(found) != 0 {
		t.Fatalf("Slice should have 0 elements, has %d!", len(found))
	}
	if len(scraper.remarks) != 0 {
		t.Fatalf("Should have 0 remarks, has %d", len(scraper.remarks))
	}

	scraper.Debug = true
	found = n.Find(selector)
	if len(found) != 0 {
		t.Fatalf("Slice should have 0 elements, has %d!", len(found))
	}
	if len(scraper.remarks) != 1 {
		t.Fatalf("Should have 1 remark, has %d", len(scraper.remarks))
	}

	err_remark := <-scraper.remarks
	expected_remark := "foo: End of Selector"
	if err_remark != expected_remark {
		t.Fatalf(
			"Expected error: \"%s\"\n\nGot error: \"%s\"",
			expected_remark,
			err_remark,
		)
	}
}

func TestNode_Attr_NoElements(t *testing.T) {
	n := setup(t)
	attrs := n.Find("foo").Attr("bar")
	if len(attrs) != 0 {
		t.Fatalf("Should have 0 results, got %d", len(attrs))
	}
}

func TestNode_Attr_ElementsButNoAttr(t *testing.T) {
	n := setup(t)
	attrs := n.Find("p").Attr("bar")
	if len(attrs) != 0 {
		t.Fatalf("Should have 0 results, got %d", len(attrs))
	}
}

// Should only ever return one result for each element
func TestNode_Attr_OneElementMultipleAttrs(t *testing.T) {
	n := setup(t)
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
	n := setup(t)
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
	n := setup(t)
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
