package allmusic

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ynori7/hulksmash/anonymizer"
	"github.com/ynori7/music/config"
)

func Test_GetNewReleases(t *testing.T) {
	//given
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		dat, err := os.ReadFile("testdata/newreleases.html")
		require.NoError(t, err, "There was an error reading the test data file")
		rw.Write(dat)
	}))
	defer server.Close()

	conf := config.Config{MainGenres: []string{"Rock", "Rap"}}
	newReleasesClient := ReleasesClient{
		httpClient:    server.Client(),
		conf:          conf,
		reqAnonymizer: anonymizer.New(12345),
	}

	//when
	releases, err := newReleasesClient.GetPotentiallyInterestingNewReleases(server.URL)

	//then
	require.NoError(t, err, "There was an error getting the releases")
	assert.Equal(t, 89, len(releases))
}
