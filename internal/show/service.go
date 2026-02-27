package show

import (
	"context"

	"github.com/keithics/devops-dashboard/api/internal/db/sqlc"
)

func NewHandler(q *sqlc.Queries) *Handler {
	return &Handler{
		svc: &Service{q: q},
	}
}

func (s *Service) CreateShow(ctx context.Context, req createShowRequest) (sqlc.Show, error) {
	externalIDs, err := marshalExternalID(req.ExternalID)
	if err != nil {
		return sqlc.Show{}, err
	}

	return s.q.CreateShow(ctx, sqlc.CreateShowParams{
		TitlePreferred: req.TitlePreferred,
		TitleOriginal:  req.TitleOriginal,
		AltTitles:      req.AltTitles,
		Type:           req.Type,
		Status:         req.Status,
		Synopsis:       req.Synopsis,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		PosterUrl:      req.PosterUrl,
		BannerUrl:      req.BannerUrl,
		SeasonCount:    req.SeasonCount,
		EpisodeCount:   req.EpisodeCount,
		ExternalIds:    externalIDs,
	})
}

func (s *Service) ListShows(ctx context.Context) ([]sqlc.Show, error) {
	return s.q.ListShows(ctx)
}

func (s *Service) GetShowByID(ctx context.Context, showID string) (sqlc.Show, error) {
	return s.q.GetShowByID(ctx, showID)
}

func (s *Service) UpdateShow(ctx context.Context, showID string, req updateShowRequest) (sqlc.Show, error) {
	externalIDs, err := marshalExternalID(req.ExternalID)
	if err != nil {
		return sqlc.Show{}, err
	}

	return s.q.UpdateShow(ctx, sqlc.UpdateShowParams{
		InternalShowID: showID,
		TitlePreferred: req.TitlePreferred,
		TitleOriginal:  req.TitleOriginal,
		AltTitles:      req.AltTitles,
		Type:           req.Type,
		Status:         req.Status,
		Synopsis:       req.Synopsis,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		PosterUrl:      req.PosterUrl,
		BannerUrl:      req.BannerUrl,
		SeasonCount:    req.SeasonCount,
		EpisodeCount:   req.EpisodeCount,
		ExternalIds:    externalIDs,
	})
}

func (s *Service) DeleteShow(ctx context.Context, showID string) error {
	_, err := s.q.DeleteShow(ctx, showID)
	return err
}

func (s *Service) ListWorkerData(ctx context.Context) ([]workerDataResponse, error) {
	shows, err := s.q.ListShows(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]workerDataResponse, 0, len(shows))
	for _, item := range shows {
		externalID, err := unmarshalExternalID(item.ExternalIds)
		if err != nil {
			return nil, err
		}

		episodes, err := s.q.ListEpisodesByShowID(ctx, item.InternalShowID)
		if err != nil {
			return nil, err
		}

		mappedEpisodes := make([]workerEpisode, 0, len(episodes))
		for _, ep := range episodes {
			mappedEpisodes = append(mappedEpisodes, workerEpisode{
				EpisodeNumber: ep.EpisodeNumber,
				AirDate:       ep.AirDate,
			})
		}

		response = append(response, workerDataResponse{
			InternalShowID: item.InternalShowID,
			Show: workerShowResponse{
				ExternalID: externalID,
				AltTitles:  normalizeAltTitles(item.AltTitles),
			},
			Episodes: mappedEpisodes,
		})
	}

	return response, nil
}
