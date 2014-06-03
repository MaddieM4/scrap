package scrap

import "strings"

// Used to determine if a URL fits a route
type StringTest func(string) bool

func StringTestExact(pattern string) StringTest {
	return func(url string) bool {
		return url == pattern
	}
}

func StringTestPrefix(pattern string) StringTest {
	return func(url string) bool {
		return strings.HasPrefix(url, pattern)
	}
}
