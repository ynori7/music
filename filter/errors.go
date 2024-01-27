package filter

import "fmt"

var (
	ErrNotInterestingGenre  = fmt.Errorf("artist is not an interesting genre")
	ErrNotHighEnoughRatings = fmt.Errorf("artist doesn't have high enough ratings")
	ErrAlbumNotFound        = fmt.Errorf("newest album was not found in the list")
)
