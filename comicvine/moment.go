package comicvine

import (
	"fmt"
	"time"
)

// @return a list of months, of the form YYYY-MM
func months(since int) []string {
	t := time.Now()
	cutoffDate := time.Date(since, 1, 1, 0, 0, 0, 0, time.Local)

	var months []string

	for t.After(cutoffDate) {
		t = t.AddDate(0, -1, 0)
		str := fmt.Sprintf("%d-%02d", t.Year(), int(t.Month()))
		months = append(months, str)
	}

	return months
}
