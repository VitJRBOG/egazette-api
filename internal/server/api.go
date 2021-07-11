package server

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"

	db "github.com/VitJRBOG/RSSFeeder/internal/db"
	rss "github.com/VitJRBOG/RSSFeeder/internal/rss"
	natgeocollector "github.com/VitJRBOG/RSSFeeder/internal/sources/natgeo/collector"
	natgeoparser "github.com/VitJRBOG/RSSFeeder/internal/sources/natgeo/parser"
	tgblogcollector "github.com/VitJRBOG/RSSFeeder/internal/sources/tgblog/collector"
	tgblogparser "github.com/VitJRBOG/RSSFeeder/internal/sources/tgblog/parser"
	vkapi "github.com/VitJRBOG/RSSFeeder/internal/sources/vk/api"
	vkcollector "github.com/VitJRBOG/RSSFeeder/internal/sources/vk/collector"
)

func getRSSFeed(dbase *sql.DB, id int) ([]byte, Error) {
	var source = db.Source{
		ID: id,
	}

	sources, err := source.SelectFrom(dbase)
	if err != nil {
		log.Print(err.Error())
		return nil, Error{
			Message: "error getting source from database",
			Code:    http.StatusInternalServerError,
		}
	}

	if len(sources) == 0 {
		return nil, Error{
			Message: fmt.Sprintf("source with the id = %d was not found", id),
			Code:    http.StatusNotFound,
		}
	}

	var rssFeed rss.RSS

	switch {
	case strings.Contains(sources[0].URL, "vk.com"):
		rssFeed, err = rssFromVk(dbase, sources[0])
		if err != nil {
			log.Print(err.Error())
			return nil, Error{
				Message: "error getting RSS feed from VK source",
				Code:    http.StatusInternalServerError,
			}
		}
	case strings.Contains(sources[0].URL, "nationalgeographic.com"):
		rssFeed, err = rssFromNationalGeographic(sources[0])
		if err != nil {
			log.Print(err.Error())
			return nil, Error{
				Message: "error getting RSS feed from National Geographic articles",
				Code:    http.StatusInternalServerError,
			}
		}
	case strings.Contains(sources[0].URL, "telegram.org"):
		rssFeed, err = rssFromTelegramBlog(sources[0])
		if err != nil {
			log.Print(err.Error())
			return nil, Error{
				Message: "error getting RSS feed from Telegram Blog articles",
				Code:    http.StatusInternalServerError,
			}
		}
	}

	data, err := xml.Marshal(rssFeed)
	if err != nil {
		log.Printf("\n%s\n%s", err.Error(), debug.Stack())
		return nil, Error{
			Message: "error marshalling of RSS feed",
			Code:    http.StatusInternalServerError,
		}
	}

	return data, Error{}
}

func rssFromVk(dbase *sql.DB, source db.Source) (rss.RSS, error) {
	var vkAccess = db.VKAccess{
		SourceID: source.ID,
	}

	vkAccesses, err := vkAccess.SelectFrom(dbase)
	if err != nil {
		return rss.RSS{}, err
	}

	if len(vkAccesses) == 0 {
		return rss.RSS{}, fmt.Errorf("access token for source \"%s\" was not found", source.Name)
	}

	community, err := vkapi.GetCommunityInfo(url.Values{
		"access_token": {vkAccesses[0].AccessToken},
		"group_ids":    {strconv.Itoa(-(vkAccesses[0].VKID))},
		"fields":       {"description"},
		"v":            {"5.131"},
		"lang":         {"1"},
	})
	if err != nil {
		return rss.RSS{}, err
	}

	wallPosts, err := vkapi.GetWallPosts(url.Values{
		"access_token": {vkAccesses[0].AccessToken},
		"owner_id":     {strconv.Itoa(vkAccesses[0].VKID)},
		"count":        {"10"},
		"filter":       {"all"},
		"v":            {"5.131"},
	})
	if err != nil {
		return rss.RSS{}, err
	}

	rssFeed, err := vkcollector.ComposeRSS(community, wallPosts)
	if err != nil {
		return rss.RSS{}, err
	}

	return rssFeed, nil
}

func rssFromNationalGeographic(source db.Source) (rss.RSS, error) {
	articles, err := natgeoparser.GetArticles(source.URL)
	if err != nil {
		return rss.RSS{}, err
	}

	rssFeed, err := natgeocollector.ComposeRSS(articles)
	if err != nil {
		return rss.RSS{}, err
	}

	return rssFeed, nil
}

func rssFromTelegramBlog(source db.Source) (rss.RSS, error) {
	articles, err := tgblogparser.GetArticles(source.URL)
	if err != nil {
		return rss.RSS{}, err
	}

	rssFeed, err := tgblogcollector.ComposeRSS(articles)
	if err != nil {
		return rss.RSS{}, err
	}

	return rssFeed, nil
}

type VKRSSSource struct {
	SourceName    string `json:"source_name"`
	URL           string `json:"url"`
	VKAccessToken string `json:"access_token"`
	VKID          int    `json:"vk_id"`
}

func addVKRSSSource(dbase *sql.DB, vkRSSSource VKRSSSource) ([]byte, Error) {
	var source = db.Source{
		Name: vkRSSSource.SourceName,
		URL:  vkRSSSource.URL,
	}
	var vkAccess = db.VKAccess{
		AccessToken: vkRSSSource.VKAccessToken,
		VKID:        vkRSSSource.VKID,
	}

	source, vkAccess, err := db.AddNewVKSource(source, vkAccess, dbase)
	if err != nil {
		log.Print(err.Error())
		return nil, Error{
			Message: "error adding a new VK source",
			Code:    http.StatusInternalServerError,
		}
	}

	var values = map[string]interface{}{
		"feed_id": source.ID,
	}

	data, err := json.Marshal(values)
	if err != nil {
		log.Printf("\n%s\n%s", err.Error(), debug.Stack())
		return nil, Error{
			Message: "error adding a new VK source",
			Code:    http.StatusInternalServerError,
		}
	}

	return data, Error{}
}
