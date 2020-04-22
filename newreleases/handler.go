package newreleases

import (
	"fmt"
	"io/ioutil"
	"time"

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

func (h newReleasesHandler) GenerateNewReleasesReport(week string) (string, error) {
	logger := log.WithFields(log.Fields{"Logger": "GenerateNewReleasesReport"})

	//Fetch the new releases (filtered by top-level genre)
	newReleases, err := allmusic.NewReleasesClient(h.config).GetPotentiallyInterestingNewReleases(allmusic.GetNewReleasesUrlForWeek(week))
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Error("Error fetching new releases")
		return "", err
	}

	//Fetch the discographies and filter the releases
	filterer := filter.NewFilterer(h.config, allmusic.NewDiscographyClient(), newReleases)
	interestingDiscographies := filterer.FilterAndEnrich()

	//Build HTML output
	template := view.NewHtmlTemplate(interestingDiscographies)
	out, err := template.ExecuteHtmlTemplate()
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Error("Error generating html")
		return "", err
	}

	//Save HTML output to file
	dateString := week
	if week == "" {
		dateString = time.Now().Format("20060102") //yyyyMMdd
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s-%s.html", config.CliConf.OutputPath, h.config.Title, dateString), []byte(out), 0644)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Warn("Error saving html to file")
		return "", err
	}

	return out, nil
}
