package config

import (
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	MainGenres []string `yaml:"main_genres,flow"`
	SubGenres  []string `yaml:"sub_genres,flow"`
}

/**
 * Parse the contents of the YAML file into the Config object.
 */
func (c *Config) Parse(data []byte) error {
	return yaml.Unmarshal(data, &c)
}

func (c *Config) IsInterestingMainGenre(genre string) bool {
	return stringContainsListItem(genre, c.MainGenres)
}

func (c *Config) IsInterestingSubGenre(genre string) bool {
	return stringContainsListItem(genre, c.SubGenres)
}

func stringContainsListItem(str string, list []string) bool {
	for _, s := range list {
		if strings.Contains(str, s) {
			return true
		}
	}
	return false
}
