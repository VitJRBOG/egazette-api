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

// InsertArticle inserts a new record into the "article" table.
func InsertArticle(dbConn Connection, sourceName string, article models.Article) error {
	query := "INSERT INTO article (url, date, title, description, cover_url, source_id) " +
		"SELECT $1, $2, $3, $4, $5, (SELECT id FROM source WHERE name=$6) " +
		"WHERE NOT EXISTS (SELECT id FROM article WHERE url=$1)"

	// FIXME: need to describe an inserting for the multiple records.

	_, err := dbConn.Conn.Exec(query,
		article.URL, article.Date, article.Title, article.Description, article.CoverURL, sourceName)
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

		err := rows.Scan(&id, &article.URL, &article.Date, &article.Title,
			&article.Description, &article.CoverURL, &sourceID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows from the 'article' table: %s", err)
		}

		articles = append(articles, article)
	}

	return articles, nil
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

		err := rows.Scan(&id, &source.Name, &source.Description, &source.HomeURL, &source.APIName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows from the 'source' table: %s", err)
		}

		sources = append(sources, source)
	}

	return sources, nil
}
