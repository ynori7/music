package filter

import (
	"errors"
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/ynori7/music/allmusic"
	"github.com/ynori7/music/config"
)

const WorkerCount = 5

type Filterer struct {
	conf              config.Config
	potentialReleases []allmusic.NewRelease
}

func NewFilterer(conf config.Config, releases []allmusic.NewRelease) Filterer {
	return Filterer{
		conf:              conf,
		potentialReleases: releases,
	}
}

func (f Filterer) FilterAndEnrich() []allmusic.Discography {
	logger := log.WithFields(log.Fields{"Logger": "FilterAndEnrich"})

	resultsChan := make(chan allmusic.Discography, WorkerCount)
	errorChan := make(chan error, WorkerCount)

	//Spawn workers to process in parallel
	workers := make([]chan allmusic.NewRelease, WorkerCount)
	for i := 0; i < WorkerCount; i++ {
		workers[i] = make(chan allmusic.NewRelease, len(f.potentialReleases)/WorkerCount)
		go f.enrichAndFilterWorker(resultsChan, errorChan, workers[i])
	}

	//Assign an equal number of releases to be checked by each worker
	var i = 0
	for _, s := range f.potentialReleases {
		workers[i] <- s
		i = (i + 1) % WorkerCount
	}

	//Process results
	discographies := make([]allmusic.Discography, 0)
	for i := 0; i < len(f.potentialReleases); i++ {
		select {
		case r := <-resultsChan:
			logger.WithFields(log.Fields{"Artist": r.Artist.Name}).Debug("Found interesting artist")
			discographies = append(discographies, r)
		case err := <-errorChan:
			unwrappedErr :=  errors.Unwrap(err)
			switch unwrappedErr {
			case ErrAlbumNotFound, ErrNotHighEnoughRatings, ErrNotInterestingGenre:
				logger.WithFields(log.Fields{"error": err}).Info("Filtered an artist")
			default:
				logger.WithFields(log.Fields{"error": err}).Error("Error looking up artist data")
			}
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

func (f Filterer) enrichAndFilterWorker(successes chan allmusic.Discography, errors chan error, jobs chan allmusic.NewRelease) {
	for j := range jobs {
		discography, err := allmusic.GetArtistDiscography(j.ArtistLink)
		if err != nil {
			errors <- fmt.Errorf("%w: %s", err, discography.Artist.Name)
			continue
		}

		//validate genres
		if !f.artistHasInterestingGenre(discography.Artist.Genres) {
			errors <- fmt.Errorf("%w: %s", ErrNotInterestingGenre, discography.Artist.Name)
			continue
		}

		//validate ratings
		if discography.BestRating < 8 {
			errors <- fmt.Errorf("%w: %s", ErrNotHighEnoughRatings, discography.Artist.Name)
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
			errors <- fmt.Errorf("%w: %s - %s", ErrAlbumNotFound, discography.Artist.Name, j.NewAlbumTitle) //it was probably a single or an EP
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
