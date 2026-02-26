package metadata

import (
	"context"

	showmodel "github.com/keithics/devops-dashboard/api/internal/show"
)

type ProviderName string

const (
	ProviderAniDB   ProviderName = "anidb"
	ProviderAniList ProviderName = "anilist"
	ProviderTVDB    ProviderName = "tvdb"
)

type SearchOpts struct {
	Page  int
	Limit int
}

type ListEpisodesOpts struct {
	Page         int
	Limit        int
	SeasonNumber *int64
}

type SearchHit struct {
	showmodel.Show
}

type Show struct {
	showmodel.Show
}

type Episode struct {
	Provider       ProviderName `json:"provider"`
	ExternalID     string       `json:"externalId"`
	SeasonNumber   int64        `json:"seasonNumber"`
	EpisodeNumber  int64        `json:"episodeNumber"`
	Title          string       `json:"title"`
	AirDate        *string      `json:"airDate,omitempty"`
	RuntimeMinutes *int64       `json:"runtimeMinutes,omitempty"`
}

type Provider interface {
	Search(ctx context.Context, query string, opts SearchOpts) ([]SearchHit, error)
	GetShow(ctx context.Context, externalID string) (Show, error)
	ListEpisodes(ctx context.Context, externalID string, opts ListEpisodesOpts) ([]Episode, error)
}

type Registry struct {
	providers map[ProviderName]Provider
}

type Service struct {
	registry *Registry
}
