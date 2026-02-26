package metadata

import (
	jobshow "github.com/keithics/devops-dashboard/api/internal/job/show"
	worker "github.com/keithics/devops-dashboard/api/internal/worker/metadata"
)

const (
	ctxProviderTypeKey = "metadata.provider.type"
	ctxQueryKey        = "metadata.search.query"
	ctxSearchOptsKey   = "metadata.search.opts"
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
	Provider       string   `json:"provider"`
	ExternalID     string   `json:"externalId"`
	TitlePreferred string   `json:"titlePreferred"`
	TitleOriginal  *string  `json:"titleOriginal,omitempty"`
	Type           string   `json:"type,omitempty"`
	Score          *float64 `json:"score,omitempty"`
	Description    *string  `json:"description,omitempty"`
	BannerURL      *string  `json:"bannerUrl,omitempty"`
}

type AddShowRequest struct {
	Provider       string   `json:"provider"`
	ExternalID     string   `json:"externalId"`
	TitlePreferred string   `json:"titlePreferred"`
	TitleOriginal  *string  `json:"titleOriginal,omitempty"`
	Type           *string  `json:"type,omitempty"`
	Score          *float64 `json:"score,omitempty"`
	Description    *string  `json:"description,omitempty"`
	BannerURL      *string  `json:"bannerUrl,omitempty"`
}

type AddShowResponse struct {
	InternalSearchResultID string `json:"internalSearchResultId"`
	InternalJobShowID      string `json:"internalJobShowId"`
	Status                 string `json:"status"`
	RetryCount             int32  `json:"retryCount"`
}
