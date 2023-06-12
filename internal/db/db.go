package db

import (
	"database/sql"
	"egazette-api/internal/models"
	"fmt"

	_ "github.com/lib/pq" // Postgres driver
)

// Connection stores DB connection.
type Connection struct {
	Conn *sql.DB
}

// NewConnection creates a connection to the PostgreSQL database and returns the struct with it.
func NewConnection(dsn string) (Connection, error) {
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		return Connection{}, fmt.Errorf("unable to create a database connection: %s", err.Error())
	}

	err = dbConn.Ping()
	if err != nil {
		return Connection{}, fmt.Errorf("unable connect to database: %s", err.Error())
	}

	return Connection{
		Conn: dbConn,
	}, nil
}

// InsertArticles inserts the new records into the "article" table.
func InsertArticles(dbConn Connection, sourceAPIName string, articles []models.Article) error {
	beginOfQuery := "INSERT INTO article (url, pub_date, title, description, cover_url, add_date, source_id) " +
		"SELECT * FROM ("
	endOfQuery := ") AS newarticle " +
		"WHERE NOT EXISTS (" +
		"SELECT id FROM article WHERE article.url = newarticle.url" +
		")"

	queryWithValues := ""
	values := []interface{}{}
	n := 1

	for i, article := range articles {
		a := fmt.Sprintf(
			"SELECT $%d AS url, $%d AS pub_date, $%d AS title, $%d AS description, $%d AS cover_url, $%d AS add_date, "+
				"(SELECT id FROM source WHERE api_name=$%d) AS source_id", n, n+1, n+2, n+3, n+4, n+5, n+6)
		n += 7
		queryWithValues += a
		if i < len(articles)-1 {
			queryWithValues += " UNION ALL "
		}
		values = append(values, article.URL, article.PubDate, article.Title, article.Description, article.CoverURL, article.AddDate,
			sourceAPIName)
	}

	query := beginOfQuery + queryWithValues + endOfQuery

	stmt, err := dbConn.Conn.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare the query: %s", err)
	}

	_, err = stmt.Exec(values...)
	if err != nil {
		return fmt.Errorf("failed to insert a record into the 'article' table: %s", err)
	}

	return nil
}

// SelectArticles selects records from the "article" table.
func SelectArticles(dbConn Connection, sourceName string) ([]models.Article, error) {
	query := "SELECT * FROM article WHERE source_id=(SELECT id FROM source WHERE name=$1)"

	rows, err := dbConn.Conn.Query(query, sourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to select a records from the 'article' table: %s", err)
	}

	articles := []models.Article{}

	for rows.Next() {
		article := models.Article{}

		// FIXME: has not thought of anything better
		// than writing data into variables that are not used anywhere.

		var id, sourceID int

		err := rows.Scan(&id, &article.URL, &article.PubDate, &article.Title,
			&article.Description, &article.CoverURL, &article.AddDate, &sourceID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows from the 'article' table: %s", err)
		}

		articles = append(articles, article)
	}

	return articles, nil
}

// DeleteOldArticles removes old articles from the "article" table by their "article.source_id"
// if their number is greater than "source.max_articles".
func DeleteOldArticles(dbConn Connection, source models.Source) error {
	query := "DELETE FROM article " +
		"WHERE source_id = (SELECT id FROM source WHERE api_name = $1) " +
		"AND id NOT IN (SELECT id FROM article " +
		"WHERE source_id = (SELECT id FROM source WHERE api_name = $1) " +
		"ORDER BY add_date DESC, pub_date DESC LIMIT $2)"

	_, err := dbConn.Conn.Exec(query, source.APIName, source.MaxArticles)
	if err != nil {
		return fmt.Errorf("failed to delete a records from the 'article' table: %s", err)
	}

	return nil
}

// SelectSources selects records from the "source" table.
func SelectSources(dbConn Connection) ([]models.Source, error) {
	query := "SELECT * FROM source"

	rows, err := dbConn.Conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to select a records from the 'source' table: %s", err)
	}

	sources := []models.Source{}

	for rows.Next() {
		source := models.Source{}

		var id int

		err := rows.Scan(&id, &source.Name, &source.Description, &source.HomeURL,
			&source.APIName, &source.MaxArticles)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows from the 'source' table: %s", err)
		}

		sources = append(sources, source)
	}

	return sources, nil
}
