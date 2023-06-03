package models

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
