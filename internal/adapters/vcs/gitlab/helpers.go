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
