package gitlab

import (
	"devmetrics/internal/domain/vcs"
	"github.com/xanzy/go-gitlab"
	"strconv"
)

func (a *Adapter) mapCommit(glCommit *gitlab.Commit, repoID string) vcs.Commit {
	if glCommit == nil {
		return vcs.Commit{}
	}

	return vcs.Commit{
		SHA:          glCommit.ID,
		Message:      glCommit.Message,
		AuthorName:   glCommit.AuthorName,
		AuthorEmail:  glCommit.AuthorEmail,
		CommittedAt:  glCommit.CommittedDate.UTC(),
		ChangedFiles: glCommit.Stats.Total,
		Additions:    glCommit.Stats.Additions,
		Deletions:    glCommit.Stats.Deletions,
		RepositoryID: repoID,
	}
}

func (a *Adapter) mapPullRequest(mr *gitlab.MergeRequest, repoID string) vcs.PullRequest {
	if mr == nil {
		return vcs.PullRequest{}
	}

	changesCount, _ := strconv.Atoi(mr.ChangesCount)

	return vcs.PullRequest{
		Number:       mr.IID,
		Title:        mr.Title,
		State:        mr.State,
		CreatedAt:    mr.CreatedAt.UTC(),
		UpdatedAt:    mr.UpdatedAt.UTC(),
		ClosedAt:     mr.ClosedAt,
		MergedAt:     mr.MergedAt,
		AuthorName:   mr.Author.Username,
		ReviewCount:  mr.UserNotesCount,
		CommitCount:  mr.DivergedCommitsCount,
		ChangedFiles: changesCount,
		Additions:    0,
		Deletions:    0,
		RepositoryID: repoID,
	}
}

func (a *Adapter) mapRepository(project *gitlab.Project) *vcs.Repository {
	if project == nil {
		return nil
	}

	return &vcs.Repository{
		ID:            strconv.Itoa(project.ID),
		Name:          project.Name,
		FullName:      project.NameWithNamespace,
		DefaultBranch: project.DefaultBranch,
		CreatedAt:     project.CreatedAt.UTC(),
		UpdatedAt:     project.LastActivityAt.UTC(),
		Description:   project.Description,
		Language:      "", // GitLab API doesn't provide primary language info
		Private:       project.Visibility != "public",
	}
}
