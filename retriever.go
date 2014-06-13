package scrap

import (
	"net/http"
)

// In a production environment, you will always use HttpRetriever.
type Retriever func(ScraperRequest) (*http.Response, error)

// Retrieves pages via HTTP or HTTPS, depending on URL.
func HttpRetriever(req ScraperRequest) (*http.Response, error) {
	return http.Get(req.Url)
}
