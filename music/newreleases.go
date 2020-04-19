package music

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/ynori7/MusicNewReleases/config"
)

func GetPotentiallyInterestingNewReleases(conf config.Config) (NewReleases, error) {
	// Request the HTML page.
	res, err := http.Get("https://www.allmusic.com/newreleases/all")
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

	newReleases := make(NewReleases, 0)

	// Find the new releases
	doc.Find(".all-new-releases tr[data-type-filter=\"NEW\"]").Each(func(i int, s *goquery.Selection) {
		genre := s.Find(".genre a").Text()
		if !conf.IsInterestingMainGenre(genre) {
			return
		}
		// For each item found, get the band and title
		band := s.Find(".artist a")
		if band.Text() == "" {
			return //sometimes there is no artist page
		}
		bandLink, _ := band.Attr("href")

		newReleases = append(newReleases, bandLink+"/discography")
	})

	return newReleases, nil
}
