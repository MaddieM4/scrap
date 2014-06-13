package scrap

import (
	"code.google.com/p/go.net/html"
	"errors"
	"net/http"
	"time"
)

// Represents a response from the server, containing the original
// http.Response object, extra contextual/statistic data, and various
// convenience functions.
type ServerResponse struct {
	Response *http.Response
	Request  ScraperRequest
	Stats    ResponseStats
}

// Some statistics on the completed request, for example, the time
// required to retrieve the file from the network.
type ResponseStats struct {
	Start    time.Time
	Duration time.Duration
}

// Get a ServerResponse based on a ScraperRequest and a Retriever.
//
// Handles things like ResponseStats, so that Retrievers can just focus
// on handing off an http.Response.
func GetResponse(req ScraperRequest, ret Retriever) (ServerResponse, error) {
	response := ServerResponse{}
	response.Request = req
	response.Stats.Start = time.Now()

	http_resp, err := ret(req)
	if err != nil {
		return response, err
	}

	response.Response = http_resp
	response.Stats.Duration = time.Since(response.Stats.Start)
	return response, nil
}

// Get the HTML node contents of the scraped data.
//
// This consumes resp.Response.Body.
func (resp ServerResponse) Parse() (Node, error) {
	if resp.Response == nil || resp.Response.Body == nil {
		return Node{}, errors.New("Could not parse, nil pointer")
	}
	raw_node, err := html.Parse(resp.Response.Body)
	if err != nil {
		return Node{}, err
	} else {
		return Node{raw_node, &resp.Request}, nil
	}
}
