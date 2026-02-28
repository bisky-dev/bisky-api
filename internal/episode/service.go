package episode

import (
	"context"

	"github.com/keithics/devops-dashboard/api/internal/db/sqlc"
	"github.com/keithics/devops-dashboard/api/internal/hooks"
)

func NewHandler(q *sqlc.Queries) *Handler {
	return NewHandlerWithHooks(q, hooks.NoopDispatcher{})
}

func NewHandlerWithHooks(q *sqlc.Queries, dispatcher hooks.Dispatcher) *Handler {
	if dispatcher == nil {
		dispatcher = hooks.NoopDispatcher{}
	}
	return &Handler{
		svc: &Service{
			q:     q,
			hooks: dispatcher,
		},
	}
}

func (s *Service) CreateEpisode(ctx context.Context, req createEpisodeRequest) (sqlc.Episode, error) {
	if err := s.hooks.DispatchPre(ctx, hooks.EventEpisodeCreatePre, req); err != nil {
		return sqlc.Episode{}, err
	}

	externalIDs, err := marshalExternalIDs(req.ExternalIDs)
	if err != nil {
		return sqlc.Episode{}, err
	}

	created, err := s.q.CreateEpisode(ctx, sqlc.CreateEpisodeParams{
		ShowID:         req.ShowID,
		SeasonNumber:   req.SeasonNumber,
		EpisodeNumber:  req.EpisodeNumber,
		Title:          req.Title,
		AirDate:        req.AirDate,
		RuntimeMinutes: req.RuntimeMinutes,
		ExternalIds:    externalIDs,
	})
	if err != nil {
		return sqlc.Episode{}, err
	}

	s.hooks.DispatchPost(ctx, hooks.EventEpisodeCreatePost, created)
	return created, nil
}

func (s *Service) ListEpisodes(ctx context.Context) ([]sqlc.Episode, error) {
	return s.q.ListEpisodes(ctx)
}

func (s *Service) GetEpisodeByID(ctx context.Context, episodeID string) (sqlc.Episode, error) {
	return s.q.GetEpisodeByID(ctx, episodeID)
}

func (s *Service) UpdateEpisode(ctx context.Context, episodeID string, req updateEpisodeRequest) (sqlc.Episode, error) {
	if err := s.hooks.DispatchPre(ctx, hooks.EventEpisodeUpdatePre, map[string]any{
		"internalEpisodeId": episodeID,
		"request":           req,
	}); err != nil {
		return sqlc.Episode{}, err
	}

	externalIDs, err := marshalExternalIDs(req.ExternalIDs)
	if err != nil {
		return sqlc.Episode{}, err
	}

	updated, err := s.q.UpdateEpisode(ctx, sqlc.UpdateEpisodeParams{
		InternalEpisodeID: episodeID,
		ShowID:            req.ShowID,
		SeasonNumber:      req.SeasonNumber,
		EpisodeNumber:     req.EpisodeNumber,
		Title:             req.Title,
		AirDate:           req.AirDate,
		RuntimeMinutes:    req.RuntimeMinutes,
		ExternalIds:       externalIDs,
	})
	if err != nil {
		return sqlc.Episode{}, err
	}

	s.hooks.DispatchPost(ctx, hooks.EventEpisodeUpdatePost, updated)
	return updated, nil
}

func (s *Service) DeleteEpisode(ctx context.Context, episodeID string) error {
	if err := s.hooks.DispatchPre(ctx, hooks.EventEpisodeDeletePre, map[string]any{
		"internalEpisodeId": episodeID,
	}); err != nil {
		return err
	}

	_, err := s.q.DeleteEpisode(ctx, episodeID)
	if err != nil {
		return err
	}

	s.hooks.DispatchPost(ctx, hooks.EventEpisodeDeletePost, map[string]any{
		"internalEpisodeId": episodeID,
	})
	return nil
}
