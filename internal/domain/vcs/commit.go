package vcs

import "time"

type Commit struct {
	SHA          string
	Message      string
	AuthorName   string
	AuthorEmail  string
	CommittedAt  time.Time
	ChangedFiles int
	Additions    int
	Deletions    int
	RepositoryID string
}
