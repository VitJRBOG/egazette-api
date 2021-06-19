package collector

import (
	natgeoparser "github.com/VitJRBOG/RSSMaker/internal/natgeo/parser"
	rss "github.com/VitJRBOG/RSSMaker/internal/rss"
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

		r.Channel.Items = append(r.Channel.Items, rssItem)
	}

	return r
}
