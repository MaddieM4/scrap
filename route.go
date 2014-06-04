package scrap

type Route struct {
	Selector StringTest
	Action   RouteAction
}

// Run for each parsed page
type RouteAction func(req ScraperRequest, root Node)
