package scraper

import (
	"code.google.com/p/go.net/html"
	"fmt"
	"net/http"
	"sync"
)

type Scraper struct {
	StartingUrl string
	Debug       bool

	routes  []Route
	remarks chan string
	seen    map[string]bool
	wg      *sync.WaitGroup
}

func New(max_queued int) Scraper {
	return Scraper{
		routes:  make([]Route, 0),
		remarks: make(chan string, max_queued),
		seen:    make(map[string]bool),
		wg:      new(sync.WaitGroup),
	}
}

func (s *Scraper) Run() {
	s.NewRequest(s.StartingUrl)
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

func (s *Scraper) NewRequest(url string) {
	if s.seen[url] {
		return
	}
	s.seen[url] = true
	s.wg.Add(1)

	sr := ScraperRequest{
		Url:     url,
		Action:  nil,
		scraper: s,
	}

	for _, r := range s.routes {
		if r.Selector.Test(url) {
			sr.Debug("Found a route")
			sr.Action = r.Action
			break
		}
	}

	if sr.Action != nil {
		go func() {
			defer s.wg.Done()

			resp, err := http.Get(url)
			if err != nil {
				sr.Remark(err.Error())
				return
			}
			doc, err := html.Parse(resp.Body)
			if err != nil {
				sr.Remark(err.Error())
				return
			}

			sr.Action.Run(sr, doc)
		}()
	} else {
		sr.Debug("No route found")
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
