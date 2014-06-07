package scrap

import (
	"sync"
	"time"
)

// A binding between a URL-matching function, and an action to perform
// on pages where the URL matches.
//
// You always want to set up your Routes before starting to scrape, or
// else none of the scraped pages will match.
type Route struct {
	Selector StringTest
	Action   RouteAction
}

// Does the given url match this Route? Used by the Scraper to select
// the first matching Route.
func (r Route) Matches(url string) bool {
	return r.Selector(url)
}

// Runs r.Action in a goroutine, subscribing it on the WaitGroup
func (r Route) Run(req ScraperRequest, ret Retriever, wg *sync.WaitGroup) {
	start_time := time.Now()
	n, err := ret(req)
	req.Stats.Duration = time.Since(start_time)

	if err != nil {
		req.Debug.Println(err.Error())
		return
	}

	// Add to wg in calling goroutine
	wg.Add(1)

	// Decrement wg in spawned goroutine, after performing action
	go func() {
		defer wg.Done()
		r.Action(req, n)
	}()
}

// Callback that's run for each parsed page.
type RouteAction func(req ScraperRequest, root Node)

// A slice of Routes. Order is important!
type RouteSet struct {
	Routes []Route
}

func NewRouteSet() *RouteSet {
	return &RouteSet{
		make([]Route, 0),
	}
}

// Add a new Route at the end of the set.
func (rs *RouteSet) Append(r Route) {
	rs.Routes = append(rs.Routes, r)
}

// Shorthand to add a new Route at the end of the set, where exact
// URL matching is used (see StringTestExact).
func (rs *RouteSet) AppendExact(url string, action RouteAction) {
	rs.Append(Route{
		Selector: StringTestExact(url),
		Action:   action,
	})
}

// Shorthand to add a new Route at the end of the set, where prefix
// URL matching is used (see StringTestPrefix).
func (rs *RouteSet) AppendPrefix(prefix string, action RouteAction) {
	rs.Append(Route{
		Selector: StringTestPrefix(prefix),
		Action:   action,
	})
}

// Return the first Route where the URL matches according to the
// Route's matching function. Order is important!
func (rs *RouteSet) MatchUrl(url string) (Route, bool) {
	for _, r := range rs.Routes {
		if r.Matches(url) {
			return r, true
		}
	}
	return Route{}, false
}
