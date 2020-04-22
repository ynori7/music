package filter

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ynori7/music/allmusic"
	"github.com/ynori7/music/config"
)

const WorkerCount = 5

type Filterer struct {
	conf              config.Config
	potentialReleases []allmusic.NewRelease
	discographyClient allmusic.DiscographyClient
}

func NewFilterer(conf config.Config, discographyClient allmusic.DiscographyClient, releases []allmusic.NewRelease) Filterer {
	return Filterer{
		conf:              conf,
		potentialReleases: releases,
		discographyClient: discographyClient,
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
			unwrappedErr := errors.Unwrap(err)
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
		discography, err := f.discographyClient.GetArtistDiscography(j.ArtistLink)
		if err != nil {
			errors <- fmt.Errorf("%w: %s", err, j.ArtistLink)
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
		newestReleases := f.findNewReleases(discography, j.NewAlbumTitles)
		if len(newestReleases) == 0 {
			errors <- fmt.Errorf("%w: %s - %s", ErrAlbumNotFound, discography.Artist.Name, strings.Join(j.NewAlbumTitles, ",")) //it was probably a single or an EP
			continue
		}

		//push all the new releases (an artist can theoretically have more than one)
		for _, r := range newestReleases {
			discography.NewestRelease = r
			successes <- *discography
		}
	}
}

func (f Filterer) findNewReleases(discography *allmusic.Discography, releaseTitles []string) []allmusic.Album {
	newestReleases := make([]allmusic.Album, 0, len(releaseTitles))
	for _, album := range discography.Albums {
		if inArray(album.Title, releaseTitles) {
			newestReleases = append(newestReleases, album) //ensure we've selected the right one, just to be safe. Sometimes they aren't sorted properly
		}
	}

	return newestReleases
}

func (f Filterer) artistHasInterestingGenre(genres []string) bool {
	for _, g := range genres {
		if f.conf.IsInterestingSubGenre(g) {
			return true
		}
	}
	return false
}

func inArray(needle string, haystack []string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
