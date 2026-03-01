package show

import (
	"context"

	"github.com/keithics/devops-dashboard/api/internal/db/sqlc"
	"github.com/keithics/devops-dashboard/api/internal/hooks"
	normalizeutil "github.com/keithics/devops-dashboard/api/internal/utils/normalize"
)

func NewHandler(q *sqlc.Queries) *Handler {
	return NewHandlerWithHooks(q, hooks.NoopDispatcher{})
}

func NewHandlerWithHooks(q *sqlc.Queries, dispatcher hooks.Dispatcher) *Handler {
	if dispatcher == nil {
		dispatcher = hooks.NoopDispatcher{}
	}

	service := NewService(q, dispatcher)

	return &Handler{
		svc: service,
	}
}

func NewService(q *sqlc.Queries, dispatcher hooks.Dispatcher) *Service {
	if dispatcher == nil {
		dispatcher = hooks.NoopDispatcher{}
	}
	return &Service{
		q:     q,
		hooks: dispatcher,
	}
}

func (h *Handler) Service() *Service {
	return h.svc
}

func (s *Service) CreateShow(ctx context.Context, req Show) (sqlc.Show, error) {
	createReq := createShowRequest(req)

	if err := s.hooks.DispatchPre(ctx, hooks.EventShowCreatePre, req); err != nil {
		return sqlc.Show{}, err
	}

	externalIDs, err := marshalExternalID(createReq.ExternalID)
	if err != nil {
		return sqlc.Show{}, err
	}

	created, err := s.q.CreateShow(ctx, sqlc.CreateShowParams{
		TitlePreferred: createReq.TitlePreferred,
		TitleOriginal:  createReq.TitleOriginal,
		AltTitles:      createReq.AltTitles,
		Type:           createReq.Type,
		Status:         createReq.Status,
		Synopsis:       createReq.Synopsis,
		StartDate:      createReq.StartDate,
		EndDate:        createReq.EndDate,
		PosterUrl:      createReq.PosterUrl,
		BannerUrl:      createReq.BannerUrl,
		SeasonCount:    createReq.SeasonCount,
		EpisodeCount:   createReq.EpisodeCount,
		ExternalIds:    externalIDs,
	})
	if err != nil {
		return sqlc.Show{}, err
	}

	s.hooks.DispatchPost(ctx, hooks.EventShowCreatePost, created)
	return created, nil
}

func (s *Service) ListShows(ctx context.Context) ([]sqlc.Show, error) {
	return s.q.ListShows(ctx)
}

func (s *Service) GetShowByID(ctx context.Context, showID string) (sqlc.Show, error) {
	return s.q.GetShowByID(ctx, showID)
}

func (s *Service) UpdateShow(ctx context.Context, showID string, req updateShowRequest) (sqlc.Show, error) {
	if err := s.hooks.DispatchPre(ctx, hooks.EventShowUpdatePre, map[string]any{
		"internalShowId": showID,
		"request":        req,
	}); err != nil {
		return sqlc.Show{}, err
	}

	externalIDs, err := marshalExternalID(req.ExternalID)
	if err != nil {
		return sqlc.Show{}, err
	}

	updated, err := s.q.UpdateShow(ctx, sqlc.UpdateShowParams{
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
	if err != nil {
		return sqlc.Show{}, err
	}

	s.hooks.DispatchPost(ctx, hooks.EventShowUpdatePost, updated)
	return updated, nil
}

func (s *Service) DeleteShow(ctx context.Context, showID string) error {
	if err := s.hooks.DispatchPre(ctx, hooks.EventShowDeletePre, map[string]any{
		"internalShowId": showID,
	}); err != nil {
		return err
	}

	_, err := s.q.DeleteShow(ctx, showID)
	if err != nil {
		return err
	}

	s.hooks.DispatchPost(ctx, hooks.EventShowDeletePost, map[string]any{
		"internalShowId": showID,
	})
	return nil
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
				AltTitles:  normalizeutil.Strings(item.AltTitles),
			},
			Episodes: mappedEpisodes,
		})
	}

	return response, nil
}
