package collector

import (
	"fmt"
	"log"
	"runtime/debug"
	"time"

	rss "github.com/VitJRBOG/RSSFeeder/internal/rss"
	tgblogparser "github.com/VitJRBOG/RSSFeeder/internal/sources/tgblog/parser"
)

func ComposeRSS(articles []tgblogparser.Article) (rss.RSS, error) {
	var r rss.RSS

	r = getSourceInfo(r)
	r = parseArticles(r, articles)

	return r, nil
}

func getSourceInfo(r rss.RSS) rss.RSS {
	r.Channel.Title = "Telegram blog"
	r.Channel.Link = "https://telegram.org"

	return r
}

func parseArticles(r rss.RSS, articles []tgblogparser.Article) rss.RSS {
	for _, article := range articles {
		var rssItem = rss.Item{
			Title: article.Title,
			Link:  r.Channel.Link + article.Link,
			Description: fmt.Sprintf("<img src=\"%s\"><br>%s",
				article.CoverURL, article.Description),
			Date: getDateInReadableFormat(article.Date),
		}

		r.Channel.Items = append(r.Channel.Items, rssItem)
	}

	return r
}

func getDateInReadableFormat(date string) string {
	dateFormat := "Jan 2, 2006"
	t, err := time.Parse(dateFormat, date)
	if err != nil {
		log.Printf("%s\n%s\n\n", err.Error(), debug.Stack())
		return ""
	}

	d := "Mon, Jan 2 2006 15:04:05 -0700"
	return t.Format(d)
}
