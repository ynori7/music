package config

import "flag"

var CliConf CliConfig

type CliConfig struct {
	ConfigFile     string
	NewReleaseWeek string //optional
}

func init() {
	configFile := flag.String("config", "", "the path to the configuration yaml")
	newReleaseWeek := flag.String("new-release-week", "", "the new release week (always a Thursday) in the format YYYYMMDD")

	flag.Parse()

	CliConf.ConfigFile = *configFile
	CliConf.NewReleaseWeek = *newReleaseWeek
}