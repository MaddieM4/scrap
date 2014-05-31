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

type NodeSet []Node

// Turn a slice of []*html.Node into a NodeSet.
func WrapNodes(raw_nodes []*html.Node, req *ScraperRequest) NodeSet {
	nodes := make(NodeSet, len(raw_nodes))
	for i := range raw_nodes {
		nodes[i] = Node{raw_nodes[i], req}
	}
	return nodes
}

// Find a set of descendent nodes based on CSS3 selector.
func (n *Node) Find(sel string) NodeSet {
	chain, err := selector.Selector(sel)
	if err != nil {
		n.req.Debug(err.Error())
		return nil
	}
	return WrapNodes(chain.Find(n.Node), n.req)
}

// Return a slice of attr values for each element in the Nodeset,
// where the attr name is equivalent to the one given.
func (ns NodeSet) Attr(name string) []string {
	results := make([]string, 0)
	for _, n := range ns {
		// Only include one result per element
		var found bool
		var result string
		for _, attr := range n.Node.Attr {
			if attr.Key == name {
				found = true
				result = attr.Val
			}
		}
		if found {
			results = append(results, result)
		}
	}
	return results
}
