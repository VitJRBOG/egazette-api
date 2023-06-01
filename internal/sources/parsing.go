package sources

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// GetHTMLNode fetch the DOM at the specified URL and returns it converted to html.Node.
func GetHTMLNode(targetURL string) (*html.Node, error) {
	dom, err := fetchDOM(targetURL)
	if err != nil {
		return nil, err
	}

	htmlNode, err := convertToHTMLNode(dom, targetURL)
	if err != nil {
		return nil, err
	}

	return htmlNode, nil
}

func fetchDOM(targetURL string) ([]byte, error) {
	response, err := http.Get(targetURL)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch a DOM from %s: %s", targetURL, err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read respose.Body of %s: %s", targetURL, err.Error())
	}

	return body, nil
}

func convertToHTMLNode(dom []byte, targetURL string) (*html.Node, error) {
	htmlNode, err := html.Parse(strings.NewReader(string(dom)))
	if err != nil {
		return nil, fmt.Errorf("unable to parse the DOM of %s: %s", targetURL, err.Error())
	}

	return htmlNode, nil
}

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
