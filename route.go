package scrap

type Route struct {
	Selector StringTest
	Action   RouteAction
}

// Run for each parsed page
type RouteAction interface {
	Run(req ScraperRequest, root Node)
}
