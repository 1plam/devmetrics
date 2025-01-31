package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"

	vcshandler "devmetrics/internal/api/rest/handlers/vcs"
)

type Routes struct {
	vcsHandler *vcshandler.Handler
}

func NewRoutes(vcsHandler *vcshandler.Handler) *Routes {
	return &Routes{
		vcsHandler: vcsHandler,
	}
}

func (r *Routes) Setup(app *fiber.App) {
	// API version group
	api := app.Group("/api/v1")

	// Setup different route groups
	r.setupVCSRoutes(api)
	r.setupHealthRoutes(api)
}

func (r *Routes) setupVCSRoutes(api fiber.Router) {
	vcsGroup := api.Group("/vcs")

	vcsGroup.Get("/:provider/repositories/:owner/:name", r.vcsHandler.GetRepository)
	vcsGroup.Get("/:provider/repositories/:owner/:name/commits", r.vcsHandler.GetCommits)
	vcsGroup.Get("/:provider/repositories/:owner/:name/pull-requests", r.vcsHandler.GetPullRequests)
}

func (r *Routes) setupHealthRoutes(api fiber.Router) {
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now(),
		})
	})
}
