package metadata

import (
	"context"
	"strings"

	jobshow "github.com/keithics/devops-dashboard/api/internal/job/show"
	worker "github.com/keithics/devops-dashboard/api/internal/metadata/provider"
	normalizeutil "github.com/keithics/devops-dashboard/api/internal/utils/normalize"
)

func NewService(workerService *worker.Service, jobShowService *jobshow.Service) *Service {
	return &Service{
		worker:  workerService,
		jobShow: jobShowService,
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

func (s *Service) AddShow(ctx context.Context, req AddShowRequest) (AddShowResponse, error) {
	provider, err := providerFromExternalID(req.ExternalID)
	if err != nil {
		return AddShowResponse{}, err
	}

	episodes, err := s.worker.ListEpisodes(ctx, provider, req.ExternalID, worker.ListEpisodesOpts{
		Page:  1,
		Limit: 100,
	})
	if err != nil {
		return AddShowResponse{}, err
	}

	episodeInputs := make([]jobshow.EpisodeInput, 0, len(episodes))
	for _, episode := range episodes {
		episodeInputs = append(episodeInputs, jobshow.EpisodeInput{
			ExternalID:     episode.ExternalID,
			SeasonNumber:   episode.SeasonNumber,
			EpisodeNumber:  episode.EpisodeNumber,
			Title:          episode.Title,
			AirDate:        episode.AirDate,
			RuntimeMinutes: episode.RuntimeMinutes,
		})
	}

	created, err := s.jobShow.EnqueueFromSearchResult(ctx, jobshow.EnqueueFromSearchResultParams{
		ExternalID:     req.ExternalID,
		TitlePreferred: req.TitlePreferred,
		TitleOriginal:  req.TitleOriginal,
		AltTitles:      req.AltTitles,
		Type:           req.Type,
		Status:         req.Status,
		Synopsis:       req.Synopsis,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		PosterURL:      req.PosterUrl,
		BannerURL:      req.BannerUrl,
		SeasonCount:    req.SeasonCount,
		EpisodeCount:   req.EpisodeCount,
		Episodes:       episodeInputs,
	})
	if err != nil {
		return AddShowResponse{}, err
	}

	return AddShowResponse{
		InternalShowID:    created.InternalShowID,
		InternalJobShowID: created.InternalJobShowID,
		Status:            created.Status,
		RetryCount:        created.RetryCount,
	}, nil
}

func providerFromExternalID(externalID string) (worker.ProviderName, error) {
	value := strings.ToLower(strings.TrimSpace(externalID))
	switch {
	case strings.HasPrefix(value, "anidb:"):
		return worker.ProviderAniDB, nil
	case strings.HasPrefix(value, "anilist:"):
		return worker.ProviderAniList, nil
	case strings.HasPrefix(value, "tvdb:"):
		return worker.ProviderTVDB, nil
	default:
		return "", errExternalIDMustBePrefixed
	}
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
