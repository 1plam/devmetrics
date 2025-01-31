package vcs

import "time"

type PullRequest struct {
	Number       int
	Title        string
	State        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ClosedAt     *time.Time
	MergedAt     *time.Time
	AuthorName   string
	ReviewCount  int
	CommitCount  int
	ChangedFiles int
	Additions    int
	Deletions    int
	RepositoryID string
}
