package scrap

import (
	"errors"
	"io"
	"log"
	"sync"
)

type ScraperConfig struct {
	Retriever Retriever
	Remarks   io.Writer
	Debug     io.Writer
}

func (sc ScraperConfig) Validate() error {
	if sc.Retriever == nil {
		return errors.New("ScraperConfig not valid if Retriever == nil")
	}
	if sc.Remarks == nil {
		return errors.New("ScraperConfig not valid if Remarks == nil")
	}
	if sc.Debug == nil {
		return errors.New("ScraperConfig not valid if Debug == nil")
	}
	return nil
}

type Scraper struct {
	config ScraperConfig
	Routes *RouteSet
	seen   map[string]bool
	wg     *sync.WaitGroup
}

// May return an error if config validation fails.
func NewScraper(config ScraperConfig) (Scraper, error) {
	err := config.Validate()
	if err != nil {
		return Scraper{}, err
	}
	return Scraper{
		config: config,
		Routes: NewRouteSet(),
		seen:   make(map[string]bool),
		wg:     new(sync.WaitGroup),
	}, nil
}

// Creates a new ScraperRequest with its properties all initialized.
func (s *Scraper) CreateRequest(url string) ScraperRequest {
	return ScraperRequest{
		Url:          url,
		RequestQueue: s,
		Remarks:      log.New(s.config.Remarks, url+": ", 0),
		Debug:        log.New(s.config.Debug, url+": ", 0),
	}
}

// Scrape a URL based on the given ScraperRequest.
func (s *Scraper) DoRequest(req ScraperRequest) {
	if s.seen[req.Url] {
		return
	}
	s.seen[req.Url] = true

	route, ok := s.Routes.MatchUrl(req.Url)
	if ok {
		req.Debug.Println("Found a route")
		route.Run(req, s.config.Retriever, s.wg)
	} else {
		req.Debug.Println("No route found")
	}
}

// Wait for all outstanding queued items to finish. You almost always
// want to do this, so that your main function doesn't end (thus ending
// the entire process) while you still have all your goroutines out in
// limbo.
func (s *Scraper) Wait() {
	s.wg.Wait()
}

// Convenience function to create and queue a new request.
func (s *Scraper) Scrape(url string) {
	s.DoRequest(s.CreateRequest(url))
}
