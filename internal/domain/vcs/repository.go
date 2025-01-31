package vcs

import "time"

type Repository struct {
	ID            string
	Name          string
	FullName      string
	DefaultBranch string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Description   string
	Language      string
	Private       bool
}
