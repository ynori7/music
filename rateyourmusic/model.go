package rateyourmusic

type Discography struct {
	Artist        string
	Url           string
	Albums        []Album
	AverageRating float64
	RatingCount   int
}

type Album struct {
	Title         string
	AverageRating float64
	RatingCount   int
	Year          string
}
