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

type Source struct {
	ID   int
	Name string
	URL  string
}

func (s *Source) SelectFrom(dbase *sql.DB) ([]Source, error) {
	query := "SELECT * FROM source"

	var values []interface{}

	if s.ID != 0 {
		query += " WHERE id = ?"
		values = append(values, s.ID)
	}

	if s.Name != "" {
		if strings.Contains(query, "WHERE") {
			query += "AND name = ?"
		} else {
			query += " WHERE name = ?"
		}
		values = append(values, s.Name)
	}

	if s.URL != "" {
		if strings.Contains(query, "WHERE") {
			query += "AND url = ?"
		} else {
			query += " WHERE url = ?"
		}
		values = append(values, s.URL)
	}

	rows, err := dbase.Query(query, values...)
	if err != nil {
		return []Source{}, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("%s\n%s\n", err, debug.Stack())
		}
	}()

	var sources []Source

	for rows.Next() {
		var source Source

		if err := rows.Scan(&source.ID, &source.Name, &source.URL); err != nil {
			return []Source{}, err
		}

		sources = append(sources, source)
	}

	if err := rows.Err(); err != nil {
		return []Source{}, err
	}

	return sources, nil
}

func (s *Source) InsertInto(query string, dbase *sql.DB, tx *sql.Tx) (int, error) {
	var id int64

	switch {
	case dbase == nil && tx != nil:
		result, err := tx.Exec(query, s.Name, s.URL)
		if err != nil {
			return 0, err
		}

		id, err = result.LastInsertId()
		if err != nil {
			return 0, err
		}
	case tx == nil && dbase != nil:
		result, err := dbase.Exec(query, s.Name, s.URL)
		if err != nil {
			return 0, err
		}

		id, err = result.LastInsertId()
		if err != nil {
			return 0, err
		}
	}

	return int(id), nil
}

type VKAccess struct {
	ID          int
	SourceID    int
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

	if v.SourceID != 0 {
		if strings.Contains(query, "WHERE") {
			query += "AND source_id = ?"
		} else {
			query += " WHERE source_id = ?"
		}
		values = append(values, v.SourceID)
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

		if err := rows.Scan(&vkAccess.ID, &vkAccess.SourceID,
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

func (v *VKAccess) InsertInto(query string, dbase *sql.DB, tx *sql.Tx) (int, error) {
	var id int64

	switch {
	case dbase == nil && tx != nil:
		result, err := tx.Exec(query, v.SourceID, v.AccessToken, v.VKID)
		if err != nil {
			return 0, err
		}

		id, err = result.LastInsertId()
		if err != nil {
			return 0, err
		}
	case tx == nil && dbase != nil:
		result, err := dbase.Exec(query, v.SourceID, v.AccessToken, v.VKID)
		if err != nil {
			return 0, err
		}

		id, err = result.LastInsertId()
		if err != nil {
			return 0, err
		}
	}

	return int(id), nil
}

func AddNewVKSource(s Source, v VKAccess, dbase *sql.DB) (Source, VKAccess, error) {
	tx, err := dbase.Begin()
	if err != nil {
		return Source{}, VKAccess{}, err
	}

	insertIntoFeed := "INSERT INTO source(name, url) VALUES(?, ?)"

	feedID, err := s.InsertInto(insertIntoFeed, nil, tx)
	if err != nil {
		return Source{}, VKAccess{}, err
	}

	s.ID = feedID
	v.SourceID = feedID

	insertIntoVkAccess := "INSERT INTO vk_access(source_id, access_token, vk_id) VALUES(?, ?, ?)"

	vkAccessID, err := v.InsertInto(insertIntoVkAccess, nil, tx)
	if err != nil {
		return Source{}, VKAccess{}, err
	}

	v.ID = vkAccessID

	err = tx.Commit()
	if err != nil {
		return Source{}, VKAccess{}, err
	}

	return s, v, nil
}
