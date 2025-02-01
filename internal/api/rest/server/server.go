package server

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"devmetrics/internal/api/rest/handlers/vcs/github"
	"devmetrics/internal/api/rest/handlers/vcs/gitlab"
	"devmetrics/internal/api/rest/middleware"
	"devmetrics/internal/api/rest/routes"
	"devmetrics/internal/config"
)

type Server struct {
	app    *fiber.App
	config *config.Config
	routes *routes.Routes
	addr   string
}

func NewServer(
	config *config.Config,
	githubHandler *github.Handler,
	gitlabHandler *gitlab.Handler,
) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	addr := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
	routes := routes.NewRoutes(githubHandler, gitlabHandler)

	return &Server{
		app:    app,
		config: config,
		routes: routes,
		addr:   addr,
	}
}

func (s *Server) setupMiddleware() {
	s.app.Use(logger.New())
	s.app.Use(middleware.Cors())
	s.app.Use(middleware.RequestID())
}

func (s *Server) setupRoutes() {
	s.routes.Setup(s.app)
}

func (s *Server) Start() error {
	s.setupMiddleware()
	s.setupRoutes()

	return s.app.Listen(s.addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}
