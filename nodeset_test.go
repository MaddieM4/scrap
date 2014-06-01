package scrap

import "testing"

func TestNodeSet_Attr_NoElements(t *testing.T) {
	n := setupNode(t)
	attrs := n.Find("foo").Attr("bar")
	if len(attrs) != 0 {
		t.Fatalf("Should have 0 results, got %d", len(attrs))
	}
}

func TestNodeSet_Attr_ElementsButNoAttr(t *testing.T) {
	n := setupNode(t)
	attrs := n.Find("p").Attr("bar")
	if len(attrs) != 0 {
		t.Fatalf("Should have 0 results, got %d", len(attrs))
	}
}

// Should only ever return one result for each element
func TestNodeSet_Attr_OneElementMultipleAttrs(t *testing.T) {
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

func TestNodeSet_Attr_OneElementMultipleAttrsFiltered(t *testing.T) {
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

func TestNodeSet_Attr_MultipleElements(t *testing.T) {
	n := setupNode(t)
	attrs := n.Find("a").Attr("href")
	expected := []string{
		"/first",
		"/second",
		"/third",
	}
	compare(t, expected, attrs)
}
