package sources

import (
	"strings"

	"golang.org/x/net/html"
)

// FindTag looks for an HTML tag to the specified class name.
func FindTag(className string, doc *html.Node) *html.Node {
	return getElementByClass(doc, className)
}

func getElementByClass(n *html.Node, k string) *html.Node {
	return traverse(n, k)
}

func traverse(n *html.Node, k string) *html.Node {
	if checkKey(n, k) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := traverse(c, k)
		if result != nil {
			return result
		}
	}

	return nil
}

func checkKey(n *html.Node, k string) bool {
	if n.Type == html.ElementNode {
		s, ok := getAttribute(n, "class")
		if ok && strings.Contains(s, k) {
			return true
		}
	}
	return false
}

func getAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}
