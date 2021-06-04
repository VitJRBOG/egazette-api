package collector

import (
	"fmt"
	"time"

	rss "bitbucket.org/VitJRBOG/rss_maker/internal/rss"
	vkapi "bitbucket.org/VitJRBOG/rss_maker/internal/vk/api"
)

func ComposeRSS(community vkapi.Community, wallPosts []vkapi.WallPost) (rss.RSS, error) {
	var r rss.RSS

	r.Channel.Title = community.Name
	r.Channel.Link = "https://vk.com/" + community.ScreenName
	r.Channel.Description = community.Description

	for _, wallPost := range wallPosts {
		var rssItem rss.Item

		rssItem.Title = getWallPostTitle(wallPost.Text)
		rssItem.Description = wallPost.Text
		var err error
		rssItem.Date, err = getDateInReadableFormat(int64(wallPost.Date))
		if err != nil {
			return rss.RSS{}, err
		}
		rssItem.Link = fmt.Sprintf("https://vk.com/wall%d_%d", wallPost.OwnerID, wallPost.ID)
		if len(wallPost.Attachments) > 0 {
			for _, attachment := range wallPost.Attachments {
				if attachment.Type != "photo" {
					continue
				}

				// FIXME: нужно указывать прямую ссылку на фото
				// c https://sun1-87.userapi.com и прочим,
				// иначе ридер ее не видит
				rssItem.Enclosure.URL = fmt.Sprintf("https://vk.com/photo%d_%d",
					attachment.Photo.OwnerID, attachment.Photo.ID)
				rssItem.Enclosure.Type = "image/jpeg"
			}
		}

		r.Channel.Items = append(r.Channel.Items, rssItem)
	}

	return r, nil
}

func getWallPostTitle(text string) string {
	runes := []rune(text)

	if len(runes) == 0 {
		return text
	}

	for i, r := range runes {
		if i == 80 {
			for j := i; j > 0; j-- {
				if runes[j] == ' ' {
					return fmt.Sprintf("%s...", string(runes[:j]))
				}
			}

			return fmt.Sprintf("%s...", string(runes[:77]))
		}

		if r == '\n' {
			return string(runes[:i])
		}
	}

	return text
}

func getDateInReadableFormat(ts int64) (string, error) {
	loc, err := time.LoadLocation("Asia/Yekaterinburg")
	if err != nil {
		return "", err
	}
	t := time.Unix(ts, 0).In(loc)
	dateFormat := "Mon, Jan 2 2006 15:04:05 -0700"
	return t.Format(dateFormat), nil
}
