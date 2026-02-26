package metadata

import (
	jobshow "github.com/keithics/devops-dashboard/api/internal/job/show"
	worker "github.com/keithics/devops-dashboard/api/internal/metadata/provider"
)

const (
	ctxProviderTypeKey = "metadata.provider.type"
	ctxQueryKey        = "metadata.search.query"
	ctxExternalIDKey   = "metadata.external.id"
	ctxSearchOptsKey   = "metadata.search.opts"
	ctxEpisodesOptsKey = "metadata.episodes.opts"
	ctxAddShowRequest  = "metadata.add.show.request"
)

type Handler struct {
	svc *Service
}

type Service struct {
	worker  *worker.Service
	jobShow *jobshow.Service
}

type SearchHitResponse struct {
	ExternalID     string   `json:"externalId,omitempty"`
	TitlePreferred string   `json:"titlePreferred"`
	TitleOriginal  *string  `json:"titleOriginal,omitempty"`
	AltTitles      []string `json:"altTitles"`
	Type           string   `json:"type"`
	Status         string   `json:"status"`
	Synopsis       *string  `json:"synopsis,omitempty"`
	StartDate      *string  `json:"startDate,omitempty"`
	EndDate        *string  `json:"endDate,omitempty"`
	PosterUrl      *string  `json:"posterUrl,omitempty"`
	BannerUrl      *string  `json:"bannerUrl,omitempty"`
	SeasonCount    *int64   `json:"seasonCount,omitempty"`
	EpisodeCount   *int64   `json:"episodeCount,omitempty"`
}

type AddShowRequest = SearchHitResponse

type ShowResponse struct {
	ExternalID     string   `json:"externalId,omitempty"`
	TitlePreferred string   `json:"titlePreferred"`
	TitleOriginal  *string  `json:"titleOriginal,omitempty"`
	AltTitles      []string `json:"altTitles"`
	Type           string   `json:"type"`
	Status         string   `json:"status"`
	Synopsis       *string  `json:"synopsis,omitempty"`
	StartDate      *string  `json:"startDate,omitempty"`
	EndDate        *string  `json:"endDate,omitempty"`
	PosterUrl      *string  `json:"posterUrl,omitempty"`
	BannerUrl      *string  `json:"bannerUrl,omitempty"`
	SeasonCount    *int64   `json:"seasonCount,omitempty"`
	EpisodeCount   *int64   `json:"episodeCount,omitempty"`
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

type AddShowResponse struct {
	InternalShowID    string `json:"internalShowId"`
	InternalJobShowID string `json:"internalJobShowId"`
	Status            string `json:"status"`
	RetryCount        int32  `json:"retryCount"`
}
