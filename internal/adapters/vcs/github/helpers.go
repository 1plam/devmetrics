package github

import (
	"time"
)

func isInTimeRange(t *time.Time, since, until time.Time) bool {
	if t == nil {
		return false
	}
	return t.After(since) && t.Before(until)
}
