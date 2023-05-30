package jpl

import (
	"egazette-api/internal/sources"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"golang.org/x/net/html"
)

// TargetURL stores source URL.
const TargetURL = "https://www.jpl.nasa.gov/news"

// Article stores data about the article from the source.
type Article struct {
	URL         string
	Date        string
	Title       string
	Description string
	CoverURL    string
}

func (a *Article) extractURL(tag *html.Node) {
	articleURL := tag.Attr[0].Val

	a.URL = fmt.Sprintf("%s%s", strings.Replace(TargetURL, "news/", "", 1), articleURL)
}

func (a *Article) extractDate(tag *html.Node) {
	a.Date = tag.Data
}

func (a *Article) extractTitle(tag *html.Node) {
	a.Title = tag.Data
}

func (a *Article) extractDescription(tag *html.Node) {
	a.Description = tag.Data
}

func (a *Article) extractCoverURL(tag *html.Node) {
	a.CoverURL = tag.Attr[0].Val
}

func getArticleData() ([]Article, error) {
	dom, err := fetchDOM(TargetURL)
	if err != nil {
		return []Article{}, err
	}

	htmlNode, err := convertToHTMLNode(dom)
	if err != nil {
		return []Article{}, err
	}

	tagOfArticlesList := sources.FindTag(htmlNode, "id", "search_results")

	articles := composeArticles(tagOfArticlesList.FirstChild)

	return articles, nil
}

func composeArticles(tagOfArticleAnnouncement *html.Node) []Article {
	var articles []Article

	for {
		if tagOfArticleAnnouncement == nil {
			break
		}

		article := extractTagAttributes(tagOfArticleAnnouncement)
		articles = append(articles, article)

		tagOfArticleAnnouncement = tagOfArticleAnnouncement.NextSibling
	}

	return articles
}

func extractTagAttributes(tagOfArticleAnnouncement *html.Node) Article {
	article := Article{}

	tagOfArticleURL := tagOfArticleAnnouncement.FirstChild.FirstChild.FirstChild
	article.extractURL(tagOfArticleURL)

	tagOfArticleInfo := tagOfArticleAnnouncement.FirstChild.FirstChild.FirstChild.FirstChild

	tagOfArticleTitle := tagOfArticleInfo.FirstChild.NextSibling.NextSibling.FirstChild.NextSibling.NextSibling.FirstChild
	article.extractTitle(tagOfArticleTitle)

	tagOfArticleDescription := tagOfArticleInfo.FirstChild.NextSibling.NextSibling.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.FirstChild
	article.extractDescription(tagOfArticleDescription)

	tagOfArticleDate := tagOfArticleInfo.FirstChild.NextSibling.NextSibling.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.FirstChild
	article.extractDate(tagOfArticleDate)

	tagOfArticleCoverURL := tagOfArticleInfo.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.FirstChild.FirstChild.FirstChild.FirstChild
	article.extractCoverURL(tagOfArticleCoverURL)

	return article
}

func convertToHTMLNode(dom []byte) (*html.Node, error) {
	htmlNode, err := html.Parse(strings.NewReader(string(dom)))
	if err != nil {
		return nil, fmt.Errorf("\n%s\n%s", err.Error(), debug.Stack())
	}

	return htmlNode, nil
}

func fetchDOM(targetURL string) ([]byte, error) {
	response, err := http.Get(targetURL)
	if err != nil {
		return nil, fmt.Errorf("\n%s\n%s", err.Error(), debug.Stack())
	}

	defer func() {
		if err != nil {
			log.Printf("%s\n%s\n", err, debug.Stack())
		}
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("\n%s\n%s", err.Error(), debug.Stack())
	}

	return body, nil
}
