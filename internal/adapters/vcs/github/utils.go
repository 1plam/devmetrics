package github

import (
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func isInTimeRange(t *time.Time, since, until time.Time) bool {
	if t == nil {
		return false
	}
	return t.After(since) && t.Before(until)
}

func (a *Adapter) extractTotalFromLink(links string) int64 {
	parts := strings.Split(links, ",")

	var lastPage int64 = 1
	var perPage int64 = 30 // default GitHub value

	for _, part := range parts {
		part = strings.TrimSpace(part)

		// Extract URL
		urlStart := strings.Index(part, "<") + 1
		urlEnd := strings.Index(part, ">")
		if urlStart <= 0 || urlEnd <= urlStart {
			continue
		}

		urlStr := part[urlStart:urlEnd]

		// Parse the URL
		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			log.Printf("Failed to parse URL from Link header: %v", err)
			continue
		}

		// Get query parameters
		query := parsedURL.Query()

		// Get per_page value
		if perPageStr := query.Get("per_page"); perPageStr != "" {
			if pp, err := strconv.ParseInt(perPageStr, 10, 64); err == nil {
				perPage = pp
			}
		}

		// Check if this is the last page link
		if strings.Contains(part, `rel="last"`) {
			if pageStr := query.Get("page"); pageStr != "" {
				if p, err := strconv.ParseInt(pageStr, 10, 64); err == nil {
					lastPage = p
				}
			}
		}
	}

	// If we have both next and last pointing to the same or lower page,
	// we're probably on the only page
	if lastPage <= 1 {
		return perPage
	}

	total := lastPage * perPage
	return total
}
