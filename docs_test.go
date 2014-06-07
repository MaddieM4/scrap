package scrap

import (
	"fmt"
	"os"
)

func Example() {
	s, err := NewScraper(ScraperConfig{
		// Real code would use HttpRetriever as the Retriever.
		Retriever: testHtmlRetriever,
		Bucket:    NewCountBucket(1),
		Remarks:   os.Stdout,
		Debug:     os.Stdout,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	base_url := "http://www.example.com/"
	s.Routes.AppendPrefix(base_url, func(req ScraperRequest, root Node) {
		// Verify that there is only one <head> element
		num_heads := len(root.Find("head"))
		if num_heads != 1 {
			req.Remarks.Printf("%d heads, expected 1!\n", num_heads)
		}

		// Queue up any links on this page, for further scraping
		root.Find("a").Queue()
	})

	s.Scrape(base_url)
	s.Wait()

	// Output:
	// http://www.example.com/: Found a route
	// http://www.example.com/first: Found a route
	// http://www.example.com/second: Found a route
	// http://www.example.com/third: Found a route
}
