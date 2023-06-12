package harvester

import (
	"egazette-api/internal/db"
	"egazette-api/internal/loggers"
	"egazette-api/internal/models"
	"egazette-api/internal/sources/jpl"
	"egazette-api/internal/sources/natgeo"
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
		cleaning(dbConn, sources)

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
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("JPL harvester recovered from the panic: %s", r)
			}
		}()

		err := harvestTheJPLArticles(dbConn, sources)
		if err != nil {
			log.Printf("failed to harvest an articles from the JPL website: %s", err)
		}
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Vestirama harvester recovered from the panic: %s", r)
			}
		}()

		err := harvestTheVestiramaArticles(dbConn, sources)
		if err != nil {
			log.Printf("failed to harvest an articles from the Vestirama website: %s", err)
		}
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("NatGeo harvester recovered from the panic: %s", r)
			}
		}()

		err := harvestTheNatGeoArticles(dbConn, sources)
		if err != nil {
			log.Printf("failed to harvest an articles from the NatGeo website: %s", err)
		}
	}(&wg)

	wg.Wait()
}

func harvestTheJPLArticles(dbConn db.Connection, sources []models.Source) error {
	source := models.FindSourceByAPIName(sources, "jpl")

	articles, err := jpl.GetArticleData()
	if err != nil {
		return err
	}

	err = db.InsertArticles(dbConn, source.APIName, articles)
	if err != nil {
		return err
	}

	return nil
}

func harvestTheVestiramaArticles(dbConn db.Connection, sources []models.Source) error {
	source := models.FindSourceByAPIName(sources, "vestirama")

	articles, err := vestirama.GetArticleData()
	if err != nil {
		return err
	}

	err = db.InsertArticles(dbConn, source.APIName, articles)
	if err != nil {
		return err
	}

	return nil
}

func harvestTheNatGeoArticles(dbConn db.Connection, sources []models.Source) error {
	source := models.FindSourceByAPIName(sources, "natgeo")

	articles, err := natgeo.GetArticleData()
	if err != nil {
		return err
	}

	err = db.InsertArticles(dbConn, source.APIName, articles)
	if err != nil {
		return err
	}

	return nil
}

func cleaning(dbConn db.Connection, sources []models.Source) {
	err := clearTheJPLArticlesList(dbConn, sources)
	if err != nil {
		log.Printf("failed to clear JPL articles: %s", err)
	}

	err = clearTheVestiramaArticlesList(dbConn, sources)
	if err != nil {
		log.Printf("failed to clear Vestirama articles: %s", err)
	}

	err = clearTheNatGeoArticlesList(dbConn, sources)
	if err != nil {
		log.Printf("failed to clear NatGeo articles: %s", err)
	}
}

func clearTheJPLArticlesList(dbConn db.Connection, sources []models.Source) error {
	source := models.FindSourceByAPIName(sources, "jpl")

	err := db.DeleteOldArticles(dbConn, source)
	if err != nil {
		return err
	}

	return nil
}

func clearTheVestiramaArticlesList(dbConn db.Connection, sources []models.Source) error {
	source := models.FindSourceByAPIName(sources, "vestirama")

	err := db.DeleteOldArticles(dbConn, source)
	if err != nil {
		return err
	}

	return nil
}

func clearTheNatGeoArticlesList(dbConn db.Connection, sources []models.Source) error {
	source := models.FindSourceByAPIName(sources, "natgeo")

	err := db.DeleteOldArticles(dbConn, source)
	if err != nil {
		return err
	}

	return nil
}
