package rateyourmusic

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetArtistDiscography_KingDiamond(t *testing.T) {
	//given
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		dat, err := ioutil.ReadFile("testdata/king-diamond.html")
		require.NoError(t, err, "There was an error reading the test data file")
		rw.Write(dat)
	}))
	defer server.Close()

	dicographyClient := DiscographyClient{server.Client()}

	//when
	discography, err := dicographyClient.GetArtistDiscography(server.URL)

	//then
	require.NoError(t, err, "There was an error getting the discography")
	assert.Equal(t, "King Diamond", discography.Artist)
	assert.Equal(t, 13, len(discography.Albums))
}

func Test_GetArtistDiscography_blah(t *testing.T) {
	artists := []string{"Joe Satriani",
		"RJD2",
		"The Black Dahlia Murder",
		"Shabazz Palaces",
		"Wildfire",
		"DJ Screw",
		"Oranssi Pazuzu",
	}

	discographyClient := NewDiscographyClient()
	for _, artist := range artists {
		fmt.Println(discographyClient.GetArtistDiscography(BuildRymLink(artist)))
	}
}