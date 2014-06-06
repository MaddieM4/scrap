package scrap

import "sync"

type Route struct {
	Selector StringTest
	Action   RouteAction
}

func (r Route) Matches(url string) bool {
	return r.Selector(url)
}

// Runs r.Action in a goroutine, subscribing it on the WaitGroup
func (r Route) Run(req ScraperRequest, ret Retriever, wg *sync.WaitGroup) {
	n, err := ret(req)
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

// Run for each parsed page
type RouteAction func(req ScraperRequest, root Node)

type RouteSet struct {
	Routes []Route
}

func NewRouteSet() *RouteSet {
	return &RouteSet{
		make([]Route, 0),
	}
}

func (rs *RouteSet) Append(r Route) {
	rs.Routes = append(rs.Routes, r)
}

func (rs *RouteSet) AppendExact(url string, action RouteAction) {
	rs.Append(Route{
		Selector: StringTestExact(url),
		Action:   action,
	})
}

func (rs *RouteSet) AppendPrefix(prefix string, action RouteAction) {
	rs.Append(Route{
		Selector: StringTestPrefix(prefix),
		Action:   action,
	})
}

func (rs *RouteSet) MatchUrl(url string) (Route, bool) {
	for _, r := range rs.Routes {
		if r.Matches(url) {
			return r, true
		}
	}
	return Route{}, false
}
