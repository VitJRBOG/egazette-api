package vestirama

import (
	"egazette-api/internal/sources"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// TargetURL stores source URL.
const TargetURL = "https://vestirama.ru/novosti/"

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
		name:       "Vestirama",
		homeURL:    "https://vestirama.ru",
		dateFormat: "2.01.2006 15:04",
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

	a.url = fmt.Sprintf("%s%s", strings.Replace(TargetURL, "novosti/", "", 1), articleURL)
}

func (a *Article) extractDate(tag *html.Node) {
	date := strings.Trim(tag.Data, "'")
	date = strings.TrimLeft(date, " ")

	a.date = date
}

func (a *Article) extractTitle(tag *html.Node) {
	a.title = tag.Data
}

func (a *Article) extractCoverURL(tag *html.Node) {
	coverURL := tag.Attr[0].Val

	fullSizeImageURL := a.composeURLToFullSizeCoverImage(coverURL)

	a.coverURL = strings.Replace(fmt.Sprintf("%s%s", TargetURL, fullSizeImageURL), "novosti/", "", 1)
}

func (a *Article) composeURLToFullSizeCoverImage(coverURL string) string {
	fullSizeImageURL := strings.Replace(coverURL, "/assets/cache_image/", "", 1)

	begin := strings.Index(fullSizeImageURL, "_")
	end := strings.Index(fullSizeImageURL[begin:], ".")

	sizeOfSmallImageInFileName := fullSizeImageURL[begin : begin+end]

	fullSizeImageURL = strings.Replace(fullSizeImageURL, sizeOfSmallImageInFileName, "", 1)

	return fullSizeImageURL
}

// GetArticleData parses articles from the website of source and returns them.
func GetArticleData() ([]Article, error) {
	htmlNode, err := sources.GetHTMLNode(TargetURL)
	if err != nil {
		return []Article{}, err
	}

	tagOfArticlesList := sources.FindTag(htmlNode, "class", "news").FirstChild.NextSibling

	articles := composeArticles(tagOfArticlesList.FirstChild.NextSibling)

	return articles, nil
}

func composeArticles(tagOfArticleAnnouncement *html.Node) []Article {
	var articles []Article

	for {
		if tagOfArticleAnnouncement == nil {
			break
		}
		if tagOfArticleAnnouncement.Attr[0].Val != "news-item __" {
			break
		}

		article := extractTagAttributes(tagOfArticleAnnouncement)
		articles = append(articles, article)

		tagOfArticleAnnouncement = tagOfArticleAnnouncement.NextSibling.NextSibling
	}

	return articles
}

func extractTagAttributes(tagOfArticleAnnouncement *html.Node) Article {
	article := Article{}

	tagOfArticleCoverURL := tagOfArticleAnnouncement.FirstChild.NextSibling.FirstChild
	article.extractCoverURL(tagOfArticleCoverURL)

	tagOfArticleTextInfo := tagOfArticleAnnouncement.FirstChild.NextSibling.NextSibling.NextSibling

	tagOfArticleURL := tagOfArticleTextInfo.FirstChild.NextSibling.NextSibling.NextSibling
	article.extractURL(tagOfArticleURL)

	tagOfArticleTitle := tagOfArticleTextInfo.FirstChild.NextSibling.NextSibling.NextSibling.FirstChild
	article.extractTitle(tagOfArticleTitle)

	tagOfArticleDate := tagOfArticleTextInfo.FirstChild.NextSibling.FirstChild.NextSibling.FirstChild
	article.extractDate(tagOfArticleDate)

	return article
}
