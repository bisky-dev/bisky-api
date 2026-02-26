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

type externalIDs struct {
	Anilist *int64 `json:"anilist,omitempty"`
	Tvdb    *int64 `json:"tvdb,omitempty"`
}

type Show struct {
	TitlePreferred string      `json:"titlePreferred"`
	TitleOriginal  *string     `json:"titleOriginal,omitempty"`
	AltTitles      []string    `json:"altTitles"`
	Type           string      `json:"type"`
	Status         string      `json:"status"`
	Synopsis       *string     `json:"synopsis,omitempty"`
	StartDate      *string     `json:"startDate,omitempty"`
	EndDate        *string     `json:"endDate,omitempty"`
	PosterUrl      *string     `json:"posterUrl,omitempty"`
	BannerUrl      *string     `json:"bannerUrl,omitempty"`
	SeasonCount    *int64      `json:"seasonCount,omitempty"`
	EpisodeCount   *int64      `json:"episodeCount,omitempty"`
	ExternalIDs    externalIDs `json:"externalIds"`
}

type createShowRequest = Show

type updateShowRequest = Show

type showResponse struct {
	InternalShowID string `json:"internalShowId"`
	Show
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
