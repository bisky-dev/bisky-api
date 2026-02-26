package episode

import (
	"time"

	"github.com/keithics/devops-dashboard/api/internal/db/sqlc"
)

const (
	ctxCreateEpisodeRequestKey = "episode.create.request"
	ctxUpdateEpisodeRequestKey = "episode.update.request"
	ctxEpisodeIDKey            = "episode.id"
)

type Handler struct {
	svc *Service
}

type Service struct {
	q *sqlc.Queries
}

type externalIDs struct {
	Anilist *int64 `json:"anilist,omitempty"`
	Tvdb    *int64 `json:"tvdb,omitempty"`
}

type createEpisodeRequest struct {
	ShowID         string      `json:"showId"`
	SeasonNumber   int64       `json:"seasonNumber"`
	EpisodeNumber  int64       `json:"episodeNumber"`
	Title          string      `json:"title"`
	AirDate        *string     `json:"airDate"`
	RuntimeMinutes *int64      `json:"runtimeMinutes"`
	ExternalIDs    externalIDs `json:"externalIds"`
}

type updateEpisodeRequest struct {
	ShowID         string      `json:"showId"`
	SeasonNumber   int64       `json:"seasonNumber"`
	EpisodeNumber  int64       `json:"episodeNumber"`
	Title          string      `json:"title"`
	AirDate        *string     `json:"airDate"`
	RuntimeMinutes *int64      `json:"runtimeMinutes"`
	ExternalIDs    externalIDs `json:"externalIds"`
}

type episodeResponse struct {
	InternalEpisodeID string      `json:"internalEpisodeId"`
	ShowID            string      `json:"showId"`
	SeasonNumber      int64       `json:"seasonNumber"`
	EpisodeNumber     int64       `json:"episodeNumber"`
	Title             string      `json:"title"`
	AirDate           *string     `json:"airDate,omitempty"`
	RuntimeMinutes    *int64      `json:"runtimeMinutes,omitempty"`
	ExternalIDs       externalIDs `json:"externalIds"`
	CreatedAt         time.Time   `json:"createdAt"`
	UpdatedAt         time.Time   `json:"updatedAt"`
}
