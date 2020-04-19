package main

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/ynori7/MusicNewReleases/config"
	"github.com/ynori7/MusicNewReleases/filter"
	"github.com/ynori7/MusicNewReleases/music"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)
	logger := log.WithFields(log.Fields{"Logger": "main"})

	if len(os.Args) < 2 {
		logger.Fatal("You must specify the path to the config file")
	}

	//Get the config
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Error reading config file")
	}

	var conf config.Config
	if err := conf.Parse(data); err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Error parsing config")
	}

	newReleases, err := music.GetPotentiallyInterestingNewReleases(conf)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Error fetching new releases")
	}

	filterer := filter.NewFilterer(conf, newReleases)
	interestingDiscographies := filterer.FilterAndEnrich()

	for _, r := range interestingDiscographies {
		fmt.Println(r.Artist.Name)
		fmt.Println(r.Artist.Genres)
		fmt.Printf("Best Rating: %d, Avg Rating: %d\n", r.BestRating, r.AverageRating)
		fmt.Println(r.Albums[len(r.Albums)-1]) //the newest one
	}
}
