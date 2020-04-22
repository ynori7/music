package rateyourmusic

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
	"strings"
)

const BaseUrl = "https://rateyourmusic.com/artist/"

type DiscographyClient struct {
	httpClient *http.Client
}

func NewDiscographyClient() DiscographyClient {
	return DiscographyClient{
		httpClient: &http.Client{},
	}
}

func BuildRymLink(artistName string) string {
	return BaseUrl + strings.Replace(strings.ToLower(artistName), " ", "-", -1)
}

func (dc DiscographyClient) GetArtistDiscography(link string) (*Discography, error) {
	// Request the HTML page.
	res, err := dc.httpClient.Get(link)
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

	discography := new(Discography)
	discography.Url = link
	discography.Artist = strings.TrimSpace(doc.Find(".artist_name_hdr").First().Text())

	// Scan the discography
	ratingCountTotal := 0
	var ratingSum float64

	doc.Find("#disco_type_s .disco_release").Each(func(i int, s *goquery.Selection) {
		album := Album{}

		album.Title = strings.TrimSpace(s.Find(".disco_mainline a.album").Text())
		if album.Title == "" {
			return //this was probably the sort row at the top
		}

		album.Year = strings.TrimSpace(s.Find(".disco_year_ymd").Text())
		rating := strings.TrimSpace(s.Find(".disco_avg_rating").Text())
		album.AverageRating, _ = strconv.ParseFloat(rating, 64)
		ratingCount := strings.Replace(strings.TrimSpace(s.Find(".disco_ratings").Text()), ",", "", -1)
		album.RatingCount, _ = strconv.Atoi(ratingCount)

		ratingCountTotal += album.RatingCount
		ratingSum += album.AverageRating * float64(album.RatingCount)

		discography.Albums = append(discography.Albums, album)
	})

	discography.RatingCount = ratingCountTotal
	discography.AverageRating = ratingSum / float64(ratingCountTotal)

	return discography, nil
}