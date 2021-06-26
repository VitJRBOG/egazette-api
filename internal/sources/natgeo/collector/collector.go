package collector

import (
	"log"
	"runtime/debug"
	"time"

	rss "github.com/VitJRBOG/RSSMaker/internal/rss"
	natgeoparser "github.com/VitJRBOG/RSSMaker/internal/sources/natgeo/parser"
)

func ComposeRSS(articles []natgeoparser.Article) (rss.RSS, error) {
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

func parseArticles(r rss.RSS, articles []natgeoparser.Article) rss.RSS {
	for _, article := range articles {
		var rssItem rss.Item

		rssItem.Title = article.Title
		rssItem.Link = article.Link
		rssItem.Description = article.Description
		rssItem.Date = getDateInReadableFormat(article.Date)

		r.Channel.Items = append(r.Channel.Items, rssItem)
	}

	return r
}

func getDateInReadableFormat(date string) string {
	dateFormat := "January 2, 2006"
	t, err := time.Parse(dateFormat, date)
	if err != nil {
		log.Printf("%s\n%s\n\n", err.Error(), debug.Stack())
		return ""
	}

	d := "Mon, Jan 2 2006 15:04:05 -0700"
	return t.Format(d)
}
