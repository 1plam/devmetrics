package main

import (
	"context"
	"devmetrics/internal/api/rest/server"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	adapter "devmetrics/internal/adapters/vcs"
	"devmetrics/internal/api/rest/handlers/vcs/github"
	"devmetrics/internal/api/rest/handlers/vcs/gitlab"
	"devmetrics/internal/api/rest/routes"
	"devmetrics/internal/app"
	"devmetrics/internal/config"
	domain "devmetrics/internal/domain/vcs"
	"devmetrics/internal/services/vcs"
	"github.com/gofiber/fiber/v2/log"
	"go.uber.org/dig"
)

func main() {
	if err := run(buildContainer()); err != nil {
		log.Fatal(err)
	}
}

func providers() []interface{} {
	return []interface{}{
		// Config
		config.NewConfig,
		provideVCSConfig,

		// VCS
		adapter.NewFactory,
		provideVCSService,

		// HTTP Handlers
		provideGitHubHandler,
		provideGitLabHandler,
		provideRoutes,
		server.NewServer,

		// Application
		app.NewApplication,
	}
}

func provideVCSConfig(cfg *config.Config) config.VCSConfig {
	return cfg.VCS
}

func provideVCSService(factory *adapter.Factory) (*vcs.Service, error) {
	providers, err := factory.CreateProviders()
	if err != nil {
		return vcs.NewService(make(map[domain.ProviderType]domain.Provider)), nil
	}
	return vcs.NewService(providers), nil
}

func provideGitHubHandler(service *vcs.Service) *github.Handler {
	return github.NewHandler(service)
}

func provideGitLabHandler(service *vcs.Service) *gitlab.Handler {
	return gitlab.NewHandler(service)
}

func provideRoutes(githubHandler *github.Handler, gitlabHandler *gitlab.Handler) *routes.Routes {
	return routes.NewRoutes(githubHandler, gitlabHandler)
}

func buildContainer() *dig.Container {
	container := dig.New()

	for _, provider := range providers() {
		if err := container.Provide(provider); err != nil {
			log.Fatalf("dependency injection error: %v", err)
		}
	}

	return container
}

func run(container *dig.Container) error {
	return container.Invoke(func(app *app.Application) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		errChan := make(chan error, 1)
		go func() {
			errChan <- app.Start(ctx)
		}()

		return handleShutdown(ctx, app, errChan)
	})
}

func handleShutdown(ctx context.Context, app *app.Application, errChan <-chan error) error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return fmt.Errorf("application error: %w", err)
	case _ = <-shutdown:
		return app.Shutdown(ctx)
	}
}
