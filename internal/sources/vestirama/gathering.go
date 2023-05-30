package vestirama

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
const TargetURL = "https://vestirama.ru/novosti/"

// Article stores data about the article from the source.
type Article struct {
	URL      string
	Date     string
	Title    string
	CoverURL string
}

func (a *Article) extractURL(tag *html.Node) {
	articleURL := tag.Attr[0].Val

	a.URL = fmt.Sprintf("%s%s", strings.Replace(TargetURL, "novosti/", "", 1), articleURL)
}

func (a *Article) extractDate(tag *html.Node) {
	date := []rune(tag.Data)

	if date[len(date)-1] == '\'' {
		date = date[:len(date)-1]
	}

	a.Date = string(date)
}

func (a *Article) extractTitle(tag *html.Node) {
	a.Title = tag.Data
}

func (a *Article) extractCoverURL(tag *html.Node) {
	coverURL := tag.Attr[0].Val

	fullSizeImageURL := a.composeURLToFullSizeCoverImage(coverURL)

	a.CoverURL = strings.Replace(fmt.Sprintf("%s%s", TargetURL, fullSizeImageURL), "novosti/", "", 1)
}

func (a *Article) composeURLToFullSizeCoverImage(coverURL string) string {
	fullSizeImageURL := strings.Replace(coverURL, "/assets/cache_image/", "", 1)

	begin := strings.Index(fullSizeImageURL, "_")
	end := strings.Index(fullSizeImageURL[begin:], ".")

	sizeOfSmallImageInFileName := fullSizeImageURL[begin : begin+end]

	fullSizeImageURL = strings.Replace(fullSizeImageURL, sizeOfSmallImageInFileName, "", 1)

	return fullSizeImageURL
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
