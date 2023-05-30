package sources

import (
	"strings"

	"golang.org/x/net/html"
)

// FindTag looks for an HTML tag to the specified attribute key and value.
func FindTag(htmlNode *html.Node, attrKey, attrValue string) *html.Node {
	return traverse(htmlNode, attrKey, attrValue)
}

func traverse(htmlNode *html.Node, attrKey, attrValue string) *html.Node {
	if checkKey(htmlNode, attrKey, attrValue) {
		return htmlNode
	}

	for c := htmlNode.FirstChild; c != nil; c = c.NextSibling {
		result := traverse(c, attrKey, attrValue)
		if result != nil {
			return result
		}
	}

	return nil
}

func checkKey(htmlNode *html.Node, attrKey, attrValue string) bool {
	if htmlNode.Type == html.ElementNode {
		s, ok := getAttribute(htmlNode, attrKey)
		if ok && strings.Contains(s, attrValue) {
			return true
		}
	}
	return false
}

func getAttribute(htmlNode *html.Node, attrKey string) (string, bool) {
	for _, attr := range htmlNode.Attr {
		if attr.Key == attrKey {
			return attr.Val, true
		}
	}
	return "", false
}
