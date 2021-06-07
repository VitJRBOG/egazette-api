package db

import (
	"database/sql"
	"fmt"
	"log"
	"runtime/debug"
	"strings"

	config "github.com/VitJRBOG/RSSMaker/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

func Connect(connectionData config.DBConn) (*sql.DB, error) {
	c := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		connectionData.Login, connectionData.Password,
		connectionData.Address, connectionData.DBName)
	db, err := sql.Open("mysql", c)
	if err != nil {
		return nil, err
	}

	return db, nil
}

type Feed struct {
	ID         int
	SourceName string
	URL        string
}

func (f *Feed) SelectFrom(dbase *sql.DB) ([]Feed, error) {
	query := "SELECT * FROM feed"

	var values []interface{}

	if f.ID != 0 {
		query += " WHERE id = ?"
		values = append(values, f.ID)
	}

	if f.SourceName != "" {
		if strings.Contains(query, "WHERE") {
			query += "AND source_name = ?"
		} else {
			query += " WHERE source_name = ?"
		}
		values = append(values, f.SourceName)
	}

	if f.URL != "" {
		if strings.Contains(query, "WHERE") {
			query += "AND url = ?"
		} else {
			query += " WHERE url = ?"
		}
		values = append(values, f.URL)
	}

	rows, err := dbase.Query(query, values...)
	if err != nil {
		return []Feed{}, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("%s\n%s\n", err, debug.Stack())
		}
	}()

	var feeds []Feed

	for rows.Next() {
		var feed Feed

		if err := rows.Scan(&feed.ID, &feed.SourceName, &feed.URL); err != nil {
			return []Feed{}, err
		}

		feeds = append(feeds, feed)
	}

	if err := rows.Err(); err != nil {
		return []Feed{}, err
	}

	return feeds, nil
}

type VKAccess struct {
	ID          int
	FeedID      int
	AccessToken string
	VKID        int
}

func (v *VKAccess) SelectFrom(dbase *sql.DB) ([]VKAccess, error) {
	query := "SELECT * FROM vk_access"

	var values []interface{}

	if v.ID != 0 {
		query += " WHERE id = ?"
		values = append(values, v.ID)
	}

	if v.FeedID != 0 {
		if strings.Contains(query, "WHERE") {
			query += "AND feed_id = ?"
		} else {
			query += " WHERE feed_id = ?"
		}
		values = append(values, v.FeedID)
	}

	if v.AccessToken != "" {
		if strings.Contains(query, "WHERE") {
			query += "AND access_token = ?"
		} else {
			query += " WHERE access_token = ?"
		}
		values = append(values, v.AccessToken)
	}

	if v.VKID != 0 {
		if strings.Contains(query, "WHERE") {
			query += "AND vk_id = ?"
		} else {
			query += " WHERE vk_id = ?"
		}
		values = append(values, v.VKID)
	}

	rows, err := dbase.Query(query, values...)
	if err != nil {
		return []VKAccess{}, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("%s\n%s\n", err, debug.Stack())
		}
	}()

	var vkAccesses []VKAccess

	for rows.Next() {
		var vkAccess VKAccess

		if err := rows.Scan(&vkAccess.ID, &vkAccess.FeedID,
			&vkAccess.AccessToken, &vkAccess.VKID); err != nil {
			return []VKAccess{}, err
		}

		vkAccesses = append(vkAccesses, vkAccess)
	}

	if err := rows.Err(); err != nil {
		return []VKAccess{}, err
	}

	return vkAccesses, nil
}

func AddNewVKSource(f Feed, v VKAccess, dbase *sql.DB) (Feed, VKAccess, error) {
	tx, err := dbase.Begin()
	if err != nil {
		return Feed{}, VKAccess{}, err
	}

	insertIntoFeed := "INSERT INTO feed(source_name, url) VALUES(?, ?)"

	resultFeed, err := tx.Exec(insertIntoFeed, f.SourceName, f.URL)
	if err != nil {
		return Feed{}, VKAccess{}, err
	}

	feedID, err := resultFeed.LastInsertId()
	if err != nil {
		return Feed{}, VKAccess{}, err
	}

	insertIntoVkAccess := "INSERT INTO vk_access(feed_id, access_token, vk_id) VALUES(?, ?, ?)"

	vkAccessResult, err := tx.Exec(insertIntoVkAccess, feedID, v.AccessToken, v.VKID)
	if err != nil {
		return Feed{}, VKAccess{}, err
	}

	vkAccessID, err := vkAccessResult.LastInsertId()
	if err != nil {
		return Feed{}, VKAccess{}, err
	}

	err = tx.Commit()
	if err != nil {
		return Feed{}, VKAccess{}, err
	}

	f.ID = int(feedID)
	v.ID = int(vkAccessID)

	return f, v, nil
}
