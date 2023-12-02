package allmusic

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ynori7/music/config"
)

const newReleasesUrl = "https://www.allmusic.com/newreleases/all"

type ReleasesClient struct {
	httpClient *http.Client
	conf config.Config
}

func NewReleasesClient(conf config.Config) ReleasesClient {
	return ReleasesClient{
		httpClient: &http.Client{},
		conf: conf,
	}
}

func GetNewReleasesUrlForWeek(week string) string {
	url := newReleasesUrl
	if week != "" {
		url = url + "/" + week
	}
	return url
}

func (rc ReleasesClient) GetPotentiallyInterestingNewReleases(url string) ([]NewRelease, error) {
	// Request the HTML page.
	res, err := rc.httpClient.Get(url)
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

	newReleases := make([]NewRelease, 0)

	// Find the new releases
	doc.Find("#nrTable tr[data-type-filter=\"NEW\"]").Each(func(i int, s *goquery.Selection) {
		genre := s.Find(".genre a").Text()
		if !rc.conf.IsInterestingMainGenre(genre) {
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

		newReleases = append(newReleases, NewRelease{
			ArtistLink:    bandLink,
			NewAlbumTitle: album,
		})
	})

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
