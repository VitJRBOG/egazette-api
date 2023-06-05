package natgeo

import (
	"egazette-api/internal/models"
	"egazette-api/internal/sources"

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

	// TODO: describe the extraction of the date of publication
	// TODO: describe the extraction of the article cover URL
	// TODO: describe the extraction of the description

	return article, nil
}

func extractURL(tag *html.Node) string {
	return tag.Attr[3].Val
}

func extractTitle(tag *html.Node) string {
	return tag.Attr[1].Val
}
