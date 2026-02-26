package episode

import (
	"context"

	"github.com/keithics/devops-dashboard/api/internal/db/sqlc"
)

func NewHandler(q *sqlc.Queries) *Handler {
	return &Handler{svc: &Service{q: q}}
}

func (s *Service) CreateEpisode(ctx context.Context, req createEpisodeRequest) (sqlc.Episode, error) {
	externalIDs, err := marshalExternalIDs(req.ExternalIDs)
	if err != nil {
		return sqlc.Episode{}, err
	}

	return s.q.CreateEpisode(ctx, sqlc.CreateEpisodeParams{
		ShowID:         req.ShowID,
		SeasonNumber:   req.SeasonNumber,
		EpisodeNumber:  req.EpisodeNumber,
		Title:          req.Title,
		AirDate:        req.AirDate,
		RuntimeMinutes: req.RuntimeMinutes,
		ExternalIds:    externalIDs,
	})
}

func (s *Service) ListEpisodes(ctx context.Context) ([]sqlc.Episode, error) {
	return s.q.ListEpisodes(ctx)
}

func (s *Service) GetEpisodeByID(ctx context.Context, episodeID string) (sqlc.Episode, error) {
	return s.q.GetEpisodeByID(ctx, episodeID)
}

func (s *Service) UpdateEpisode(ctx context.Context, episodeID string, req updateEpisodeRequest) (sqlc.Episode, error) {
	externalIDs, err := marshalExternalIDs(req.ExternalIDs)
	if err != nil {
		return sqlc.Episode{}, err
	}

	return s.q.UpdateEpisode(ctx, sqlc.UpdateEpisodeParams{
		InternalEpisodeID: episodeID,
		ShowID:            req.ShowID,
		SeasonNumber:      req.SeasonNumber,
		EpisodeNumber:     req.EpisodeNumber,
		Title:             req.Title,
		AirDate:           req.AirDate,
		RuntimeMinutes:    req.RuntimeMinutes,
		ExternalIds:       externalIDs,
	})
}

func (s *Service) DeleteEpisode(ctx context.Context, episodeID string) error {
	_, err := s.q.DeleteEpisode(ctx, episodeID)
	return err
}
