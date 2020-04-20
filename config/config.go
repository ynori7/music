package config

import (
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Title      string
	MainGenres []string  `yaml:"main_genres,flow"`
	SubGenres  SubGenres `yaml:"sub_genres"`
	Email      Email
}

type SubGenres struct {
	FuzzyMatches []string `yaml:"fuzzy_matches,flow"`
	ExactMatches []string `yaml:"exact_matches,flow"`
}

type Email struct {
	Enabled    bool
	PrivateKey string `yaml:"private_key"`
	PublicKey  string `yaml:"public_key"`
	From       EmailRecipient
	To         EmailRecipient
}

type EmailRecipient struct {
	Address string
	Name    string
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
	return stringContainsListItem(genre, c.SubGenres.FuzzyMatches) || isContainedInList(genre, c.SubGenres.ExactMatches)
}

func stringContainsListItem(str string, list []string) bool {
	for _, s := range list {
		if strings.Contains(str, s) {
			return true
		}
	}
	return false
}

func isContainedInList(str string, list []string) bool {
	for _, s := range list {
		if str == s {
			return true
		}
	}
	return false
}
