package config

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	testConfig := []byte(`title: "rap-and-metal"
main_genres:
  - "Rock"
  - "Rap"
  - "R&B" #sometimes rap albums show up as this
sub_genres:
  fuzzy_matches:
    - "Metal"
    - "Grindcore"
    - "Hardcore"
    - "Rap"
    - "Virtuoso"
  exact_matches:
    - "Grunge" #because we don't want to match "Post-Grunge"`)

	c := Config{}

	err := c.Parse(testConfig)
	require.NoError(t, err, "It should parse the config successfully")

	assert.Equal(t, "rap-and-metal", c.Title)
	assert.Equal(t, 3, len(c.MainGenres))
	assert.Equal(t, 5, len(c.SubGenres.FuzzyMatches))
	assert.Equal(t, 1, len(c.SubGenres.ExactMatches))
}

func Test_IsInterestingMainGenre(t *testing.T) {
	testcases := map[string]struct {
		List     []string
		Genre    string
		Expected bool
	}{
		"No match": {
			List:     []string{"Jazz", "Country", "Pop/Rock"},
			Genre:    "Rap",
			Expected: false,
		},
		"Exact match": {
			List:     []string{"Jazz", "Rap", "Rock"},
			Genre:    "Rap",
			Expected: true,
		},
		"Fuzzy match": {
			List:     []string{"Jazz", "Rap", "Rock"},
			Genre:    "Pop/Rock",
			Expected: true,
		},
		"Empty list": {
			List:     []string{},
			Genre:    "Rock",
			Expected: false,
		},
	}

	for testcase, testdata := range testcases {
		c := Config{MainGenres: testdata.List}
		res := c.IsInterestingMainGenre(testdata.Genre)
		assert.Equal(t, testdata.Expected, res, testcase)
	}
}

func Test_IsInterestingSubGenre(t *testing.T) {
	testcases := map[string]struct {
		FuzzyList []string
		ExactList []string
		Genre     string
		Expected  bool
	}{
		"No match": {
			FuzzyList: []string{"Jazz", "Country", "Pop/Rock"},
			ExactList: []string{"Grunge"},
			Genre:     "Rap",
			Expected:  false,
		},
		"Exact match in fuzzy list": {
			FuzzyList: []string{"Jazz", "Rap", "Rock"},
			ExactList: []string{"Grunge"},
			Genre:     "Rap",
			Expected:  true,
		},
		"Exact match in exact list": {
			FuzzyList: []string{"Jazz", "Rap", "Rock"},
			ExactList: []string{"Grunge"},
			Genre:     "Grunge",
			Expected:  true,
		},
		"Fuzzy match": {
			FuzzyList: []string{"Jazz", "Rap", "Rock"},
			ExactList: []string{"Grunge"},
			Genre:     "Pop/Rock",
			Expected:  true,
		},
		"Fuzzy match of exact list": {
			FuzzyList: []string{"Jazz", "Rap", "Rock"},
			ExactList: []string{"Grunge"},
			Genre:     "Post-Grunge",
			Expected:  false,
		},
		"Empty list": {
			FuzzyList: []string{},
			ExactList: []string{},
			Genre:     "Rock",
			Expected:  false,
		},
	}

	for testcase, testdata := range testcases {
		c := Config{}
		c.SubGenres.FuzzyMatches = testdata.FuzzyList
		c.SubGenres.ExactMatches = testdata.ExactList

		res := c.IsInterestingSubGenre(testdata.Genre)
		assert.Equal(t, testdata.Expected, res, testcase)
	}
}
