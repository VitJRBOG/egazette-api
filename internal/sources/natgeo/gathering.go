package natgeo

import (
	"egazette-api/internal/models"
	"egazette-api/internal/sources"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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

	tagOfArticlePubDate := sources.FindTag(htmlNode, "class", "Byline__Meta Byline__Meta--publishDate").FirstChild

	dateLayout, date, err := extractDate(tagOfArticlePubDate)
	if err != nil {
		return models.Article{}, fmt.Errorf("failed to extract the pub date of '%s' article: %s",
			article.URL, err)
	}
	err = article.SetDate(dateLayout, date)
	if err != nil {
		return models.Article{}, fmt.Errorf("failed to convert the pub date of '%s' article: %s",
			article.URL, err)
	}

	tagOfArticleDescription := sources.FindTag(htmlNode, "class", "Article__Headline__Desc")
	article.Description = extractDescription(tagOfArticleDescription.FirstChild)

	tagOfArticleCover := sources.FindTag(htmlNode, "class", "Image__Wrapper ")
	article.CoverURL = extractCoverURL(tagOfArticleCover)

	return article, nil
}

func extractDate(tag *html.Node) (string, string, error) {
	dateLayout := "January 2, 2006 -0700"

	date := fmt.Sprintf("%s +0000", strings.Replace(tag.Data, "Published ", "", 1))

	return dateLayout, date, nil
}

func extractDescription(tag *html.Node) string {
	return tag.Data
}

func extractCoverURL(tag *html.Node) string {
	tagOfArticleCoverURL := tag.FirstChild.LastChild.PrevSibling

	tagAttr := tagOfArticleCoverURL.Attr[0].Val

	return strings.Split(tagAttr, ",")[0]
}
