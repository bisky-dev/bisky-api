package metadata

import (
	"time"

	worker "github.com/keithics/devops-dashboard/api/internal/metadata/provider"
	"github.com/keithics/devops-dashboard/api/internal/show"
)

const (
	ctxProviderTypeKey = "metadata.provider.type"
	ctxQueryKey        = "metadata.search.query"
	ctxExternalIDKey   = "metadata.external.id"
	ctxSearchOptsKey   = "metadata.search.opts"
	ctxDiscoverOptsKey = "metadata.discover.opts"
	ctxEpisodesOptsKey = "metadata.episodes.opts"
)

type Handler struct {
	svc *Service
}

type Service struct {
	worker  *worker.Service
	showSvc *show.Service
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

type DiscoverResponse struct {
	Trending        []ShowResponse `json:"trending"`
	Popular         []ShowResponse `json:"popular"`
	TopRated        []ShowResponse `json:"topRated"`
	Upcoming        []ShowResponse `json:"upcoming"`
	CurrentlyAiring []ShowResponse `json:"currentlyAiring"`
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
	InternalShowID string `json:"internalShowId"`
	ShowResponse
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
