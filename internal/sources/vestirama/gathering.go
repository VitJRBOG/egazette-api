package vestirama

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
const TargetURL = "https://vestirama.ru/novosti/"

// GetArticleData parses articles from the website of source and returns them.
func GetArticleData() ([]models.Article, error) {
	htmlNode, err := sources.GetHTMLNode(TargetURL)
	if err != nil {
		return []models.Article{}, err
	}

	tagOfArticlesList := sources.FindTag(htmlNode, "class", "news").FirstChild.NextSibling

	articles, err := composeArticles(tagOfArticlesList.FirstChild.NextSibling)
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
		if tagOfArticleAnnouncement.Attr[0].Val != "news-item __" {
			break
		}

		article, err := extractTagAttributes(tagOfArticleAnnouncement)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)

		tagOfArticleAnnouncement = tagOfArticleAnnouncement.NextSibling.NextSibling
	}

	return articles, nil
}

func extractTagAttributes(tagOfArticleAnnouncement *html.Node) (models.Article, error) {
	article := models.Article{}

	tagOfArticleCoverURL := tagOfArticleAnnouncement.FirstChild.NextSibling.FirstChild
	article.CoverURL = extractCoverURL(tagOfArticleCoverURL)

	tagOfArticleTextInfo := tagOfArticleAnnouncement.FirstChild.NextSibling.NextSibling.NextSibling

	tagOfArticleURL := tagOfArticleTextInfo.FirstChild.NextSibling.NextSibling.NextSibling
	article.URL = extractURL(tagOfArticleURL)

	tagOfArticleTitle := tagOfArticleTextInfo.FirstChild.NextSibling.NextSibling.NextSibling.FirstChild
	article.Title = extractTitle(tagOfArticleTitle)

	tagOfArticleDate := tagOfArticleTextInfo.FirstChild.NextSibling.FirstChild.NextSibling.FirstChild
	err := article.SetDate(extractDate(tagOfArticleDate))
	if err != nil {
		return models.Article{}, fmt.Errorf("failed to extract tag attributes of the '%s' article: %s",
			article.Title, err)
	}

	article.AddDate = strconv.FormatInt((time.Now().UTC().Unix()), 10)

	return article, nil
}

func extractURL(tag *html.Node) string {
	articleURL := tag.Attr[0].Val

	return fmt.Sprintf("%s%s", strings.Replace(TargetURL, "novosti/", "", 1), articleURL)
}

func extractDate(tag *html.Node) (string, string) {
	dateLayout := "2.01.2006 15:04 -0700"

	date := strings.Trim(tag.Data, "'")
	date = strings.TrimLeft(date, " ")
	date += " +0500"

	return dateLayout, date
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
