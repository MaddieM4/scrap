package scrap

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Scraper struct {
	StartingUrl string
	Retriever   Retriever
	Debug       bool

	routes  []Route
	remarks chan string
	seen    map[string]bool
	wg      *sync.WaitGroup
}

func NewScraper() Scraper {
	return Scraper{
		routes:  make([]Route, 0),
		remarks: make(chan string, 100),
		seen:    make(map[string]bool),
		wg:      new(sync.WaitGroup),
	}
}

func (s *Scraper) Run() {
	s.DoRequest(s.CreateRequest(s.StartingUrl))
	done_with_remarks := make(chan int)

	go func() {
		for {
			remark, ok := <-s.remarks
			if ok {
				fmt.Println(remark)
			} else {
				close(done_with_remarks)
				return
			}
		}
	}()
	s.wg.Wait()
	close(s.remarks)
	<-done_with_remarks // Wait for all remarks to print
}

func (s *Scraper) CreateRequest(url string) ScraperRequest {
	var debug_writer io.Writer
	if s.Debug {
		debug_writer = os.Stderr
	} else {
		debug_writer = ioutil.Discard
	}
	return ScraperRequest{
		Url:          url,
		RequestQueue: s,
		Remarks:      log.New(os.Stdout, url+": ", 0),
		Debug:        log.New(debug_writer, url+": ", 0),
	}
}

func (s *Scraper) DoRequest(req ScraperRequest) {
	if s.seen[req.Url] {
		return
	}
	s.seen[req.Url] = true
	s.wg.Add(1)
	var route *Route

	for r := range s.routes {
		if s.routes[r].Selector(req.Url) {
			req.Debug.Println("Found a route")
			route = &s.routes[r]
			break
		}
	}

	if route != nil {
		go func() {
			defer s.wg.Done()

			doc, err := s.Retriever(req)
			if err != nil {
				req.Debug.Println(err.Error())
				return
			}

			route.Action.Run(req, doc)
		}()
	} else {
		req.Debug.Println("No route found")
		s.wg.Done()
	}
}

// Add a route to a scraper
func (s *Scraper) Route(r Route) {
	s.routes = append(s.routes, r)
}

// Matches URLs exactly
func (s *Scraper) RouteExact(url string, action RouteAction) {
	s.Route(Route{
		Selector: StringTestExact(url),
		Action:   action,
	})
}

// Matches URLs with a specific prefix
func (s *Scraper) RoutePrefix(prefix string, action RouteAction) {
	s.Route(Route{
		Selector: StringTestPrefix(prefix),
		Action:   action,
	})
}
