package filter

import (
	"errors"
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/ynori7/music/allmusic"
	"github.com/ynori7/music/config"
	"github.com/ynori7/workerpool"
)

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

	//Process results
	discographies := make([]allmusic.Discography, 0)

	//Set up worker pool
	workerPool := workerpool.NewWorkerPool(5,
		func(result interface{}) {
			r := result.(allmusic.Discography)
			logger.WithFields(log.Fields{"Artist": r.Artist.Name}).Debug("Found interesting artist")
			discographies = append(discographies, r)
		},
		func(err error) {
			unwrappedErr := errors.Unwrap(err)
			switch unwrappedErr {
			case ErrAlbumNotFound, ErrNotHighEnoughRatings, ErrNotInterestingGenre:
				logger.WithFields(log.Fields{"error": err}).Info("Filtered an artist")
			default:
				logger.WithFields(log.Fields{"error": err}).Error("Error looking up artist data")
			}
		},
		f.processNewRelease,
	)

	//Do the work
	if err := workerPool.Work(f.potentialReleases); err != nil {
		logger.WithFields(log.Fields{"error": err}).Error("Error processing jobs")
	}

	//Sort the results
	sort.Slice(discographies, func(i, j int) bool {
		return discographies[i].Score > discographies[j].Score
	})

	return discographies
}

func (f Filterer) processNewRelease(job interface{}) (result interface{}, err error) {
	j := job.(allmusic.NewRelease)

	discography, err := f.discographyClient.GetArtistDiscography(j.ArtistLink)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, j.ArtistLink)
	}

	//validate genres
	if !f.artistHasInterestingGenre(discography.Artist.Genres) {
		return nil, fmt.Errorf("%w: %s", ErrNotInterestingGenre, discography.Artist.Name)
	}

	//validate ratings
	if discography.BestRating < 8 {
		return nil, fmt.Errorf("%w: %s", ErrNotHighEnoughRatings, discography.Artist.Name)
	}

	//filtering out singles and EPs
	newestRelease := f.findNewRelease(discography, j.NewAlbumTitle)
	if newestRelease == nil {
		return nil, fmt.Errorf("%w: %s - %s", ErrAlbumNotFound, discography.Artist.Name, j.NewAlbumTitle) //it was probably a single or an EP
	}

	//push the new release
	discography.NewestRelease = *newestRelease
	return *discography, nil
}

func (f Filterer) findNewRelease(discography *allmusic.Discography, releaseTitle string) *allmusic.Album {
	for _, album := range discography.Albums {
		if album.Title == releaseTitle {
			return &album //ensure we've selected the right one, just to be safe. Sometimes they aren't sorted properly
		}
	}

	return nil
}

func (f Filterer) artistHasInterestingGenre(genres []string) bool {
	for _, g := range genres {
		if f.conf.IsInterestingSubGenre(g) {
			return true
		}
	}
	return false
}
