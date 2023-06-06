package jpl

import (
	"egazette-api/internal/models"
	"egazette-api/internal/sources"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// TargetURL stores source URL.
const TargetURL = "https://www.jpl.nasa.gov/news"

// GetArticleData parses articles from the website of source and returns them.
func GetArticleData() ([]models.Article, error) {
	htmlNode, err := sources.GetHTMLNode(TargetURL)
	if err != nil {
		return []models.Article{}, err
	}

	tagOfArticlesList := sources.FindTag(htmlNode, "id", "search_results")

	articles, err := composeArticles(tagOfArticlesList.FirstChild)
	if err != nil {
		return []models.Article{}, err
	}

	models.SortArticlesByDate(articles)

	return articles, nil
}

func composeArticles(tagOfArticleAnnouncement *html.Node) ([]models.Article, error) {
	var articles []models.Article

	for {
		if tagOfArticleAnnouncement == nil {
			break
		}

		article, err := extractTagAttributes(tagOfArticleAnnouncement)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)

		tagOfArticleAnnouncement = tagOfArticleAnnouncement.NextSibling
	}

	return articles, nil
}

func extractTagAttributes(tagOfArticleAnnouncement *html.Node) (models.Article, error) {
	article := models.Article{}

	tagOfArticleURL := tagOfArticleAnnouncement.FirstChild.FirstChild.FirstChild
	article.URL = extractURL(tagOfArticleURL)

	tagOfArticleInfo := tagOfArticleAnnouncement.FirstChild.FirstChild.FirstChild.FirstChild

	tagOfArticleTitle := tagOfArticleInfo.FirstChild.NextSibling.NextSibling.FirstChild.NextSibling.NextSibling.FirstChild
	article.Title = extractTitle(tagOfArticleTitle)

	tagOfArticleDescription := tagOfArticleInfo.FirstChild.NextSibling.NextSibling.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.FirstChild
	article.Description = extractDescription(tagOfArticleDescription)

	tagOfArticleDate := tagOfArticleInfo.FirstChild.NextSibling.NextSibling.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.FirstChild
	err := article.SetDate(extractDate(tagOfArticleDate))
	if err != nil {
		return models.Article{}, fmt.Errorf("failed to extract tag attributes of the '%s' article: %s",
			article.Title, err)
	}

	tagOfArticleCoverURL := tagOfArticleInfo.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.FirstChild.FirstChild.FirstChild.FirstChild
	article.CoverURL = extractCoverURL(tagOfArticleCoverURL)

	article.AddDate = strconv.FormatInt((time.Now().UTC().Unix()), 10)

	return article, nil
}

func extractURL(tag *html.Node) string {
	articleURL := tag.Attr[0].Val

	return fmt.Sprintf("%s%s", strings.Replace(TargetURL, "/news", "", 1), articleURL)
}

func extractDate(tag *html.Node) (string, string) {
	dateLayout := "January 2, 2006 -0700"
	date := fmt.Sprintf("%s +0000", tag.Data)

	return dateLayout, date
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
