package main

import (
	"code.google.com/p/go.net/html"
	"github.com/campadrenalin/scrap/scraper"
)

type HomeAction struct{}

func (ha HomeAction) Run(req scraper.ScraperRequest, root *html.Node) {
	req.Remark("Homepage")
	req.QueueAnchors("a", root)
}

type DefaultAction struct{}

func (da DefaultAction) Run(req scraper.ScraperRequest, root *html.Node) {
	req.Remark("Page is in site")
	req.QueueAnchors("a", root)
	req.QueueAnother("http://localhost/bar/")
}

func main() {
	scrapey := scraper.New(50)
	scrapey.Debug = true
	scrapey.StartingUrl = "http://localhost/"
	scrapey.RouteExact("http://localhost/", HomeAction{})
	scrapey.RoutePrefix("http://localhost/", DefaultAction{})
	scrapey.Run()
}
