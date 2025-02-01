package github

type RepoRequest struct {
	Owner string `params:"owner" validate:"required"`
	Name  string `params:"name" validate:"required"`
}
