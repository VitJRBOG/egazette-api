package rss

import (
	"egazette-api/internal/models"
	"fmt"
	"time"
)

// ComposeRSSFeed fetches article data from source and return it in RSS format.
func ComposeRSSFeed(source models.Source, articles []models.Article) (RSS, error) {
	r := RSS{}

	r = putSourceInfo(r, source)
	r, err := putArticleData(r, source, articles)
	if err != nil {
		return RSS{}, err
	}

	return r, nil
}

func putSourceInfo(r RSS, source models.Source) RSS {
	r.Channel.Title = source.Name
	r.Channel.Link = source.HomeURL

	return r
}

func putArticleData(r RSS, source models.Source, articles []models.Article) (RSS, error) {
	for _, article := range articles {
		date, err := prepareDateForRSS(source.Name, source.DateFormat, article.Date)
		if err != nil {
			return RSS{}, err
		}

		var rssItem = Item{
			Title:       article.Title,
			Link:        article.URL,
			Description: fmt.Sprintf("<img src=\"%s\"><br>%s", article.CoverURL, article.Description),
			Date:        date,
		}

		r.Channel.Items = append(r.Channel.Items, rssItem)
	}

	return r, nil
}

func prepareDateForRSS(sourceName, referenceDateFormat, referenceDate string) (string, error) {
	t, err := time.Parse(referenceDateFormat, referenceDate)
	if err != nil {
		return "", fmt.Errorf("failed conversion %s date to RSS date: %s", sourceName, err.Error())
	}

	d := "Mon, Jan 2 2006 15:04:05 -0700"
	return t.Format(d), nil
}
