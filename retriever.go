package scrap

import (
	"code.google.com/p/go.net/html"
	"errors"
	"io"
	"net/http"
)

type Retriever func(ScraperRequest) (Node, error)

func parseReader(req ScraperRequest, r io.Reader) (Node, error) {
	raw_node, err := html.Parse(r)
	if err != nil {
		return Node{}, err
	} else {
		return Node{raw_node, &req}, nil
	}
}

// Retrieves pages via HTTP or HTTPS, depending on URL.
func HttpRetriever(req ScraperRequest) (Node, error) {
	resp, err := http.Get(req.Url)
	if err != nil {
		return Node{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Node{}, errors.New(resp.Status)
	}

	return parseReader(req, resp.Body)
}
