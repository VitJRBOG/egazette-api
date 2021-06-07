package server

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"
	"strconv"

	db "github.com/VitJRBOG/RSSMaker/internal/db"
	vkapi "github.com/VitJRBOG/RSSMaker/internal/vk/api"
	vkcollector "github.com/VitJRBOG/RSSMaker/internal/vk/collector"
)

func getRSSFeed(dbase *sql.DB, id int) ([]byte, error) {
	var feed = db.Feed{
		ID: id,
	}

	feeds, err := feed.SelectFrom(dbase)
	if err != nil {
		return nil, err
	}

	if len(feeds) == 0 {
		return nil, fmt.Errorf("feed with the id = %d was not found", id)
	}

	var vkAccess = db.VKAccess{
		FeedID: feed.ID,
	}

	vkAccesses, err := vkAccess.SelectFrom(dbase)
	if err != nil {
		return nil, err
	}

	if len(vkAccesses) == 0 {
		return nil, fmt.Errorf("access token for feed \"%s\" was not found", feed.SourceName)
	}

	community, err := vkapi.GetCommunityInfo(url.Values{
		"access_token": {vkAccesses[0].AccessToken},
		"group_ids":    {strconv.Itoa(-(vkAccesses[0].VKID))},
		"fields":       {"description"},
		"v":            {"5.131"},
		"lang":         {"1"},
	})
	if err != nil {
		return nil, err
	}

	wallPosts, err := vkapi.GetWallPosts(url.Values{
		"access_token": {vkAccesses[0].AccessToken},
		"owner_id":     {strconv.Itoa(vkAccesses[0].VKID)},
		"count":        {"10"},
		"filter":       {"all"},
		"v":            {"5.131"},
	})
	if err != nil {
		return nil, err
	}

	rssFeed, err := vkcollector.ComposeRSS(community, wallPosts)
	if err != nil {
		return nil, err
	}

	data, err := xml.Marshal(rssFeed)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type VKRSSSource struct {
	SourceName    string `json:"source_name"`
	URL           string `json:"url"`
	VKAccessToken string `json:"access_token"`
	VKID          int    `json:"vk_id"`
}

func addVKRSSSource(dbase *sql.DB, vkRSSSource VKRSSSource) ([]byte, error) {
	var feed = db.Feed{
		SourceName: vkRSSSource.SourceName,
		URL:        vkRSSSource.URL,
	}
	var vkAccess = db.VKAccess{
		AccessToken: vkRSSSource.VKAccessToken,
		VKID:        vkRSSSource.VKID,
	}

	feed, vkAccess, err := db.AddNewVKSource(feed, vkAccess, dbase)
	if err != nil {
		return nil, err
	}

	var values = map[string]interface{}{
		"feed_id": feed.ID,
	}

	data, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	return data, nil
}
