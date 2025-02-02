package gitlab

import (
	"context"
	"fmt"
	"github.com/xanzy/go-gitlab"
	"net/url"
	"strings"
)

func (a *Adapter) GetProjectID(ctx context.Context, repo string) (int, error) {
	path := url.PathEscape(repo)

	project, _, err := a.client.Projects.GetProject(path, nil, gitlab.WithContext(ctx))
	if err != nil {
		return 0, fmt.Errorf("getting project ID: %w", err)
	}

	return project.ID, nil
}

func (a *Adapter) doesRepositoryExist(ctx context.Context, repo string) (bool, error) {
	path := url.QueryEscape(repo)

	_, resp, err := a.client.Projects.GetProject(path, &gitlab.GetProjectOptions{
		Statistics: gitlab.Bool(false),
	}, gitlab.WithContext(ctx))

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (a *Adapter) buildProjectPath(repo string) string {
	repo = strings.Trim(repo, "/")
	repo = strings.ReplaceAll(repo, "//", "/")

	return repo
}

func (a *Adapter) findLastPage(ctx context.Context, projectID interface{}, opts *gitlab.ListCommitsOptions) (int, int64, error) {
	// Start with a small request to check if we have any commits
	firstPageOpts := *opts
	firstPageOpts.Page = 1
	firstPageOpts.PerPage = 1

	commits, _, err := a.client.Commits.ListCommits(projectID, &firstPageOpts, gitlab.WithContext(ctx))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get first page: %w", err)
	}

	if len(commits) == 0 {
		return 0, 0, nil
	}

	// Initialize binary search bounds
	// Start with a reasonably large upper bound
	left, right := 1, 1000
	lastValidPage := 1
	totalCommits := int64(0)

	// Binary search for the last valid page
	for left <= right {
		mid := (left + right) / 2

		pageOpts := *opts
		pageOpts.Page = mid
		pageOpts.PerPage = 100 // Use max page size for efficiency

		commits, _, err := a.client.Commits.ListCommits(projectID, &pageOpts, gitlab.WithContext(ctx))
		if err != nil {
			// If we get a 404, page is too high
			if strings.Contains(err.Error(), "404") {
				right = mid - 1
				continue
			}
			return 0, 0, fmt.Errorf("failed to list commits during binary search: %w", err)
		}

		if len(commits) == 0 {
			// This page is empty, look in lower half
			right = mid - 1
		} else {
			// This page has commits, look in upper half
			lastValidPage = mid
			left = mid + 1
			// Update total commits for the last valid page we found
			if lastValidPage == mid {
				totalCommits = int64((mid-1)*100 + len(commits))
			}
		}
	}

	// Get exact count for the last page
	lastPageOpts := *opts
	lastPageOpts.Page = lastValidPage
	lastPageOpts.PerPage = 100

	commits, _, err = a.client.Commits.ListCommits(projectID, &lastPageOpts, gitlab.WithContext(ctx))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get last page: %w", err)
	}

	// Calculate total commits
	totalCommits = int64((lastValidPage-1)*100 + len(commits))

	return lastValidPage, totalCommits, nil
}
