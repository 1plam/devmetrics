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

func (a *Adapter) GetCommits(ctx context.Context, repo string, since, until time.Time) ([]vcs.Commit, error) {
	owner, repoName := common.ParseRepoString(repo)

	opts := &github.CommitsListOptions{
		Since: since,
		Until: until,
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allCommits []vcs.Commit
	for {
		commits, resp, err := a.client.Repositories.ListCommits(ctx, owner, repoName, opts)
		if err != nil {
			return nil, fmt.Errorf("listing commits: %w", err)
		}

		for _, commit := range commits {
			allCommits = append(allCommits, a.mapCommit(commit, repo))
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allCommits, nil
}

func (a *Adapter) GetPullRequests(ctx context.Context, repo string, since, until time.Time) ([]vcs.PullRequest, error) {
	owner, repoName := common.ParseRepoString(repo)

	opts := &github.PullRequestListOptions{
		State: "all",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allPRs []vcs.PullRequest
	for {
		prs, resp, err := a.client.PullRequests.List(ctx, owner, repoName, opts)
		if err != nil {
			return nil, fmt.Errorf("listing pull requests: %w", err)
		}

		for _, pr := range prs {
			if isInTimeRange(pr.CreatedAt, since, until) {
				details, _, err := a.client.PullRequests.Get(ctx, owner, repoName, pr.GetNumber())
				if err != nil {
					log.Printf("Failed to fetch PR details: pr_number=%d, error=%v", pr.GetNumber(), err)
					continue
				}
				allPRs = append(allPRs, a.mapPullRequest(details, repo))
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allPRs, nil
}

func (a *Adapter) GetRepository(ctx context.Context, repo string) (*vcs.Repository, error) {
	owner, repoName := common.ParseRepoString(repo)

	repository, _, err := a.client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		return nil, fmt.Errorf("getting repository: %w", err)
	}

	return a.mapRepository(repository), nil
}
