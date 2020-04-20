package filter

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/ynori7/MusicNewReleases/config"
	"github.com/ynori7/MusicNewReleases/music"
	"sort"
)

const WorkerCount = 4

type Filterer struct {
	conf              config.Config
	potentialReleases []music.NewRelease
}

func NewFilterer(conf config.Config, releases []music.NewRelease) Filterer {
	return Filterer{
		conf:              conf,
		potentialReleases: releases,
	}
}

func (f Filterer) FilterAndEnrich() []music.Discography {
	logger := log.WithFields(log.Fields{"Logger": "FilterAndEnrich"})

	resultsChan := make(chan music.Discography, WorkerCount)
	errorChan := make(chan error, WorkerCount)

	//Spawn workers to process in parallel
	workers := make([]chan music.NewRelease, WorkerCount)
	for i := 0; i < WorkerCount; i++ {
		workers[i] = make(chan music.NewRelease, len(f.potentialReleases)/WorkerCount)
		go f.enrichAndFilterWorker(resultsChan, errorChan, workers[i])
	}

	//Assign an equal number of releases to be checked by each worker
	var i = 0
	for _, s := range f.potentialReleases {
		workers[i] <- s
		i = (i + 1) % WorkerCount
	}

	//Process results
	discographies := make([]music.Discography, 0)
	for i := 0; i < len(f.potentialReleases); i++ {
		select {
		case r := <-resultsChan:
			logger.WithFields(log.Fields{"Artist": r.Artist.Name}).Debug("Found interesting artist")
			discographies = append(discographies, r)
		case err := <-errorChan:
			logger.WithFields(log.Fields{"error": err}).Info("Filtered an artist")
		}
	}

	//Sort the results
	sort.Slice(discographies, func(i, j int) bool {
		return discographies[i].Score > discographies[j].Score
	})

	//Signal workers to stop working
	for _, worker := range workers {
		close(worker)
	}

	return discographies
}

func (f Filterer) enrichAndFilterWorker(successes chan music.Discography, errors chan error, jobs chan music.NewRelease) {
	for j := range jobs {
		discography, err := music.GetArtistDiscography(j.ArtistLink)
		if err != nil {
			errors <- err
			continue
		}

		//validate genres
		if !f.artistHasInterestingGenre(discography.Artist.Genres) {
			errors <- fmt.Errorf("artist is not an interesting genre")
			continue
		}

		//validate ratings
		if discography.BestRating < 8 {
			errors <- fmt.Errorf("artist doesn't have high enough ratings")
			continue
		}

		//filtering out singles and EPs
		foundNewAlbum := false
		for _, album := range discography.Albums {
			if album.Title == j.NewAlbumTitle {
				foundNewAlbum = true
				discography.NewestRelease = album //ensure we've selected the right one, just to be safe. Sometimes they aren't sorted properly
				break
			}
		}
		if !foundNewAlbum {
			errors <- fmt.Errorf("newest album was not found in the list") //it was probably a single or an EP
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
