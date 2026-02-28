package hooksettings

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/keithics/devops-dashboard/api/internal/hooks"
)

type Handler struct {
	svc *Service
}

type Service struct {
	pool *pgxpool.Pool
}

type upsertHookItem struct {
	Event hooks.Event `json:"event"`
	URL   string      `json:"url"`
}

type upsertHooksRequest struct {
	Hooks []upsertHookItem `json:"hooks"`
}
