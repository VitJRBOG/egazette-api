package harvester

import (
	"egazette-api/internal/db"
	"egazette-api/internal/loggers"
	"egazette-api/internal/models"
	"egazette-api/internal/sources/jpl"
	"egazette-api/internal/sources/vestirama"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

// Harvesting starts the articles gatherers.
func Harvesting(wg *sync.WaitGroup, signalToExit chan os.Signal, dbConn db.Connection,
	sources []models.Source) {
	infoLogger := loggers.NewInfoLogger()

	infoLogger.Println("harvesting of articles is started")

harvy:
	for {
		harvest(dbConn, sources)

		n := rand.Intn(3600)
		waitFor := 3600 + n

		for i := 0; i < waitFor; i++ {
			select {
			case <-signalToExit:
				i = waitFor
				break harvy
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}

	loggers.NewInfoLogger().Println("harvester exited successfully")
	wg.Done()
}

func harvest(dbConn db.Connection, sources []models.Source) {
	err := harvestTheJPLArticles(dbConn, sources)
	if err != nil {
		log.Printf("failed to harvest an articles from the JPL website: %s", err)
	}

	err = harvestTheVestiramaArticles(dbConn, sources)
	if err != nil {
		log.Printf("failed to harvest an articles from the Vestirama website: %s", err)
	}
}

func harvestTheJPLArticles(dbConn db.Connection, sources []models.Source) error {
	source := models.FindSourceByAPIName(sources, "jpl")

	articles, err := jpl.GetArticleData()
	if err != nil {
		return err
	}

	for _, article := range articles {
		err := db.InsertArticle(dbConn, source.Name, article)
		if err != nil {
			return err
		}
	}

	return nil
}

func harvestTheVestiramaArticles(dbConn db.Connection, sources []models.Source) error {
	source := models.FindSourceByAPIName(sources, "vestirama")

	articles, err := vestirama.GetArticleData()
	if err != nil {
		return err
	}

	for _, article := range articles {
		err := db.InsertArticle(dbConn, source.Name, article)
		if err != nil {
			return err
		}
	}

	return nil
}
