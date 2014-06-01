package scrap

// Scraper implements this. But we leave it as an interface so that
// we can mock it in tests.
type SRQueuer interface {
	CreateRequest(string) ScraperRequest
	DoRequest(ScraperRequest)
}
