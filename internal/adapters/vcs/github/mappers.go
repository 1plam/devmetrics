package github

import (
	"devmetrics/internal/domain/vcs"
	"github.com/google/go-github/v45/github"
	"strconv"
)

func (a *Adapter) mapCommit(ghCommit *github.RepositoryCommit, repoID string) vcs.Commit {
	if ghCommit == nil || ghCommit.Commit == nil {
		return vcs.Commit{}
	}

	return vcs.Commit{
		SHA:          ghCommit.GetSHA(),
		Message:      ghCommit.Commit.GetMessage(),
		AuthorName:   ghCommit.Commit.Author.GetName(),
		AuthorEmail:  ghCommit.Commit.Author.GetEmail(),
		CommittedAt:  ghCommit.Commit.Author.GetDate(),
		ChangedFiles: ghCommit.GetStats().GetTotal(),
		Additions:    ghCommit.GetStats().GetAdditions(),
		Deletions:    ghCommit.GetStats().GetDeletions(),
		RepositoryID: repoID,
	}
}

func (a *Adapter) mapPullRequest(pr *github.PullRequest, repoID string) vcs.PullRequest {
	if pr == nil {
		return vcs.PullRequest{}
	}

	return vcs.PullRequest{
		Number:       pr.GetNumber(),
		Title:        pr.GetTitle(),
		State:        pr.GetState(),
		CreatedAt:    pr.GetCreatedAt(),
		UpdatedAt:    pr.GetUpdatedAt(),
		ClosedAt:     pr.ClosedAt,
		MergedAt:     pr.MergedAt,
		AuthorName:   pr.User.GetLogin(),
		ReviewCount:  pr.GetReviewComments(),
		CommitCount:  pr.GetCommits(),
		ChangedFiles: pr.GetChangedFiles(),
		Additions:    pr.GetAdditions(),
		Deletions:    pr.GetDeletions(),
		RepositoryID: repoID,
	}
}

func (a *Adapter) mapRepository(repo *github.Repository) *vcs.Repository {
	if repo == nil {
		return nil
	}

	return &vcs.Repository{
		ID:            strconv.FormatInt(repo.GetID(), 10),
		Name:          repo.GetName(),
		FullName:      repo.GetFullName(),
		DefaultBranch: repo.GetDefaultBranch(),
		CreatedAt:     repo.GetCreatedAt().Time,
		UpdatedAt:     repo.GetUpdatedAt().Time,
		Description:   repo.GetDescription(),
		Language:      repo.GetLanguage(),
		Private:       repo.GetPrivate(),
	}
}
