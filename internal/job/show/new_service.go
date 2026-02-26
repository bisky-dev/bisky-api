package show

import "github.com/jackc/pgx/v5/pgxpool"

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}
