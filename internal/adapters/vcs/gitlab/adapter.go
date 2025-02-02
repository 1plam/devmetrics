package gitlab

import (
	"context"
	"fmt"
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

func (a *Adapter) GetCommits(ctx context.Context, repo string, since, until time.Time, offset, limit int) ([]vcs.Commit, int64, error) {
	project, _, err := a.client.Projects.GetProject(repo, &gitlab.GetProjectOptions{}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get GitLab project for commits: %w", err)
	}

	// Count total commits by fetching until we hit a page with no commits
	totalCommits := int64(0)
	currentPage := 1
	for {
		countOpts := &gitlab.ListCommitsOptions{
			Since: gitlab.Time(since),
			Until: gitlab.Time(until),
			ListOptions: gitlab.ListOptions{
				Page:    currentPage,
				PerPage: 100, // Use max page size for counting
			},
		}

		commits, resp, err := a.client.Commits.ListCommits(project.ID, countOpts, gitlab.WithContext(ctx))
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count commits: %w", err)
		}

		if len(commits) == 0 {
			break
		}

		totalCommits += int64(len(commits))

		// If no next page, we're done
		if resp.NextPage == 0 {
			break
		}
		currentPage = resp.NextPage
	}

	// Now get the actual requested page
	opts := &gitlab.ListCommitsOptions{
		Since: gitlab.Time(since),
		Until: gitlab.Time(until),
		ListOptions: gitlab.ListOptions{
			Page:    (offset / limit) + 1,
			PerPage: limit,
		},
	}

	commits, _, err := a.client.Commits.ListCommits(project.ID, opts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list GitLab commits: %w", err)
	}

	var results []vcs.Commit
	resultChan := make(chan commitResult, len(commits))
	semaphore := make(chan struct{}, 10)

	for _, commit := range commits {
		go func(commit *gitlab.Commit) {
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			detailedCommit, _, err := a.client.Commits.GetCommit(
				project.ID,
				commit.ID,
				&gitlab.GetCommitOptions{},
				gitlab.WithContext(ctx),
			)
			if err != nil {
				resultChan <- commitResult{err: fmt.Errorf("failed to get commit details for %s: %w", commit.ID, err)}
				return
			}

			resultChan <- commitResult{commit: a.mapCommit(detailedCommit, repo)}
		}(commit)
	}

	for range commits {
		result := <-resultChan
		if result.err != nil {
			return nil, 0, result.err
		}
		results = append(results, result.commit)
	}

	return results, totalCommits, nil
}

func (a *Adapter) GetPullRequests(ctx context.Context, repo string, since, until time.Time, offset, limit int) ([]vcs.PullRequest, int64, error) {
	project, _, err := a.client.Projects.GetProject(repo, &gitlab.GetProjectOptions{}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get GitLab project for merge requests: %w", err)
	}

	opts := &gitlab.ListProjectMergeRequestsOptions{
		CreatedAfter:  gitlab.Time(since),
		CreatedBefore: gitlab.Time(until),
		ListOptions: gitlab.ListOptions{
			Page:    (offset / limit) + 1,
			PerPage: limit,
		},
	}

	mrs, resp, err := a.client.MergeRequests.ListProjectMergeRequests(project.ID, opts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list GitLab merge requests: %w", err)
	}

	var results []vcs.PullRequest
	resultChan := make(chan prResult, len(mrs))
	semaphore := make(chan struct{}, 10)

	for _, mr := range mrs {
		go func(mr *gitlab.MergeRequest) {
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			detailedMR, _, err := a.client.MergeRequests.GetMergeRequest(
				project.ID,
				mr.IID,
				&gitlab.GetMergeRequestsOptions{},
				gitlab.WithContext(ctx),
			)
			if err != nil {
				resultChan <- prResult{err: fmt.Errorf("failed to get merge request details for %d: %w", mr.IID, err)}
				return
			}

			resultChan <- prResult{pr: a.mapPullRequest(detailedMR, repo)}
		}(mr)
	}

	for range mrs {
		result := <-resultChan
		if result.err != nil {
			return nil, 0, result.err
		}
		results = append(results, result.pr)
	}

	return results, int64(resp.TotalItems), nil
}
