package allmusic

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ynori7/hulksmash/anonymizer"
	hulkhttp "github.com/ynori7/hulksmash/http"
)

const BaseUrl = "https://www.allmusic.com"

type DiscographyClient struct {
	httpClient    *http.Client
	reqAnonymizer anonymizer.Anonymizer
}

func NewDiscographyClient() DiscographyClient {
	return DiscographyClient{
		httpClient:    hulkhttp.NewClient(),
		reqAnonymizer: anonymizer.New(int64(rand.Int())),
	}
}

func (dc DiscographyClient) GetArtistDiscography(link string) (*Discography, error) {
	discography, err := dc.lookupBasicInfo(link)
	if err != nil {
		return nil, err
	}

	// Request the discography
	req, err := http.NewRequest(http.MethodGet, link+"/discographyAjax", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Referer", link)
	dc.reqAnonymizer.AnonymizeRequest(req)
	res, err := dc.httpClient.Do(req)
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

	// Scan the discography
	ratingCount := 0
	ratingSum := 0
	doc.Find("#discography tr").Each(func(i int, s *goquery.Selection) {
		album := Album{}

		album.Title = strings.TrimSpace(s.Find("td.meta a").First().Text())
		if album.Title == "" {
			return //this was probably the sort row at the top
		}

		cover := s.Find("td.cover a")
		if link, ok := cover.Attr("href"); ok {
			album.Link = BaseUrl + strings.TrimSpace(link)
		}
		coverImg := cover.Find("img")
		if img, ok := coverImg.Attr("src"); ok {
			album.Image = strings.TrimSpace(img)
		} else if img, ok := coverImg.Attr("data-src"); ok {
			album.Image = strings.TrimSpace(img)
		}

		album.Year = strings.TrimSpace(s.Find("td.year").Text())

		album.Rating = getEditorRating(s.Find("td.musicRating"))

		if album.Rating > discography.BestRating {
			discography.BestRating = album.Rating
		}
		if album.Rating != 0 {
			ratingCount++
			ratingSum += album.Rating
		}

		discography.Albums = append(discography.Albums, album)
	})

	if len(discography.Albums) == 0 {
		return nil, fmt.Errorf("artist has no albums")
	}

	// Find newest release by iterating the list backwards.
	discography.NewestRelease = discography.Albums[len(discography.Albums)-1] // Sometimes the newest album won't have a year yet, so we use the last item by default
	for i := len(discography.Albums) - 1; i >= 0; i-- {
		if discography.Albums[i].Year == fmt.Sprintf("%d", time.Now().Year()) {
			discography.NewestRelease = discography.Albums[i]
			break
		}
	}

	// Set average rating and calculate score
	if ratingCount > 0 {
		discography.AverageRating = getAverage(ratingSum, ratingCount)
	}
	discography.Score = calculateScore(discography.BestRating, discography.AverageRating, ratingCount, discography.NewestRelease.Rating)

	return discography, nil
}

func (dc DiscographyClient) lookupBasicInfo(link string) (*Discography, error) {
	// Request the HTML page.
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, err
	}
	dc.reqAnonymizer.AnonymizeRequest(req)
	res, err := dc.httpClient.Do(req)
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

	discography.Artist.Link = link
	artistNameNode := doc.Find("#artistName").First()
	artistNameNode.Find("span").Remove() // Remove child spans (e.g., follower count)
	discography.Artist.Name = strings.TrimSpace(artistNameNode.Text())

	doc.Find("#basicInfoMeta .styles a").Each(func(i int, s *goquery.Selection) {
		if genre := strings.TrimSpace(s.Text()); genre != "" {
			discography.Artist.Genres = append(discography.Artist.Genres, strings.TrimSpace(s.Text()))
		}
	})

	return discography, nil
}

func getEditorRating(s *goquery.Selection) int {
	ratingVal, _ := s.Attr("data-text")
	ratingInt := 0
	if len(ratingVal) > 0 {
		ratingInt, _ = strconv.Atoi(string(ratingVal[0]))
		if ratingInt > 0 {
			ratingInt = ratingInt + 1
		}
	}
	return ratingInt
}

func calculateScore(bestRating, averageRating, ratingCount, newestAlbumRating int) int {
	score := bestRating * 4    //40% of the score is from the best rating (gives some extra weight to the average)
	score += averageRating * 4 //another 40% is from the average rating

	if newestAlbumRating >= 8 {
		score = newestAlbumRating * 8 //if the newest album is a well-rated one, use that as the base for the score to give it extra boost
	}

	//add a little extra weight based on the total number of ratings there were
	switch {
	case ratingCount > 7:
		score += 20
	case ratingCount > 5:
		score += 15
	case ratingCount > 2:
		score += 10
	case ratingCount > 1:
		score += 5
	}

	return score
}

func getAverage(sum, count int) int {
	return int(math.RoundToEven(float64(sum) / float64(count)))
}
