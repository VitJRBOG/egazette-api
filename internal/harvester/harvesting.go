package harvester

import (
	"egazette-api/internal/db"
	"egazette-api/internal/loggers"
	"egazette-api/internal/sources/jpl"
	"egazette-api/internal/sources/vestirama"
	"log"
	"os"
	"sync"
	"time"
)

// Harvesting starts the articles gatherers.
func Harvesting(wg *sync.WaitGroup, signalToExit chan os.Signal, dbConn db.Connection) {
	infoLogger := loggers.NewInfoLogger()

	infoLogger.Println("harvesting of articles is started")

harvy:
	for {
		harvest(dbConn)

		// FIXME: a randomiser for sleeping time should be described
		// to send requests in a less conspicuous way.

		for i := 0; i < 360; i++ {
			select {
			case <-signalToExit:
				i = 360
				break harvy
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}

	loggers.NewInfoLogger().Println("harvester exited successfully")
	wg.Done()
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
