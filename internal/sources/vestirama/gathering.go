package vestirama

import (
	"egazette-api/internal/models"
	"egazette-api/internal/sources"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// TargetURL stores source URL.
const TargetURL = "https://vestirama.ru/novosti/"

// GetSourceData returns a struct with data about source.
func GetSourceData() models.Source {
	return models.Source{
		Name:       "Vestirama",
		HomeURL:    "https://vestirama.ru",
		DateFormat: "2.01.2006 15:04",
	}
}

// GetArticleData parses articles from the website of source and returns them.
func GetArticleData() ([]models.Article, error) {
	htmlNode, err := sources.GetHTMLNode(TargetURL)
	if err != nil {
		return []models.Article{}, err
	}

	tagOfArticlesList := sources.FindTag(htmlNode, "class", "news").FirstChild.NextSibling

	articles := composeArticles(tagOfArticlesList.FirstChild.NextSibling)

	return articles, nil
}

func composeArticles(tagOfArticleAnnouncement *html.Node) []models.Article {
	var articles []models.Article

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

func extractTagAttributes(tagOfArticleAnnouncement *html.Node) models.Article {
	article := models.Article{}

	tagOfArticleCoverURL := tagOfArticleAnnouncement.FirstChild.NextSibling.FirstChild
	article.CoverURL = extractCoverURL(tagOfArticleCoverURL)

	tagOfArticleTextInfo := tagOfArticleAnnouncement.FirstChild.NextSibling.NextSibling.NextSibling

	tagOfArticleURL := tagOfArticleTextInfo.FirstChild.NextSibling.NextSibling.NextSibling
	article.URL = extractURL(tagOfArticleURL)

	tagOfArticleTitle := tagOfArticleTextInfo.FirstChild.NextSibling.NextSibling.NextSibling.FirstChild
	article.Title = extractTitle(tagOfArticleTitle)

	tagOfArticleDate := tagOfArticleTextInfo.FirstChild.NextSibling.FirstChild.NextSibling.FirstChild
	article.Date = extractDate(tagOfArticleDate)

	return article
}

func extractURL(tag *html.Node) string {
	articleURL := tag.Attr[0].Val

	return fmt.Sprintf("%s%s", strings.Replace(TargetURL, "novosti/", "", 1), articleURL)
}

func extractDate(tag *html.Node) string {
	date := strings.Trim(tag.Data, "'")
	date = strings.TrimLeft(date, " ")

	return date
}

func extractTitle(tag *html.Node) string {
	return tag.Data
}

func extractCoverURL(tag *html.Node) string {
	coverURL := tag.Attr[0].Val

	fullSizeImageURL := composeURLToFullSizeCoverImage(coverURL)

	return strings.Replace(fmt.Sprintf("%s%s", TargetURL, fullSizeImageURL), "novosti/", "", 1)
}

func composeURLToFullSizeCoverImage(coverURL string) string {
	fullSizeImageURL := strings.Replace(coverURL, "/assets/cache_image/", "", 1)

	begin := strings.Index(fullSizeImageURL, "_")
	end := strings.Index(fullSizeImageURL[begin:], ".")

	sizeOfSmallImageInFileName := fullSizeImageURL[begin : begin+end]

	fullSizeImageURL = strings.Replace(fullSizeImageURL, sizeOfSmallImageInFileName, "", 1)

	return fullSizeImageURL
}
