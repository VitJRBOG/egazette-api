package parser

import (
	"io/ioutil"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"golang.org/x/net/html"
)

type Article struct {
	Link        string
	Date        string
	Title       string
	Description string
}

func (a *Article) composeInfo(articleTag *html.Node) {
	a.Title = articleTag.Attr[2].Val
	a.Link = articleTag.Attr[3].Val
}

func GetArticles(u string) ([]Article, error) {
	doc, err := fetchHTMLNode(u)
	if err != nil {
		return []Article{}, err
	}

	externalTag := extractExternalTag(doc)

	articles := composeArticles(externalTag)

	return articles, nil
}

func extractExternalTag(doc *html.Node) *html.Node {
	return getElementByClass(doc, "FrameBackgroundFull FrameBackgroundFull--grey")
}

func composeArticles(externalTag *html.Node) []Article {
	var articles []Article

	articles = extractTagAttributes(articles, externalTag)
	// articles = extractBigTagAttrubites(articles, externalTag.NextSibling)
	articles = extractTagAttributes(articles, externalTag.NextSibling.NextSibling)
	articles = extractLastTagAttributes(articles, externalTag.NextSibling.NextSibling)

	return articles
}

func extractTagAttributes(articles []Article, externalTag *html.Node) []Article {
	commonTag := externalTag.FirstChild.NextSibling.FirstChild.FirstChild.FirstChild
	for {
		if commonTag == nil {
			break
		}
		articleTag := commonTag.FirstChild.FirstChild.FirstChild
		var a Article
		a.composeInfo(articleTag)
		articles = append(articles, a)

		commonTag = commonTag.NextSibling
	}
	return articles
}

// func extractBigTagAttrubites(articles []Article, externalTag *html.Node) []Article {
// 	articleURLTag := getElementByClass(externalTag, "AnchorLink BgImagePromo__Container__Text__Link")
// 	if articleURLTag != nil {
// 		articleTitleTag := articleURLTag.NextSibling.NextSibling

// 		var a Article

// 		// TODO: найти способ получать текст из тега <h2>
// 		a.Title = fmt.Sprintf("%v", articleTitleTag)
// 		a.Link = articleURLTag.Attr[5].Val
// 		articles = append(articles, a)
// 	}

// 	return articles
// }

func extractLastTagAttributes(articles []Article, externalTag *html.Node) []Article {
	commonTag := externalTag.FirstChild.NextSibling.FirstChild.NextSibling.FirstChild.FirstChild
	for {
		if commonTag == nil {
			break
		}
		articleTag := commonTag.FirstChild.FirstChild.FirstChild
		var a Article
		a.composeInfo(articleTag)
		articles = append(articles, a)

		commonTag = commonTag.NextSibling
	}
	return articles
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

func fetchHTMLNode(u string) (*html.Node, error) {
	body, err := sendRequest(u)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func sendRequest(u string) ([]byte, error) {
	response, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			log.Printf("%s\n%s\n", err, debug.Stack())
		}
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
