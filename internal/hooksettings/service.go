package hooksettings

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/keithics/devops-dashboard/api/internal/hooks"
)

func NewHandler(pool *pgxpool.Pool) (*Handler, error) {
	if err := hooks.EnsureStore(context.Background(), pool); err != nil {
		return nil, err
	}
	return &Handler{
		svc: &Service{pool: pool},
	}, nil
}

func (s *Service) List(ctx context.Context) ([]hooks.Config, error) {
	return hooks.ListConfigs(ctx, s.pool)
}

func (s *Service) ListKeys() []hooks.Event {
	return hooks.AllEvents()
}

func (s *Service) Upsert(ctx context.Context, req upsertHooksRequest) ([]hooks.Config, error) {
	if len(req.Hooks) == 0 {
		return nil, errors.New("hooks payload is required")
	}

	updated := make([]hooks.Config, 0, len(req.Hooks))
	for _, item := range req.Hooks {
		config, err := hooks.UpsertConfig(ctx, s.pool, item.Event, item.URL)
		if err != nil {
			return nil, err
		}
		updated = append(updated, config)
	}
	return updated, nil
}
