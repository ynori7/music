package main

import (
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/ynori7/MusicNewReleases/config"
	"github.com/ynori7/MusicNewReleases/filter"
	"github.com/ynori7/MusicNewReleases/music"
	"github.com/ynori7/MusicNewReleases/view"
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

	//Fetch the new releases (filtered by top-level genre)
	newReleases, err := music.GetPotentiallyInterestingNewReleases(conf)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Error fetching new releases")
	}

	//Fetch the discographies and filter the releases
	filterer := filter.NewFilterer(conf, newReleases)
	interestingDiscographies := filterer.FilterAndEnrich()

	//Build HTML output
	template := view.NewHtmlTemplate(interestingDiscographies)
	out, err := template.ExecuteHtmlTemplate()
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Error generating html")
	}

	//Save HTML output to file
	err = ioutil.WriteFile("10-04-2020.html", []byte(out), 0644)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Error saving html to file")
	}

}
