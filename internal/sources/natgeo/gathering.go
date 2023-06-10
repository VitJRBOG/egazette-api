package natgeo

import (
	"egazette-api/internal/models"
	"egazette-api/internal/sources"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/net/html"
)

// TargetURL stores source URL.
const TargetURL = "https://www.nationalgeographic.com/pages/topic/latest-stories"

// GetArticleData parses articles from the website of source and returns them.
func GetArticleData() ([]models.Article, error) {
	htmlNode, err := sources.GetHTMLNode(TargetURL)
	if err != nil {
		return []models.Article{}, err
	}

	tagOfPromoTileArticles := sources.FindTag(htmlNode, "class", "GridPromoTile")

	promoTileArticlesContainer, err := composeArticles(tagOfPromoTileArticles.FirstChild.FirstChild)
	if err != nil {
		return []models.Article{}, err
	}

	tagOfInfiniteFeedArticles := sources.FindTag(htmlNode, "class", "InfiniteFeedModule")

	infiniteFeedArticlesContainer, err := composeArticles(tagOfInfiniteFeedArticles.FirstChild.FirstChild.FirstChild.FirstChild)
	if err != nil {
		return []models.Article{}, err
	}

	articles := []models.Article{}
	articles = append(articles, promoTileArticlesContainer...)
	articles = append(articles, infiniteFeedArticlesContainer...)

	return articles, nil
}

func composeArticles(tagOfArticleAnnouncement *html.Node) ([]models.Article, error) {
	articles := []models.Article{}

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

	article.Title = extractTitle(tagOfArticleAnnouncement)

	tagOfArticleURL := tagOfArticleAnnouncement.FirstChild.FirstChild.FirstChild
	article.URL = extractURL(tagOfArticleURL)

	article, err := extractDataFromArticlePage(article)
	if err != nil {

	}

	article.AddDate = strconv.FormatInt((time.Now().UTC().Unix()), 10)

	return article, nil
}

func extractURL(tag *html.Node) string {
	return tag.Attr[3].Val
}

func extractTitle(tag *html.Node) string {
	return tag.Attr[1].Val
}

func extractDataFromArticlePage(article models.Article) (models.Article, error) {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	htmlNode, err := sources.GetHTMLNode(article.URL)
	if err != nil {
		return models.Article{}, err
	}

	tagOfArticleDate := sources.FindTag(htmlNode, "type", "application/ld+json").FirstChild

	articleData := map[string]interface{}{}
	err = json.Unmarshal([]byte(tagOfArticleDate.Data), &articleData)
	if err != nil {
		return models.Article{}, fmt.Errorf("failed to unmarshal the '%s' article data: %s", article.URL, err)
	}

	article.Description = extractDescription(articleData)

	article.CoverURL = extractCoverURL(articleData)

	dateLayout, date := extractDate(articleData)

	if dateLayout != "" && date != "" {
		err := article.SetDate(dateLayout, date)
		if err != nil {
			return models.Article{}, fmt.Errorf("failed to extract pub date of the '%s' article: %s",
				article.URL, err)
		}
	}

	return article, nil
}

func extractDate(articleData map[string]interface{}) (string, string) {
	dateLayout := "2006-01-02T15:04:05.000Z"

	date := ""

	if value, ok := articleData["datePublished"]; ok {
		if date, ok = value.(string); ok {
			return dateLayout, date
		}
	}

	return "", date
}

func extractDescription(articleData map[string]interface{}) string {
	if value, ok := articleData["description"]; ok {
		if description, ok := value.(string); ok {
			return description
		}
	}

	return ""
}

func extractCoverURL(articleData map[string]interface{}) string {
	if innerMap, ok := articleData["image"]; ok {
		if innerMap, ok := innerMap.(map[string]interface{}); ok {
			if value, ok := innerMap["url"]; ok {
				if coverURL, ok := value.(string); ok {
					return coverURL
				}
			}
		}
	}

	return ""
}
