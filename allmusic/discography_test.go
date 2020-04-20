package allmusic

import (
	"testing"
)

func Test_calculateScore(t *testing.T) {
	testdata := map[string]struct {
		bestRating        int
		averageRating     int
		ratingCount       int
		newestAlbumRating int
		expectedScore     int
	}{
		//Good top ratings, somewhat low average, high rating count
		"King Diamond with no new album rating": {
			bestRating:        9,
			averageRating:     6,
			ratingCount:       16,
			newestAlbumRating: 0,
			expectedScore:     80,
		},
		"King Diamond with well-rated new album": {
			bestRating:        9,
			averageRating:     6,
			ratingCount:       16,
			newestAlbumRating: 9,
			expectedScore:     92,
		},
		"King Diamond with poorly-rated new album": {
			bestRating:        9,
			averageRating:     6,
			ratingCount:       16,
			newestAlbumRating: 6,
			expectedScore:     80,
		},

		//Okay top ratings, low average, high rating count
		"Manowar with no new album rating": {
			bestRating:        8,
			averageRating:     6,
			ratingCount:       15,
			newestAlbumRating: 0,
			expectedScore:     76,
		},
		"Manowar with poorly-rated new album": {
			bestRating:        8,
			averageRating:     6,
			ratingCount:       15,
			newestAlbumRating: 5,
			expectedScore:     76,
		},
		"Manowar with well-rated new album": {
			bestRating:        8,
			averageRating:     6,
			ratingCount:       15,
			newestAlbumRating: 8,
			expectedScore:     84,
		},

		//Great top ratings, good average, high rating count
		"Eminem with no new album rating": {
			bestRating:        10,
			averageRating:     7,
			ratingCount:       14,
			newestAlbumRating: 0,
			expectedScore:     88,
		},
		"Eminem with poorly-rated new album": {
			bestRating:        10,
			averageRating:     7,
			ratingCount:       14,
			newestAlbumRating: 5,
			expectedScore:     88,
		},
		"Eminem with well-rated new album": {
			bestRating:        10,
			averageRating:     7,
			ratingCount:       14,
			newestAlbumRating: 9,
			expectedScore:     92,
		},

		//Great top ratings, great average, low rating count
		"NWA with no new album rating": {
			bestRating:        10,
			averageRating:     8,
			ratingCount:       3,
			newestAlbumRating: 0,
			expectedScore:     82,
		},
		"NWA with poorly-rated new album": {
			bestRating:        10,
			averageRating:     8,
			ratingCount:       3,
			newestAlbumRating: 6,
			expectedScore:     82,
		},
		"NWA with well-rated new album": {
			bestRating:        10,
			averageRating:     8,
			ratingCount:       3,
			newestAlbumRating: 10,
			expectedScore:     90,
		},

		//Okay top rating, low average, middle rating count
		"50 Cent with no new album rating": {
			bestRating:        8,
			averageRating:     7,
			ratingCount:       7,
			newestAlbumRating: 0,
			expectedScore:     75,
		},
		"50 Cent with poorly-rated new album": {
			bestRating:        8,
			averageRating:     7,
			ratingCount:       7,
			newestAlbumRating: 5,
			expectedScore:     75,
		},
		"50 Cent with well-rated new album": {
			bestRating:        8,
			averageRating:     7,
			ratingCount:       7,
			newestAlbumRating: 8,
			expectedScore:     79,
		},
	}

	for testdescription, testcase := range testdata {
		score := calculateScore(testcase.bestRating, testcase.averageRating, testcase.ratingCount, testcase.newestAlbumRating)
		if score != testcase.expectedScore {
			t.Errorf("%s: expected %d but got %d", testdescription, testcase.expectedScore, score)
		}
	}
}
