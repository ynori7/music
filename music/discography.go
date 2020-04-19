package music

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const BaseUrl = "https://www.allmusic.com"

func GetArtistDiscography(link string) (*Discography, error) {
	// Request the HTML page.
	res, err := http.Get(link)
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
	discography.Artist.Name = strings.TrimSpace(doc.Find(".artist-name").First().Text())

	doc.Find(".basic-info .styles a").Each(func(i int, s *goquery.Selection) {
		if genre := strings.TrimSpace(s.Text()); genre != "" {
			discography.Artist.Genres = append(discography.Artist.Genres, strings.TrimSpace(s.Text()))
		}
	})

	ratingCount := 0
	ratingSum := 0
	doc.Find(".discography tr").Each(func(i int, s *goquery.Selection) {
		album := Album{}

		album.Title = strings.TrimSpace(s.Find("td.title a").Text())
		if album.Title == "" {
			return //this was probably the sort row at the top
		}

		cover := s.Find("td.cover a")
		if link, ok := cover.Attr("href"); ok {
			album.Link = BaseUrl + strings.TrimSpace(link)
		}
		coverImg := cover.Find("img")
		if img, ok := coverImg.Attr("data-original"); ok {
			album.Image = strings.TrimSpace(img)
		}

		album.Year = strings.TrimSpace(s.Find("td.year").Text())

		album.Rating = GetEditorRating(s.Find("td.all-rating"))

		if album.Rating > discography.BestRating {
			discography.BestRating = album.Rating
		}
		if album.Rating != 0 {
			ratingCount++
			ratingSum += album.Rating
		}

		discography.Albums = append(discography.Albums, album)
	})

	if ratingCount > 0 {
		discography.AverageRating = int(math.RoundToEven(float64(ratingSum) / float64(ratingCount)))
	}

	return discography, nil
}

func GetEditorRating(s *goquery.Selection) int {
	ratingVal, _ := s.Attr("data-sort-value")
	ratingInt := 0
	if len(ratingVal) > 0 {
		ratingInt, _ = strconv.Atoi(string(ratingVal[0]))
		if ratingInt > 0 {
			ratingInt = ratingInt + 1
		}
	}
	return ratingInt
}
