package scrap

import (
	"errors"
	"io"
	"log"
	"sync"
)

type ScraperConfig struct {
	StartingUrl string
	Retriever   Retriever
	Remarks     io.Writer
	Debug       io.Writer
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

func (s *Scraper) Run() {
	s.DoRequest(s.CreateRequest(s.config.StartingUrl))
	s.wg.Wait()
}

func (s *Scraper) CreateRequest(url string) ScraperRequest {
	return ScraperRequest{
		Url:          url,
		RequestQueue: s,
		Remarks:      log.New(s.config.Remarks, url+": ", 0),
		Debug:        log.New(s.config.Debug, url+": ", 0),
	}
}

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
