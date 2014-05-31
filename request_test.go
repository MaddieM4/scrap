package scrap

import (
	"testing"
)

func TestScraperRequest_Remark(t *testing.T) {
	scraper := NewScraper()
	request := ScraperRequest{
		Url:     "foo",
		scraper: &scraper,
	}
	if len(scraper.remarks) != 0 {
		t.Fatal("Should not have remarks yet")
	}
	request.Remark("bar")
	if len(scraper.remarks) != 1 {
		t.Fatal("Should have one remark now")
	}
	recvd := <-scraper.remarks
	expected := "foo: bar"
	if recvd != expected {
		t.Fatalf("Expected %v, got %v", expected, recvd)
	}
}

func TestScraperRequest_Debug(t *testing.T) {
	scraper := NewScraper()
	request := ScraperRequest{
		Url:     "foo",
		scraper: &scraper,
	}
	request.Debug("bar")
	if len(scraper.remarks) != 0 {
		t.Fatal("Should not have remarks yet")
	}
	scraper.Debug = true
	request.Debug("bar")
	if len(scraper.remarks) != 1 {
		t.Fatal("Should have one remark now")
	}
	recvd := <-scraper.remarks
	expected := "foo: bar"
	if recvd != expected {
		t.Fatalf("Expected %v, got %v", expected, recvd)
	}
}
