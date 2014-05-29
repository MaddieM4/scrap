package scrap

import (
	"code.google.com/p/go-html-transform/css/selector"
	"code.google.com/p/go.net/html"
	"net/url"
)

// Represents a single request.
//
// This is handed to all RouteActions.
type ScraperRequest struct {
	Url    string
	Action RouteAction

	scraper *Scraper
}

func (sr ScraperRequest) Remark(r string) {
	sr.scraper.remarks <- sr.Url + ": " + r
}

func (sr ScraperRequest) Debug(r string) {
	if sr.scraper.Debug {
		sr.Remark(r)
	}
}

func (sr ScraperRequest) QueueAnother(queue_url string) {
	old_url, err := url.Parse(sr.Url)
	if err != nil {
		sr.Remark(err.Error())
		return
	}
	new_url, err := url.Parse(queue_url)
	if err != nil {
		sr.Remark(err.Error())
		return
	}
	new_url.Fragment = "" // Treat all #whatever urls the same
	final_url := old_url.ResolveReference(new_url)

	sr.scraper.NewRequest(final_url.String())
}

func (sr ScraperRequest) Find(sel string, n *html.Node) []*html.Node {
	chain, err := selector.Selector(sel)
	if err != nil {
		sr.Remark(err.Error())
		return nil
	}
	return chain.Find(n)
}

func (sr ScraperRequest) QueueAnchors(sel string, n *html.Node) {
	anchors := sr.Find(sel, n)
	for _, element := range anchors {
		for _, attr := range element.Attr {
			if attr.Key == "href" {
				sr.QueueAnother(attr.Val)
			}
		}
	}
}
