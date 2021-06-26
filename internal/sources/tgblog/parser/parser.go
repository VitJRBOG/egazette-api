package parser

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"golang.org/x/net/html"
)

//
// NOTE
// Parser for https://telegram.org/blog
//

type Article struct {
	Link        string
	Date        string
	Title       string
	Description string
	CoverURL    string
}

func (a *Article) composeInfo(tag *html.Node) {
	a.extractLink(tag)
	a.extractDate(tag)
	a.extractTitle(tag)
	a.extractDescription(tag)
	a.extractCoverURL(tag)
}

func (a *Article) extractLink(tag *html.Node) {
	a.Link = tag.Attr[1].Val
}

func (a *Article) extractCoverURL(tag *html.Node) {
	a.CoverURL = tag.FirstChild.FirstChild.NextSibling.Attr[1].Val
}

func (a *Article) extractTitle(doc *html.Node) {
	tag := findTag("dev_blog_card_title", doc)

	var buf bytes.Buffer
	w := io.Writer(&buf)
	err := html.Render(w, tag)
	if err != nil {
		log.Printf("%s\n%s\n\n", err.Error(), debug.Stack())
		return
	}

	a.Title = buf.String()

	i := strings.Index(a.Title, ">")
	a.Title = a.Title[i+1:]

	i = strings.Index(a.Title, "<")
	a.Title = a.Title[:i]
}

func (a *Article) extractDescription(doc *html.Node) {
	tag := findTag("dev_blog_card_lead", doc)

	var buf bytes.Buffer
	w := io.Writer(&buf)
	err := html.Render(w, tag)
	if err != nil {
		log.Printf("%s\n%s\n\n", err.Error(), debug.Stack())
		return
	}

	a.Description = buf.String()

	i := strings.Index(a.Description, ">")
	a.Description = a.Description[i+1:]

	i = strings.Index(a.Description, "<")
	a.Description = a.Description[:i]
}

func (a *Article) extractDate(doc *html.Node) {
	tag := findTag("dev_blog_card_date", doc)

	var buf bytes.Buffer
	w := io.Writer(&buf)
	err := html.Render(w, tag)
	if err != nil {
		log.Printf("%s\n%s\n\n", err.Error(), debug.Stack())
		return
	}

	a.Date = buf.String()

	i := strings.Index(a.Date, ">")
	a.Date = a.Date[i+1:]

	i = strings.Index(a.Date, "<")
	a.Date = a.Date[:i]
}

func GetArticles(u string) ([]Article, error) {
	doc, err := fetchHTMLNode(u)
	if err != nil {
		return []Article{}, err
	}

	externalTag := findTag("tl_blog_list_cards_wrap", doc)

	articles := composeArticles(externalTag)

	return articles, nil
}

func findTag(tagName string, doc *html.Node) *html.Node {
	return getElementByClass(doc, tagName)
}

func composeArticles(externalTag *html.Node) []Article {
	var articles []Article

	articles = extractTagAttributes(articles, externalTag.FirstChild.NextSibling)

	return articles
}

func extractTagAttributes(articles []Article, articleTag *html.Node) []Article {
	for {
		if articleTag == nil {
			break
		}
		if articleTag.Attr == nil {
			break
		}
		var a Article
		a.composeInfo(articleTag)
		articles = append(articles, a)

		articleTag = articleTag.NextSibling
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
