package allmusic

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ynori7/hulksmash/anonymizer"
	hulkhttp "github.com/ynori7/hulksmash/http"
	"github.com/ynori7/music/config"
)

const newReleasesUrl = "https://www.allmusic.com/newreleases/all"

// JSON-LD structures for parsing the embedded schema.org data
type jsonLDData struct {
	Graph []musicAlbum `json:"@graph"`
}

type musicAlbum struct {
	Type     string   `json:"@type"`
	Name     string   `json:"name"`
	ByArtist []artist `json:"byArtist"`
	Genre    string   `json:"genre"`
	URL      string   `json:"url"`
}

type artist struct {
	Type string `json:"@type"`
	Name string `json:"name"`
}

type ReleasesClient struct {
	httpClient    *http.Client
	conf          config.Config
	reqAnonymizer anonymizer.Anonymizer
}

func NewReleasesClient(conf config.Config) ReleasesClient {
	return ReleasesClient{
		httpClient:    hulkhttp.NewClient(),
		conf:          conf,
		reqAnonymizer: anonymizer.New(int64(rand.Int())),
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
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rc.reqAnonymizer.AnonymizeRequest(req)

	res, err := rc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	b, _ := io.ReadAll(res.Body)
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}

	// Extract JSON-LD data from the script tag
	var jsonData jsonLDData
	doc.Find("script[type='application/ld+json']").Each(func(i int, s *goquery.Selection) {
		jsonText := s.Text()
		if err := json.Unmarshal([]byte(jsonText), &jsonData); err == nil {
			// Successfully parsed JSON-LD data
			return
		}
	})

	// Build a map of album URL -> album data from JSON
	albumDataMap := make(map[string]musicAlbum)
	for _, album := range jsonData.Graph {
		if album.Type == "MusicAlbum" {
			albumDataMap[album.URL] = album
		}
	}

	newReleases := make([]NewRelease, 0)

	// Parse the HTML table to get artist links
	doc.Find("#nrTable tr[data-type-filter=\"NEW\"]").Each(func(i int, s *goquery.Selection) {
		// Get album URL from the table
		albumLink := s.Find(".album a")
		albumURL, exists := albumLink.Attr("href")
		if !exists {
			return
		}

		// Look up the album data from JSON
		albumData, found := albumDataMap[albumURL]
		if !found {
			return // No JSON data for this album
		}

		// Filter by genre using JSON data
		if !rc.isInterestingGenre(albumData.Genre) {
			return
		}

		// Filter out compilations
		if isCompilation(albumData.Name) {
			return
		}

		// Get artist link from HTML (skip if not available)
		artistLink := s.Find(".artist a")
		if artistLink.Text() == "" {
			return // No artist link available, skip this release
		}
		bandLink, _ := artistLink.Attr("href")

		newReleases = append(newReleases, NewRelease{
			ArtistLink:    bandLink,
			NewAlbumTitle: albumData.Name,
		})
	})

	return newReleases, nil
}

func (rc ReleasesClient) isInterestingGenre(genre string) bool {
	// Genre can be comma-separated, check each one
	genres := strings.Split(genre, ",")
	for _, g := range genres {
		g = strings.TrimSpace(g)
		if rc.conf.IsInterestingMainGenre(g) {
			return true
		}
	}
	return false
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
