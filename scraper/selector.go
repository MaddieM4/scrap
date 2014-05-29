package scraper

import "strings"

// Used to determine if a URL fits a route
type StringTest interface {
	Test(string) bool
}

type StringTestExact string

func (s StringTestExact) Test(url string) bool {
	return url == string(s)
}

type StringTestPrefix string

func (s StringTestPrefix) Test(url string) bool {
	return strings.HasPrefix(url, string(s))
}
