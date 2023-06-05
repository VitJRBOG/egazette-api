package rss

import (
	"egazette-api/internal/models"
	"fmt"
	"strconv"
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
		var err error
		date := ""
		if article.Date != "" {
			date, err = prepareUnixTSForRSS(article.Date)
			if err != nil {
				return RSS{}, err
			}
		}

		description := ""
		if article.CoverURL != "" {
			description += fmt.Sprintf("<img src=\"%s\">", article.CoverURL)
		}

		if article.Description != "" {
			description += fmt.Sprintf("<br>%s", article.Description)
		}

		var rssItem = Item{
			Title:       article.Title,
			Link:        article.URL,
			Description: description,
			Date:        date,
		}

		r.Channel.Items = append(r.Channel.Items, rssItem)
	}

	return r, nil
}

func prepareUnixTSForRSS(unixTimeStamp string) (string, error) {
	t, err := strconv.ParseInt(unixTimeStamp, 10, 64)
	if err != nil {
		return "", fmt.Errorf("failed to convert unix timestamp str '%s' to int64: %s",
			unixTimeStamp, err)
	}

	rssDateLayout := "Mon, Jan 2 2006 15:04:05 -0700"
	return time.Unix(t, 0).Format(rssDateLayout), nil
}
