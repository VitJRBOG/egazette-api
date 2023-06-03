package models

import (
	"fmt"
	"strconv"
	"time"
)

// Source stores data about the source.
type Source struct {
	Name       string
	HomeURL    string
	DateFormat string
}

// Article stores data about the article from the source.
type Article struct {
	URL         string
	Date        string
	Title       string
	Description string
	CoverURL    string
}

// SetDate converts the date str into a unix timestamp str and sets it in the Date field.
func (a *Article) SetDate(referenceDateLayout, referenceDate string) error {
	date, err := time.Parse(referenceDateLayout, referenceDate)
	if err != nil {
		return fmt.Errorf("failed to convert date str '%s' to time.Time on '%s' layout: %s",
			referenceDate, referenceDateLayout, err.Error())
	}

	a.Date = strconv.FormatInt((date.Unix()), 10)

	return nil
}
