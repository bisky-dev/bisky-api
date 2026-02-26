package show

import "github.com/jackc/pgx/v5/pgxpool"

type Service struct {
	pool *pgxpool.Pool
}

type EnqueueFromSearchResultParams struct {
	Provider       string
	ExternalID     string
	TitlePreferred string
	TitleOriginal  *string
	Type           *string
	Score          *float64
	Description    *string
	BannerURL      *string
}

type EnqueueFromSearchResultResult struct {
	InternalSearchResultID string `json:"internalSearchResultId"`
	InternalJobShowID      string `json:"internalJobShowId"`
	Status                 string `json:"status"`
	RetryCount             int32  `json:"retryCount"`
}
