package newreleases

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/ynori7/music/allmusic"
	"github.com/ynori7/music/config"
	"github.com/ynori7/music/filter"
	"github.com/ynori7/music/view"
)

type newReleasesHandler struct {
	config config.Config
}

func NewReleasesHandler(
	conf config.Config,
) newReleasesHandler {
	return newReleasesHandler{
		config: conf,
	}
}

func (h newReleasesHandler) GenerateNewReleasesReport() {
	logger := log.WithFields(log.Fields{"Logger": "GenerateNewReleasesReport"})

	//Fetch the new releases (filtered by top-level genre)
	newReleases, err := allmusic.GetPotentiallyInterestingNewReleases(h.config)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Error fetching new releases")
	}

	//Fetch the discographies and filter the releases
	filterer := filter.NewFilterer(h.config, newReleases)
	interestingDiscographies := filterer.FilterAndEnrich()

	//Build HTML output
	template := view.NewHtmlTemplate(interestingDiscographies)
	out, err := template.ExecuteHtmlTemplate()
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Error generating html")
	}

	//Save HTML output to file
	err = ioutil.WriteFile("27-03-2020.html", []byte(out), 0644)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Error saving html to file")
	}
}
