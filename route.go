package scraper

import "code.google.com/p/go.net/html"

type Route struct {
	Selector StringTest
	Action   RouteAction
}

// Run for each parsed page
type RouteAction interface {
	Run(req ScraperRequest, root *html.Node)
}
