package email

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_getSubjectLine(t *testing.T) {
	testcases := map[string]struct {
		ReleaseWeek string
		Expected string
	}{
		"No release week": {
			ReleaseWeek: "",
			Expected: "Newest releases from the week of " + time.Now().Format("2006-01-02"),
		},
		"With release week": {
			ReleaseWeek: "20200327",
			Expected: "Newest releases from the week of 2020-03-27",
		},
	}

	for testcase, testdata := range testcases {
		subject := GetNewReleasesSubjectLine(testdata.ReleaseWeek)
		assert.Equal(t, testdata.Expected, subject, testcase)
	}
}