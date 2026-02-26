package metadata

import (
	worker "github.com/keithics/devops-dashboard/api/internal/worker/metadata"
)

const (
	ctxProviderTypeKey = "metadata.provider.type"
	ctxQueryKey        = "metadata.search.query"
	ctxExternalIDKey   = "metadata.external.id"
	ctxSearchOptsKey   = "metadata.search.opts"
	ctxEpisodesOptsKey = "metadata.episodes.opts"
)

type Handler struct {
	svc *Service
}

type Service struct {
	worker *worker.Service
}

type SearchHitResponse struct {
	Provider       string   `json:"provider"`
	ExternalID     string   `json:"externalId"`
	TitlePreferred string   `json:"titlePreferred"`
	TitleOriginal  *string  `json:"titleOriginal,omitempty"`
	Type           string   `json:"type,omitempty"`
	Score          *float64 `json:"score,omitempty"`
	Description    *string  `json:"description,omitempty"`
	BannerURL      *string  `json:"bannerUrl,omitempty"`
}

type ShowResponse struct {
	Provider       string  `json:"provider"`
	ExternalID     string  `json:"externalId"`
	TitlePreferred string  `json:"titlePreferred"`
	TitleOriginal  *string `json:"titleOriginal,omitempty"`
	Synopsis       *string `json:"synopsis,omitempty"`
	StartDate      *string `json:"startDate,omitempty"`
	EndDate        *string `json:"endDate,omitempty"`
	PosterURL      *string `json:"posterUrl,omitempty"`
	BannerURL      *string `json:"bannerUrl,omitempty"`
}

type EpisodeResponse struct {
	Provider       string  `json:"provider"`
	ExternalID     string  `json:"externalId"`
	SeasonNumber   int64   `json:"seasonNumber"`
	EpisodeNumber  int64   `json:"episodeNumber"`
	Title          string  `json:"title"`
	AirDate        *string `json:"airDate,omitempty"`
	RuntimeMinutes *int64  `json:"runtimeMinutes,omitempty"`
}
