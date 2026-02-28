package metadata

import (
	worker "github.com/keithics/devops-dashboard/api/internal/metadata/provider"
	showmodel "github.com/keithics/devops-dashboard/api/internal/show"
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

type SearchHitResponse = showmodel.Show

type ShowResponse = showmodel.Show

type EpisodeResponse struct {
	Provider       string  `json:"provider"`
	ExternalID     string  `json:"externalId"`
	SeasonNumber   int64   `json:"seasonNumber"`
	EpisodeNumber  int64   `json:"episodeNumber"`
	Title          string  `json:"title"`
	AirDate        *string `json:"airDate,omitempty"`
	RuntimeMinutes *int64  `json:"runtimeMinutes,omitempty"`
}
