package metadata

import (
	"context"

	jobshow "github.com/keithics/devops-dashboard/api/internal/job/show"
	worker "github.com/keithics/devops-dashboard/api/internal/worker/metadata"
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
	return s.worker.Search(ctx, provider, query, opts)
}

func (s *Service) GetShow(ctx context.Context, provider worker.ProviderName, externalID string) (worker.Show, error) {
	return s.worker.GetShow(ctx, provider, externalID)
}

func (s *Service) ListEpisodes(ctx context.Context, provider worker.ProviderName, externalID string, opts worker.ListEpisodesOpts) ([]worker.Episode, error) {
	return s.worker.ListEpisodes(ctx, provider, externalID, opts)
}

func (s *Service) AddShow(ctx context.Context, req AddShowRequest) (AddShowResponse, error) {
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
