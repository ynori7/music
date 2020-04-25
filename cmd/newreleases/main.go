package main

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/ynori7/music/config"
	"github.com/ynori7/music/email"
	"github.com/ynori7/music/newreleases"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)
	logger := log.WithFields(log.Fields{"Logger": "main"})

	//Get the cli flags
	config.ParseCliFlags()
	if config.CliConf.ConfigFile == "" {
		logger.Fatal("You must specify the path to the config file")
	}

	//Get the config
	data, err := ioutil.ReadFile(config.CliConf.ConfigFile)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Error reading config file")
	}

	var conf config.Config
	if err := conf.Parse(data); err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Error parsing config")
	}

	//Generate the report
	newReleasesHandler := newreleases.NewReleasesHandler(conf)
	report, err := newReleasesHandler.GenerateNewReleasesReport(config.CliConf.NewReleaseWeek)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Unable to generate report")
	}

	if conf.Email.Enabled {
		mailer := email.NewMailer(conf)
		if err := mailer.SendMail(email.GetNewReleasesSubjectLine(config.CliConf.NewReleaseWeek), report); err != nil {
			logger.WithFields(log.Fields{"error": err}).Error("Error sending email")
		}
	}
}
