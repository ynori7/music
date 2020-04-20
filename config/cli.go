package config

import "flag"

var CliConf CliConfig

type CliConfig struct {
	ConfigFile     string
	NewReleaseWeek string //optional
	OutputPath     string //optional
}

func ParseCliFlags() {
	configFile := flag.String("config", "", "the path to the configuration yaml")
	newReleaseWeek := flag.String("new-release-week", "", "the new release week (always a Thursday) in the format YYYYMMDD")
	output := flag.String("output", "out", "the path where output files should be saved")

	flag.Parse()

	CliConf.ConfigFile = *configFile
	CliConf.NewReleaseWeek = *newReleaseWeek
	CliConf.OutputPath = *output
}