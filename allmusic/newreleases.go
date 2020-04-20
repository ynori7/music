package allmusic

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ynori7/music/config"
)

const newReleasesUrl = "https://www.allmusic.com/newreleases/all"

func GetPotentiallyInterestingNewReleases(conf config.Config, week string) ([]NewRelease, error) {
	// Request the HTML page.
	url := newReleasesUrl
	if week != "" {
		url = url + "/" + week
	}
	res, err := http.Get(url)
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

	newReleaseSet := make(map[string]*NewRelease, 0) //used for deduplication

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
		bandLink = bandLink + "/discography"

		if _, ok := newReleaseSet[bandLink]; !ok {
			newReleaseSet[bandLink] = &NewRelease{
				ArtistLink:    bandLink + "/discography",
				NewAlbumTitles: []string{},
			}
		}
		newReleaseSet[bandLink].NewAlbumTitles = append(newReleaseSet[bandLink].NewAlbumTitles, album)
	})

	//turn the map into a list
	newReleases := make([]NewRelease, 0, len(newReleaseSet))
	for _, r := range newReleaseSet {
		newReleases = append(newReleases, *r)
	}

	return newReleases, nil
}

var compilationIndicators = []string{"Live", "Compilation", "Best of", "Interview", "From the Vault", "Collection"}

func isCompilation(title string) bool {
	for _, i := range compilationIndicators {
		if strings.Contains(title, i) {
			return true
		}
	}
	return false
}
