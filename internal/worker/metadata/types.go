package metadata

import "context"

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
	Provider       ProviderName `json:"provider"`
	ExternalID     string       `json:"externalId"`
	TitlePreferred string       `json:"titlePreferred"`
	TitleOriginal  *string      `json:"titleOriginal,omitempty"`
	Type           string       `json:"type,omitempty"`
	Score          *float64     `json:"score,omitempty"`
	Description    *string      `json:"description,omitempty"`
	BannerURL      *string      `json:"bannerUrl,omitempty"`
}

type Show struct {
	Provider       ProviderName `json:"provider"`
	ExternalID     string       `json:"externalId"`
	TitlePreferred string       `json:"titlePreferred"`
	TitleOriginal  *string      `json:"titleOriginal,omitempty"`
	Synopsis       *string      `json:"synopsis,omitempty"`
	StartDate      *string      `json:"startDate,omitempty"`
	EndDate        *string      `json:"endDate,omitempty"`
	PosterURL      *string      `json:"posterUrl,omitempty"`
	BannerURL      *string      `json:"bannerUrl,omitempty"`
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
