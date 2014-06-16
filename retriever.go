package scrap

import (
	"net/http"
)

// In a production environment, you will always use HttpRetriever.
type Retriever func(ScraperRequest) (*http.Response, error)

// Retrieves pages via HTTP or HTTPS, depending on URL.
func HttpRetriever(req ScraperRequest) (*http.Response, error) {
	var client http.Client
	request, err := http.NewRequest("GET", req.Url, nil)
	if err != nil {
		return nil, err
	}

	if req.Auth != nil {
		request.SetBasicAuth(req.Auth.Username, req.Auth.Password)
	}
	return client.Do(request)
}
