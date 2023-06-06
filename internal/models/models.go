package models

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

// Source stores data about the source.
type Source struct {
	Name        string
	Description string
	HomeURL     string
	APIName     string
}

// FindSourceByAPIName searches for an element by its 'APIName'.
func FindSourceByAPIName(sources []Source, apiName string) Source {
	for _, source := range sources {
		if source.APIName == apiName {
			return source
		}
	}

	return Source{}
}

// Article stores data about the article from the source.
type Article struct {
	URL         string
	PubDate     string
	Title       string
	Description string
	CoverURL    string
	AddDate     string
}

// SetDate converts the date str into a unix timestamp str and sets it in the Date field.
func (a *Article) SetDate(referenceDateLayout, referenceDate string) error {
	date, err := time.Parse(referenceDateLayout, referenceDate)
	if err != nil {
		return fmt.Errorf("failed to convert date str '%s' to time.Time on '%s' layout: %s",
			referenceDate, referenceDateLayout, err.Error())
	}

	a.PubDate = strconv.FormatInt((date.Unix()), 10)

	return nil
}

// SortArticlesByDate sorts the 'Article' slice by the 'Date' field.
func SortArticlesByDate(articles []Article) {
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].PubDate < articles[j].PubDate
	})
}
