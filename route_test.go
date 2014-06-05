package scrap

import "testing"

func TestRoute_Matches(t *testing.T) {
	r := Route{
		Selector: StringTestExact("http://host/foo"),
		Action:   func(ScraperRequest, Node) {},
	}
	if r.Matches("ht") {
		t.Fatal("In this test, should only match exact")
	}
	if !r.Matches("http://host/foo") {
		t.Fatal("In this test, should match exact")
	}

	r.Selector = StringTestPrefix("ht")
	if r.Matches("/foo") {
		t.Fatal("In this test, should match prefix")
	}
	if !r.Matches("http://host/bar") {
		t.Fatal("In this test, should match prefix")
	}
}

func TestRouteSet_Append(t *testing.T) {
	rs := NewRouteSet()
	var feedback string
	r1 := Route{
		Selector: StringTestExact("http://host/foo"),
		Action:   func(ScraperRequest, Node) { feedback = "r1" },
	}
	r2 := Route{
		Selector: StringTestExact("http://host/bar"),
		Action:   func(ScraperRequest, Node) { feedback = "r2" },
	}
	rs.Append(r1)
	rs.Append(r2)

	compare(t, 2, len(rs.Routes))

	rs.Routes[0].Action(ScraperRequest{}, Node{})
	compare(t, "r1", feedback)

	rs.Routes[1].Action(ScraperRequest{}, Node{})
	compare(t, "r2", feedback)

}

func TestRouteSet_AppendExact(t *testing.T) {
	rs := NewRouteSet()
	rs.AppendExact("foo", func(ScraperRequest, Node) {})

	compare(t, 1, len(rs.Routes))
	_, ok := rs.MatchUrl("foobar")
	if ok {
		t.Fatal("Should not match foobar")
	}
	_, ok = rs.MatchUrl("foo")
	if !ok {
		t.Fatal("Should match foo")
	}
}

func TestRouteSet_AppendPrefix(t *testing.T) {
	rs := NewRouteSet()
	rs.AppendPrefix("foo", func(ScraperRequest, Node) {})

	compare(t, 1, len(rs.Routes))
	_, ok := rs.MatchUrl("foobar")
	if !ok {
		t.Fatal("Should match foobar")
	}
	_, ok = rs.MatchUrl("foo")
	if !ok {
		t.Fatal("Should match foo")
	}
	_, ok = rs.MatchUrl("kungfoo")
	if ok {
		t.Fatal("Should not match kungfoo")
	}
}

func TestRouteSet_MatchUrl(t *testing.T) {
	rs := NewRouteSet()
	var feedback string
	r1 := Route{
		Selector: StringTestExact("http://host/foo"),
		Action:   func(ScraperRequest, Node) { feedback = "r1" },
	}
	r2 := Route{
		Selector: StringTestExact("http://host/bar"),
		Action:   func(ScraperRequest, Node) { feedback = "r2" },
	}
	rs.Append(r1)
	rs.Append(r2)

	match, ok := rs.MatchUrl("http://host/foo")
	if !ok {
		t.Fatal("Should get a result for foo")
	}
	match.Action(ScraperRequest{}, Node{})
	compare(t, "r1", feedback)

	match, ok = rs.MatchUrl("http://host/bar")
	if !ok {
		t.Fatal("Should get a result for bar")
	}
	match.Action(ScraperRequest{}, Node{})
	compare(t, "r2", feedback)

	match, ok = rs.MatchUrl("http://host/baz")
	if ok {
		t.Fatal("Should not get a result for baz")
	}

}
