package rss

// RSS stores data about the RSS feed.
type RSS struct {
	Channel Channel `xml:"channel"`
}

// Channel stores data about the RSS channel.
type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

// Item stores data about the RSS feed item.
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Date        string `xml:"pubDate"`
}
