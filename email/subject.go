package email

import (
	"fmt"
	"time"
)

func GetNewReleasesSubjectLine(releaseWeek string) string {
	date := ""
	if releaseWeek != "" {
		parsed, err := time.Parse("20060102", releaseWeek)
		if err == nil {
			date = parsed.Format("2006-01-02")
		}
	}

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	return fmt.Sprintf("Newest releases from the week of %s", date)
}
