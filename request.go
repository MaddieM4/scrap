package scrap

import (
	"log"
	"net/url"
)

// Represents a single request.
//
// This is handed to all RouteActions.
type ScraperRequest struct {
	Url          string
	RequestQueue SRQueuer
	Remarks      *log.Logger
	Debug        *log.Logger
}

// Queue another URL for scraping. Duplicate queued items are ignored.
func (sr ScraperRequest) QueueAnother(queue_url string) {
	old_url, err := url.Parse(sr.Url)
	if err != nil {
		sr.Debug.Println(err.Error())
		return
	}
	new_url, err := url.Parse(queue_url)
	if err != nil {
		sr.Debug.Println(err.Error())
		return
	}
	new_url.Fragment = "" // Treat all #whatever urls the same
	final_url := old_url.ResolveReference(new_url)

	new_req := sr.RequestQueue.CreateRequest(final_url.String())
	sr.RequestQueue.DoRequest(new_req)
}
