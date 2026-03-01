package metadata

import (
	"context"
	"strings"

	worker "github.com/keithics/devops-dashboard/api/internal/metadata/provider"
	normalizeutil "github.com/keithics/devops-dashboard/api/internal/utils/normalize"
)

func NewService(workerService *worker.Service) *Service {
	return &Service{
		worker: workerService,
	}
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (s *Service) Search(ctx context.Context, provider worker.ProviderName, query string, opts worker.SearchOpts) ([]worker.SearchHit, error) {
	items, err := s.worker.Search(ctx, provider, query, opts)
	if err != nil {
		return nil, err
	}
	return filterTitleContains(query, items), nil
}

func (s *Service) GetShow(ctx context.Context, provider worker.ProviderName, externalID string) (worker.Show, error) {
	return s.worker.GetShow(ctx, provider, externalID)
}

func (s *Service) ListEpisodes(ctx context.Context, provider worker.ProviderName, externalID string, opts worker.ListEpisodesOpts) ([]worker.Episode, error) {
	return s.worker.ListEpisodes(ctx, provider, externalID, opts)
}

func filterTitleContains(query string, items []worker.SearchHit) []worker.SearchHit {
	normalizedQuery := normalizeutil.LowerString(query)
	if normalizedQuery == "" {
		return []worker.SearchHit{}
	}

	filtered := make([]worker.SearchHit, 0, len(items))
	for _, item := range items {
		if !hasTitleContainsMatch(item, normalizedQuery) {
			continue
		}
		filtered = append(filtered, item)
	}
	return filtered
}

func hasTitleContainsMatch(item worker.SearchHit, normalizedQuery string) bool {
	if strings.Contains(normalizeutil.LowerString(item.TitlePreferred), normalizedQuery) {
		return true
	}
	if item.TitleOriginal != nil && strings.Contains(normalizeutil.LowerString(*item.TitleOriginal), normalizedQuery) {
		return true
	}
	for _, altTitle := range item.AltTitles {
		if strings.Contains(normalizeutil.LowerString(altTitle), normalizedQuery) {
			return true
		}
	}
	return false
}
