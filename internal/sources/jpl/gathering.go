package jpl

import (
	"egazette-api/internal/models"
	"egazette-api/internal/sources"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// TargetURL stores source URL.
const TargetURL = "https://www.jpl.nasa.gov/news"

// GetSourceData returns a struct with data about source.
func GetSourceData() models.Source {
	return models.Source{
		Name:       "Jet Propulsion Laboratory",
		HomeURL:    "https://www.jpl.nasa.gov",
		DateFormat: "January 2, 2006",
	}
}

// GetArticleData parses articles from the website of source and returns them.
func GetArticleData() ([]models.Article, error) {
	htmlNode, err := sources.GetHTMLNode(TargetURL)
	if err != nil {
		return []models.Article{}, err
	}

	tagOfArticlesList := sources.FindTag(htmlNode, "id", "search_results")

	articles := composeArticles(tagOfArticlesList.FirstChild)

	return articles, nil
}

func composeArticles(tagOfArticleAnnouncement *html.Node) []models.Article {
	var articles []models.Article

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

func extractTagAttributes(tagOfArticleAnnouncement *html.Node) models.Article {
	article := models.Article{}

	tagOfArticleURL := tagOfArticleAnnouncement.FirstChild.FirstChild.FirstChild
	article.URL = extractURL(tagOfArticleURL)

	tagOfArticleInfo := tagOfArticleAnnouncement.FirstChild.FirstChild.FirstChild.FirstChild

	tagOfArticleTitle := tagOfArticleInfo.FirstChild.NextSibling.NextSibling.FirstChild.NextSibling.NextSibling.FirstChild
	article.Title = extractTitle(tagOfArticleTitle)

	tagOfArticleDescription := tagOfArticleInfo.FirstChild.NextSibling.NextSibling.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.FirstChild
	article.Description = extractDescription(tagOfArticleDescription)

	tagOfArticleDate := tagOfArticleInfo.FirstChild.NextSibling.NextSibling.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.FirstChild
	article.Date = extractDate(tagOfArticleDate)

	tagOfArticleCoverURL := tagOfArticleInfo.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.FirstChild.FirstChild.FirstChild.FirstChild
	article.CoverURL = extractCoverURL(tagOfArticleCoverURL)

	return article
}

func extractURL(tag *html.Node) string {
	articleURL := tag.Attr[0].Val

	return fmt.Sprintf("%s%s", strings.Replace(TargetURL, "/news", "", 1), articleURL)
}

func extractDate(tag *html.Node) string {
	return tag.Data
}

func extractTitle(tag *html.Node) string {
	return tag.Data
}

func extractDescription(tag *html.Node) string {
	return tag.Data
}

func extractCoverURL(tag *html.Node) string {
	return tag.Attr[0].Val
}
