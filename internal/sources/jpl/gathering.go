package jpl

import (
	"egazette-api/internal/sources"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// TargetURL stores source URL.
const TargetURL = "https://www.jpl.nasa.gov/news"

// Source stores data about the source.
type Source struct {
	name       string
	homeURL    string
	dateFormat string
}

// Name returns a name of source.
func (s Source) Name() string {
	return s.name
}

// HomeURL returns an URL of the home page of source.
func (s Source) HomeURL() string {
	return s.homeURL
}

// DateFormat returns an format of articles publication dates.
func (s Source) DateFormat() string {
	return s.dateFormat
}

// GetSourceData returns a struct with data about source.
func GetSourceData() Source {
	return Source{
		name:       "Jet Propulsion Laboratory",
		homeURL:    "https://www.jpl.nasa.gov",
		dateFormat: "January 2, 2006",
	}
}

// Article stores data about the article from the source.
type Article struct {
	url         string
	date        string
	title       string
	description string
	coverURL    string
}

// URL returns an URL of the article.
func (a Article) URL() string {
	return a.url
}

// Date returns the date the article was published.
func (a Article) Date() string {
	return a.date
}

// Title returns the title of the article.
func (a Article) Title() string {
	return a.title
}

// Description returns the description of the article.
func (a Article) Description() string {
	return a.description
}

// CoverURL returns an URL of the article cover.
func (a Article) CoverURL() string {
	return a.coverURL
}

func (a *Article) extractURL(tag *html.Node) {
	articleURL := tag.Attr[0].Val

	a.url = fmt.Sprintf("%s%s", strings.Replace(TargetURL, "/news", "", 1), articleURL)
}

func (a *Article) extractDate(tag *html.Node) {
	a.date = tag.Data
}

func (a *Article) extractTitle(tag *html.Node) {
	a.title = tag.Data
}

func (a *Article) extractDescription(tag *html.Node) {
	a.description = tag.Data
}

func (a *Article) extractCoverURL(tag *html.Node) {
	a.coverURL = tag.Attr[0].Val
}

// GetArticleData parses articles from the website of source and returns them.
func GetArticleData() ([]Article, error) {
	htmlNode, err := sources.GetHTMLNode(TargetURL)
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
