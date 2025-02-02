package github

import (
	"context"
	"devmetrics/internal/adapters/vcs/common"
	"devmetrics/internal/config"
	"fmt"
	"log"
	"time"

	"devmetrics/internal/domain/vcs"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

type Adapter struct {
	client *github.Client
	config config.GitHubConfig
}

var _ vcs.Provider = (*Adapter)(nil)

func NewAdapter(cfg config.GitHubConfig) (*Adapter, error) {
	if cfg.Token == "" {
		return nil, fmt.Errorf("github token is required")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.Token})
	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)

	return &Adapter{
		client: client,
		config: cfg,
	}, nil
}

func (a *Adapter) GetRepository(ctx context.Context, repo string) (*vcs.Repository, error) {
	owner, repoName := common.ParseRepoString(repo)

	repository, _, err := a.client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		return nil, fmt.Errorf("getting repository: %w", err)
	}

	return a.mapRepository(repository), nil
}

func (a *Adapter) GetCommits(ctx context.Context, repo string, since, until time.Time, offset, limit int) ([]vcs.Commit, int64, error) {
	owner, repoName := common.ParseRepoString(repo)

	// Ensure limit is set and within bounds
	if limit <= 0 {
		limit = 30 // Default to 30 if not specified
	}
	if limit > 100 {
		limit = 100 // GitHub's max is 100
	}

	// Calculate which page we need to start from based on offset and limit
	page := (offset / limit) + 1

	opts := &github.CommitsListOptions{
		Since: since,
		Until: until,
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: limit,
		},
	}

	commits, resp, err := a.client.Repositories.ListCommits(ctx, owner, repoName, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("listing commits: %w", err)
	}

	var results []vcs.Commit
	for _, commit := range commits {
		results = append(results, a.mapCommit(commit, repo))
	}

	// Calculate total items
	total := int64(0)
	if links := resp.Response.Header.Get("Link"); links != "" {
		total = a.extractTotalFromLink(links)
	} else if resp.LastPage > 0 {
		// If we have LastPage but no Link header
		total = int64(resp.LastPage * opts.PerPage)
	} else if len(results) < opts.PerPage {
		// If we got less results than requested, this must be the total
		total = int64(len(results))
	}

	return results, total, nil
}

func (a *Adapter) GetPullRequests(ctx context.Context, repo string, since, until time.Time, offset, limit int) ([]vcs.PullRequest, int64, error) {
	owner, repoName := common.ParseRepoString(repo)

	page := (offset / limit) + 1

	opts := &github.PullRequestListOptions{
		State: "all",
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: limit,
		},
	}

	prs, resp, err := a.client.PullRequests.List(ctx, owner, repoName, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("listing pull requests: %w", err)
	}

	var results []vcs.PullRequest
	for _, pr := range prs {
		if isInTimeRange(pr.CreatedAt, since, until) {
			details, _, err := a.client.PullRequests.Get(ctx, owner, repoName, pr.GetNumber())
			if err != nil {
				log.Printf("Failed to fetch PR details: pr_number=%d, error=%v", pr.GetNumber(), err)
				continue
			}
			results = append(results, a.mapPullRequest(details, repo))
		}
	}

	total := int64(0)
	if links := resp.Response.Header.Get("Link"); links != "" {
		total = a.extractTotalFromLink(links)
	}

	return results, total, nil
}
