package show

import "github.com/jackc/pgx/v5/pgxpool"

type Service struct {
	pool *pgxpool.Pool
}

type EnqueueFromSearchResultParams struct {
	ExternalID     string
	TitlePreferred string
	TitleOriginal  *string
	AltTitles      []string
	Type           string
	Status         string
	Synopsis       *string
	StartDate      *string
	EndDate        *string
	PosterURL      *string
	BannerURL      *string
	SeasonCount    *int64
	EpisodeCount   *int64
}

type EnqueueFromSearchResultResult struct {
	InternalShowID    string `json:"internalShowId"`
	InternalJobShowID string `json:"internalJobShowId"`
	Status            string `json:"status"`
	RetryCount        int32  `json:"retryCount"`
}
