package scrap

import "net/url"

// Represents a single request.
//
// This is handed to all RouteActions.
type ScraperRequest struct {
	Url    string
	Action RouteAction

	scraper *Scraper
}

// Output a message. Will show up in output with the URL prepended.
func (sr ScraperRequest) Remark(r string) {
	sr.scraper.remarks <- sr.Url + ": " + r
}

// Like Remark, but only outputs the remark if scraper.Debug == true.
func (sr ScraperRequest) Debug(r string) {
	if sr.scraper.Debug {
		sr.Remark(r)
	}
}

// Queue another URL for scraping. Duplicate queued items are ignored.
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
