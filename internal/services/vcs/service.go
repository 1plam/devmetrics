package vcs

import (
	"context"
	"fmt"
	"time"

	"devmetrics/internal/domain/vcs"
)

type Service struct {
	providers map[vcs.ProviderType]vcs.Provider
}

func NewService(providers map[vcs.ProviderType]vcs.Provider) *Service {
	return &Service{
		providers: providers,
	}
}

func (s *Service) GetRepository(ctx context.Context, providerType vcs.ProviderType, repo string) (*vcs.Repository, error) {
	provider, ok := s.providers[providerType]
	if !ok {
		return nil, fmt.Errorf("unsupported VCS provider: %s", providerType)
	}

	return provider.GetRepository(ctx, repo)
}

func (s *Service) GetCommits(
	ctx context.Context,
	providerType vcs.ProviderType,
	repo string,
	since, until time.Time,
	offset, limit int,
) ([]vcs.Commit, int64, error) {
	provider, ok := s.providers[providerType]
	if !ok {
		return nil, 0, fmt.Errorf("unsupported VCS provider: %s", providerType)
	}

	commits, total, err := provider.GetCommits(ctx, repo, since, until, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get commits: %w", err)
	}

	return commits, total, nil
}

func (s *Service) GetPullRequests(
	ctx context.Context,
	providerType vcs.ProviderType,
	repo string,
	since, until time.Time,
	offset, limit int,
) ([]vcs.PullRequest, int64, error) {
	provider, ok := s.providers[providerType]
	if !ok {
		return nil, 0, fmt.Errorf("unsupported VCS provider: %s", providerType)
	}

	prs, total, err := provider.GetPullRequests(ctx, repo, since, until, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get pull requests: %w", err)
	}

	return prs, total, nil
}
