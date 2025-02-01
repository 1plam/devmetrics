package routes

import (
	"github.com/gofiber/fiber/v2"
	"time"

	"devmetrics/internal/api/rest/handlers/vcs/github"
	"devmetrics/internal/api/rest/handlers/vcs/gitlab"
)

type Routes struct {
	githubHandler *github.Handler
	gitlabHandler *gitlab.Handler
}

func NewRoutes(githubHandler *github.Handler, gitlabHandler *gitlab.Handler) *Routes {
	return &Routes{
		githubHandler: githubHandler,
		gitlabHandler: gitlabHandler,
	}
}

func (r *Routes) Setup(app *fiber.App) {
	api := app.Group("/api/v1")

	r.setupVCSRoutes(api)
	r.setupHealthRoutes(api)
}

func (r *Routes) setupVCSRoutes(api fiber.Router) {
	vcsGroup := api.Group("/vcs")

	gitlabGroup := vcsGroup.Group("/gitlab/projects")
	gitlabGroup.Get("/:id", r.gitlabHandler.GetRepository)
	gitlabGroup.Get("/:id/commits", r.gitlabHandler.GetCommits)
	gitlabGroup.Get("/:id/merge-requests", r.gitlabHandler.GetPullRequests)

	githubGroup := vcsGroup.Group("/github/repositories")
	githubGroup.Get("/:owner/:name", r.githubHandler.GetRepository)
	githubGroup.Get("/:owner/:name/commits", r.githubHandler.GetCommits)
	githubGroup.Get("/:owner/:name/pull-requests", r.githubHandler.GetPullRequests)
}

func (r *Routes) setupHealthRoutes(api fiber.Router) {
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now(),
		})
	})
}
