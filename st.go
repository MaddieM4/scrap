package scrap

import "strings"

// Used to determine if a URL fits a Route.
type StringTest func(string) bool

// Create a callback that only returns true for exact matches.
func StringTestExact(pattern string) StringTest {
	return func(url string) bool {
		return url == pattern
	}
}

// Create a callback that only returns true if the URL starts with
// the given prefix.
func StringTestPrefix(pattern string) StringTest {
	return func(url string) bool {
		return strings.HasPrefix(url, pattern)
	}
}
