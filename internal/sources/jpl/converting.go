package jpl

import (
	"egazette-api/internal/rss"
	"fmt"
	"time"
)

// ComposeRSSFeed fetches article data from source and return it in RSS format.
func ComposeRSSFeed() (rss.RSS, error) {
	articles, err := getArticleData()
	if err != nil {
		return rss.RSS{}, err
	}

	r := rss.RSS{}

	r = putSourceInfo(r)
	r, err = putArticleData(r, articles)
	if err != nil {
		return rss.RSS{}, err
	}

	return r, nil
}

func putSourceInfo(r rss.RSS) rss.RSS {
	r.Channel.Title = "Jet Propulsion Laboratory"
	r.Channel.Link = "https://www.jpl.nasa.gov"

	return r
}

func putArticleData(r rss.RSS, articles []Article) (rss.RSS, error) {
	for _, article := range articles {
		date, err := prepareDateForRSS(article.Date)
		if err != nil {
			return rss.RSS{}, err
		}

		var rssItem = rss.Item{
			Title:       article.Title,
			Link:        article.URL,
			Description: fmt.Sprintf("<img src=\"%s\"><br>%s", article.CoverURL, article.Description),
			Date:        date,
		}

		r.Channel.Items = append(r.Channel.Items, rssItem)
	}

	return r, nil
}

func prepareDateForRSS(referenceDate string) (string, error) {
	referenceDateFormat := "January 2, 2006"
	t, err := time.Parse(referenceDateFormat, referenceDate)
	if err != nil {
		return "", fmt.Errorf("failed conversion JPL date to RSS date: %s", err.Error())
	}

	d := "Mon, Jan 2 2006 15:04:05 -0700"
	return t.Format(d), nil
}
