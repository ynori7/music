package allmusic

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetArtistDiscography_KingDiamond(t *testing.T) {
	//given
	reqCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if reqCount == 0 {
			dat, err := ioutil.ReadFile("testdata/king-diamond.html")
			require.NoError(t, err, "There was an error reading the test data file")
			rw.Write(dat)
			reqCount++
		} else if reqCount == 1 {
			dat, err := ioutil.ReadFile("testdata/king-diamond-ajax.html")
			require.NoError(t, err, "There was an error reading the test data file")
			rw.Write(dat)
			reqCount++
		} else {
			t.Errorf("Unexpected request")
		}
	}))
	defer server.Close()

	dicographyClient := DiscographyClient{server.Client()}

	//when
	discography, err := dicographyClient.GetArtistDiscography(server.URL)

	//then
	require.NoError(t, err, "There was an error getting the discography")
	assert.Equal(t, "King Diamond", discography.Artist.Name)
	assert.Equal(t, 18, len(discography.Albums))
	assert.Equal(t, "https://rovimusic.rovicorp.com/image.jpg?c=fGwYdlDmR9-V_0hsFevyBN_M69_UI9rrJSVvWL2-yAg=&f=2", discography.Albums[0].Image)
	assert.Equal(t, 9, discography.BestRating)
}

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
