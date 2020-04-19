package music

type Artist struct {
	Name   string
	Genres []string
	Link   string
}

type Album struct {
	Title  string
	Link   string
	Rating int //Out of 10. A zero means there is no rating
	Image  string
	Year string
}

type Discography struct {
	Artist        Artist
	Albums        []Album
	AverageRating int
	BestRating    int
}

type NewReleases []string
