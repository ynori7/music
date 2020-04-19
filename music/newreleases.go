package music

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ynori7/MusicNewReleases/config"
)

func GetPotentiallyInterestingNewReleases(conf config.Config) (NewReleases, error) {
	// Request the HTML page.
	res, err := http.Get("https://www.allmusic.com/newreleases/all/20200410")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	newReleaseSet := make(map[string]struct{}, 0) //used for deduplication
	newReleases := make(NewReleases, 0)

	// Find the new releases
	doc.Find(".all-new-releases tr[data-type-filter=\"NEW\"]").Each(func(i int, s *goquery.Selection) {
		genre := s.Find(".genre a").Text()
		if !conf.IsInterestingMainGenre(genre) {
			return
		}

		//Try to filter out live albums and compilations
		album := s.Find(".album a").Text()
		if isCompilation(album) {
			return
		}

		// For each item found, get the band and title
		band := s.Find(".artist a")
		if band.Text() == "" {
			return //sometimes there is no artist page
		}
		bandLink, _ := band.Attr("href")
		bandLink = bandLink +"/discography"

		if _, ok := newReleaseSet[bandLink]; !ok {
			newReleases = append(newReleases, bandLink+"/discography")
			newReleaseSet[bandLink] = struct{}{}
		}
	})

	return newReleases, nil
}

var compilationIndicators = []string{"Live", "Compilation", "Best of", "Interview", "From the Vault"}
func isCompilation(title string) bool {
	for _, i := range compilationIndicators {
		if strings.Contains(title, i) {
			return true
		}
	}
	return false
}