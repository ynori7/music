package filter

import (
	"fmt"
	"time"

	"github.com/ynori7/MusicNewReleases/config"
	"github.com/ynori7/MusicNewReleases/music"
)

const WorkerCount = 3

type Filterer struct {
	conf              config.Config
	potentialReleases music.NewReleases
}

func NewFilterer(conf config.Config, releases music.NewReleases) Filterer {
	return Filterer{
		conf:              conf,
		potentialReleases: releases,
	}
}

func (f Filterer) FilterAndEnrich() []music.Discography {
	resultsChan := make(chan music.Discography, WorkerCount)
	errorChan := make(chan error, WorkerCount)

	workers := make([]chan string, WorkerCount)
	for i := 0; i < WorkerCount; i++ {
		workers[i] = make(chan string, len(f.potentialReleases)/WorkerCount)
		go f.enrichAndFilterWorker(resultsChan, errorChan, workers[i])
	}

	// Each worker will process an equal number of releases
	var i = 0
	for _, s := range f.potentialReleases {
		workers[i] <- s
		i = (i + 1) % WorkerCount
	}

	discogs := make([]music.Discography, 0)

	for i := 0; i < len(f.potentialReleases); i++ {
		select {
		case r := <-resultsChan:
			fmt.Println("Found interesting artist: ", r.Artist.Name)
			discogs = append(discogs, r)
		case err := <-errorChan:
			fmt.Println(err)
		}
	}

	for _, worker := range workers {
		close(worker)
	}

	return discogs
}

func (f Filterer) enrichAndFilterWorker(successes chan music.Discography, errors chan error, jobs chan string) {
	for j := range jobs {
		discography, err := music.GetArtistDiscography(j)
		if err != nil {
			errors <- err
			continue
		}

		if !f.artistHasInterestingGenre(discography.Artist.Genres) {
			errors <- fmt.Errorf("artist is not an interesting genre")
			continue
		}

		if discography.BestRating < 8 {
			errors <- fmt.Errorf("artist doesn't have high enough ratings")
			continue
		}

		latestAlbum := discography.Albums[len(discography.Albums) - 1]
		if latestAlbum.Year != "" && latestAlbum.Year != fmt.Sprintf("%d", time.Now().Year()) {
			errors <- fmt.Errorf("newest album is not from this year") //this can happen when the newest album was a collaboration
			continue
		}

		successes <- *discography
	}
}

func (f Filterer) artistHasInterestingGenre(genres []string) bool {
	for _, g := range genres {
		if f.conf.IsInterestingSubGenre(g) {
			return true
		}
	}
	return false
}
