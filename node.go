package scrap

import (
	"code.google.com/p/go-html-transform/css/selector"
	"code.google.com/p/go.net/html"
	//"net/url"
)

// Wrapper for html.node with selector capabilities.
type Node struct {
	*html.Node
	req *ScraperRequest
}

// Find a set of descendent nodes based on CSS3 selector.
func (n *Node) Find(sel string) NodeSet {
	chain, err := selector.Selector(sel)
	if err != nil {
		n.req.Debug.Println(err.Error())
		return nil
	}
	return WrapNodes(chain.Find(n.Node), n.req)
}

// Queue this node's 'href' attr value as a URL to scrape.
func (n Node) Queue() {
}
