package show

import (
	"errors"
	"time"

	"github.com/keithics/devops-dashboard/api/internal/db/sqlc"
)

const (
	ctxCreateShowRequestKey = "show.create.request"
	ctxUpdateShowRequestKey = "show.update.request"
	ctxShowIDKey            = "show.id"
)

var errInvalidShowID = errors.New("invalid show id")

type Handler struct {
	svc *Service
}

type Service struct {
	q *sqlc.Queries
}

type Show struct {
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

type createShowRequest = Show

type updateShowRequest = Show

type showResponse struct {
	InternalShowID string `json:"internalShowId"`
	Show
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type workerDataResponse struct {
	InternalShowID string             `json:"internalShowId"`
	Show           workerShowResponse `json:"show"`
	Episodes       []workerEpisode    `json:"episodes"`
}

type workerShowResponse struct {
	ExternalID string   `json:"externalId"`
	AltTitles  []string `json:"altTitles"`
}

type workerEpisode struct {
	EpisodeNumber int64   `json:"episodeNumber"`
	AirDate       *string `json:"airDate"`
}
