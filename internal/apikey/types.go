package apikey

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	svc *Service
}

type Service struct {
	pool *pgxpool.Pool
}

type createAPIKeyRequest struct {
	Name string `json:"name"`
}

type createAPIKeyResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Key       string    `json:"key"`
	Last4     string    `json:"last4"`
	CreatedAt time.Time `json:"createdAt"`
}
