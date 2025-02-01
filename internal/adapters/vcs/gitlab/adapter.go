package gitlab

import (
	"context"
	"fmt"
	"sync"
	"time"

	"devmetrics/internal/config"
	"devmetrics/internal/domain/vcs"
	"github.com/xanzy/go-gitlab"
)

type Adapter struct {
	client *gitlab.Client
}

type commitResult struct {
	commit vcs.Commit
	err    error
}

type prResult struct {
	pr  vcs.PullRequest
	err error
}

func NewAdapter(cfg config.GitLabConfig) (*Adapter, error) {
	client, err := gitlab.NewClient(cfg.Token, gitlab.WithBaseURL(cfg.BaseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to create GitLab client: %w", err)
	}

	return &Adapter{
		client: client,
	}, nil
}

func (a *Adapter) GetRepository(ctx context.Context, repo string) (*vcs.Repository, error) {
	project, _, err := a.client.Projects.GetProject(repo, &gitlab.GetProjectOptions{}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to get GitLab project: %w", err)
	}

	return a.mapRepository(project), nil
}

func (a *Adapter) GetCommits(ctx context.Context, repo string, since, until time.Time) ([]vcs.Commit, error) {
	project, _, err := a.client.Projects.GetProject(repo, &gitlab.GetProjectOptions{}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to get GitLab project for commits: %w", err)
	}

	opts := &gitlab.ListCommitsOptions{
		Since: gitlab.Time(since),
		Until: gitlab.Time(until),
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
		},
	}

	var allCommits []vcs.Commit
	var mu sync.Mutex
	semaphore := make(chan struct{}, 10)

	// Process a single commit
	processCommit := func(commit *gitlab.Commit, resultChan chan<- commitResult) {
		semaphore <- struct{}{}
		defer func() { <-semaphore }()

		detailedCommit, _, err := a.client.Commits.GetCommit(project.ID, commit.ID, &gitlab.GetCommitOptions{}, gitlab.WithContext(ctx))
		if err != nil {
			resultChan <- commitResult{err: fmt.Errorf("failed to get commit details for %s: %w", commit.ID, err)}
			return
		}

		resultChan <- commitResult{commit: a.mapCommit(detailedCommit, repo)}
	}

	for {
		commits, resp, err := a.client.Commits.ListCommits(project.ID, opts, gitlab.WithContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("failed to list GitLab commits: %w", err)
		}

		// Create result channel for this batch
		resultChan := make(chan commitResult, len(commits))

		// Launch goroutines for each commit
		for _, commit := range commits {
			go processCommit(commit, resultChan)
		}

		// Collect results
		for range commits {
			result := <-resultChan
			if result.err != nil {
				return nil, result.err
			}

			mu.Lock()
			allCommits = append(allCommits, result.commit)
			mu.Unlock()
		}

		if resp.CurrentPage >= resp.TotalPages {
			break
		}
		opts.Page = resp.NextPage
	}

	return allCommits, nil
}

func (a *Adapter) GetPullRequests(ctx context.Context, repo string, since, until time.Time) ([]vcs.PullRequest, error) {
	project, _, err := a.client.Projects.GetProject(repo, &gitlab.GetProjectOptions{}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to get GitLab project for merge requests: %w", err)
	}

	opts := &gitlab.ListProjectMergeRequestsOptions{
		CreatedAfter:  gitlab.Time(since),
		CreatedBefore: gitlab.Time(until),
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
		},
	}

	var allMRs []vcs.PullRequest
	var mu sync.Mutex
	semaphore := make(chan struct{}, 10)

	// Process a single merge request
	processMR := func(mr *gitlab.MergeRequest, resultChan chan<- prResult) {
		semaphore <- struct{}{}
		defer func() { <-semaphore }()

		detailedMR, _, err := a.client.MergeRequests.GetMergeRequest(project.ID, mr.IID, &gitlab.GetMergeRequestsOptions{}, gitlab.WithContext(ctx))
		if err != nil {
			resultChan <- prResult{err: fmt.Errorf("failed to get merge request details for %d: %w", mr.IID, err)}
			return
		}

		resultChan <- prResult{pr: a.mapPullRequest(detailedMR, repo)}
	}

	for {
		mrs, resp, err := a.client.MergeRequests.ListProjectMergeRequests(project.ID, opts, gitlab.WithContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("failed to list GitLab merge requests: %w", err)
		}

		// Create result channel for this batch
		resultChan := make(chan prResult, len(mrs))

		// Launch goroutines for each merge request
		for _, mr := range mrs {
			go processMR(mr, resultChan)
		}

		// Collect results
		for range mrs {
			result := <-resultChan
			if result.err != nil {
				return nil, result.err
			}

			mu.Lock()
			allMRs = append(allMRs, result.pr)
			mu.Unlock()
		}

		if resp.CurrentPage >= resp.TotalPages {
			break
		}
		opts.Page = resp.NextPage
	}

	return allMRs, nil
}
