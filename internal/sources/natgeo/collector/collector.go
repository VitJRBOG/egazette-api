package collector

import (
	"log"
	"runtime/debug"
	"time"

	rss "RSSFeeder/internal/rss"
	natgeoparser "RSSFeeder/internal/sources/natgeo/parser"
)

func ComposeRSS(articles []*natgeoparser.Article) (rss.RSS, error) {
	var r rss.RSS

	r = getSourceInfo(r)
	r = parseArticles(r, articles)

	return r, nil
}

func getSourceInfo(r rss.RSS) rss.RSS {
	r.Channel.Title = "National Geographic"
	r.Channel.Link = "https://www.nationalgeographic.com"

	return r
}

func parseArticles(r rss.RSS, articles []*natgeoparser.Article) rss.RSS {
	for _, article := range articles {
		var rssItem = rss.Item{
			Title:       article.Title,
			Link:        article.Link,
			Description: article.Description,
			Date:        getDateInReadableFormat(article.Date),
		}

		r.Channel.Items = append(r.Channel.Items, rssItem)
	}

	return r
}

func getDateInReadableFormat(date string) string {
	dateFormat := "January 2, 2006"
	t, err := time.Parse(dateFormat, date)
	if err != nil {
		log.Printf("\n%s\n%s", err.Error(), debug.Stack())
		return ""
	}

	d := "Mon, Jan 2 2006 15:04:05 -0700"
	return t.Format(d)
}
