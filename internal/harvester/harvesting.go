package harvester

import (
	"egazette-api/internal/db"
	"egazette-api/internal/loggers"
	"egazette-api/internal/sources/jpl"
	"egazette-api/internal/sources/vestirama"
	"log"
	"time"
)

// Harvesting starts the articles gatherers.
func Harvesting(dbConn db.Connection) {
	infoLogger := loggers.NewInfoLogger()

	infoLogger.Println("harvesting of articles is started")

	for {
		go harvest(dbConn)

		// FIXME: a randomiser for sleeping time should be described
		// to send requests in a less conspicuous way.

		time.Sleep(60 * time.Minute)
	}
}

func harvest(dbConn db.Connection) {
	err := harvestTheJPLArticles(dbConn)
	if err != nil {
		log.Printf("failed to harvest an articles from the JPL website: %s", err)
	}

	err = harvestTheVestiramaArticles(dbConn)
	if err != nil {
		log.Printf("failed to harvest an articles from the Vestirama website: %s", err)
	}
}

func harvestTheJPLArticles(dbConn db.Connection) error {
	articles, err := jpl.GetArticleData()
	if err != nil {
		return err
	}

	sourceName := jpl.GetSourceData().Name
	for _, article := range articles {
		err := db.InsertArticle(dbConn, sourceName, article)
		if err != nil {
			return err
		}
	}

	return nil
}

func harvestTheVestiramaArticles(dbConn db.Connection) error {
	articles, err := vestirama.GetArticleData()
	if err != nil {
		return err
	}

	sourceName := vestirama.GetSourceData().Name
	for _, article := range articles {
		err := db.InsertArticle(dbConn, sourceName, article)
		if err != nil {
			return err
		}
	}

	return nil
}
