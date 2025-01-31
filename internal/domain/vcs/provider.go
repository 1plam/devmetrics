package vcs

import (
	"context"
	"time"
)

// Provider defines the interface for VCS (Version Control System) operations
type Provider interface {
	// GetCommits retrieves commits for a repository within a time range
	GetCommits(ctx context.Context, repo string, since, until time.Time) ([]Commit, error)

	// GetPullRequests retrieves pull requests for a repository within a time range
	GetPullRequests(ctx context.Context, repo string, since, until time.Time) ([]PullRequest, error)

	// GetRepository retrieves repository information
	GetRepository(ctx context.Context, repo string) (*Repository, error)
}

// ProviderType represents the type of VCS provider (GitHub, GitLab, etc.)
type ProviderType string

const (
	ProviderGitHub    ProviderType = "github"
	ProviderGitLab    ProviderType = "gitlab"
	ProviderBitbucket ProviderType = "bitbucket"
)
