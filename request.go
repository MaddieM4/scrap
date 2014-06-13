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
	Context      RequestContext
}

// Used to convey information about the page and context where this
// request was queued. This can be useful for complex bucketing, or
// simply tracking down which pages are referring to 404 links.
type RequestContext struct {
	Referer string
}

// De-relativize a URL based on the existing request's URL.
//
// This is how "/foo/" URLs queued up for scraping are turned into
// more actionable "http://origin.host.name/foo/" URLs.
//
// Current behavior is to only give an absolute URL out, if the
// request's URL is absolute. The contextualization is only as
// good as the request URL's "absoluteness". You won't get an error
// if the result is ambiguous. THIS MAY CHANGE IN FUTURE RELEASES.
func (sr ScraperRequest) ContextualizeUrl(rel_url string) (string, error) {
	old_url, err := url.Parse(sr.Url)
	if err != nil {
		return "", err
	}

	new_url, err := url.Parse(rel_url)
	if err != nil {
		return "", err
	}

	new_url.Fragment = "" // Treat all #whatever urls the same
	final_url := old_url.ResolveReference(new_url)
	return final_url.String(), nil
}

// Queue another URL for scraping. Duplicate queued items are ignored.
func (sr ScraperRequest) QueueAnother(queue_url string) {
	abs_url, err := sr.ContextualizeUrl(queue_url)
	if err != nil {
		sr.Debug.Println(err)
		return
	}
	new_req := sr.RequestQueue.CreateRequest(abs_url)
	new_req.Context.Referer = sr.Url
	sr.RequestQueue.DoRequest(new_req)
}
