package metadata

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/keithics/devops-dashboard/api/internal/httperr"
	"github.com/keithics/devops-dashboard/api/internal/httpx"
	worker "github.com/keithics/devops-dashboard/api/internal/worker/metadata"
)

func parseProvider(raw string) (worker.ProviderName, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		return worker.ProviderAniDB, nil
	}

	switch value {
	case string(worker.ProviderAniDB):
		return worker.ProviderAniDB, nil
	case string(worker.ProviderTVDB):
		return worker.ProviderTVDB, nil
	case string(worker.ProviderAniList):
		return worker.ProviderAniList, nil
	default:
		return "", errors.New("type must be one of anidb|anilist|tvdb")
	}
}

func validateQuery(query string) error {
	if query == "" {
		return errors.New("query is required")
	}
	return nil
}

func abortProviderErr(c *gin.Context, internalMessage string, err error) {
	message := strings.ToLower(err.Error())
	if strings.Contains(message, "not implemented") || strings.Contains(message, "not supported") {
		httperr.Abort(c, httperr.BadRequest(err.Error()))
		return
	}
	httperr.Abort(c, httperr.Internal(internalMessage).WithCause(err))
}

func getSearchInput(c *gin.Context) (worker.ProviderName, string, worker.SearchOpts, bool) {
	provider, ok := httpx.AbortIfMissingContext[worker.ProviderName](c, ctxProviderTypeKey)
	if !ok {
		return "", "", worker.SearchOpts{}, false
	}
	query, ok := httpx.AbortIfMissingContext[string](c, ctxQueryKey)
	if !ok {
		return "", "", worker.SearchOpts{}, false
	}
	opts, ok := httpx.AbortIfMissingContext[worker.SearchOpts](c, ctxSearchOptsKey)
	if !ok {
		return "", "", worker.SearchOpts{}, false
	}
	return provider, query, opts, true
}

func normalizeAddShowRequest(req *AddShowRequest) {
	req.ExternalID = strings.TrimSpace(req.ExternalID)
	req.TitlePreferred = strings.TrimSpace(req.TitlePreferred)
	req.TitleOriginal = httpx.TrimmedOrNil(req.TitleOriginal)
	req.Type = strings.ToLower(strings.TrimSpace(req.Type))
	req.Status = strings.ToLower(strings.TrimSpace(req.Status))
	req.Synopsis = httpx.TrimmedOrNil(req.Synopsis)
	req.StartDate = httpx.TrimmedOrNil(req.StartDate)
	req.EndDate = httpx.TrimmedOrNil(req.EndDate)
	req.PosterUrl = httpx.TrimmedOrNil(req.PosterUrl)
	req.BannerUrl = httpx.TrimmedOrNil(req.BannerUrl)
	req.AltTitles = normalizeAltTitles(req.AltTitles)
}

func validateAddShowRequest(req AddShowRequest) error {
	if err := httpx.ValidateVar(req.ExternalID, "required,max=128", "externalId is invalid"); err != nil {
		return err
	}
	if err := httpx.ValidateVar(req.TitlePreferred, "required,max=500", "titlePreferred is invalid"); err != nil {
		return err
	}
	if err := httpx.ValidateVar(req.Type, "required,oneof=anime tv movie ova special", "type is invalid"); err != nil {
		return err
	}
	if err := httpx.ValidateVar(req.Status, "required,oneof=ongoing finished", "status is invalid"); err != nil {
		return err
	}
	if err := httpx.ValidateOptionalDate(req.StartDate, "startDate is invalid"); err != nil {
		return err
	}
	if err := httpx.ValidateOptionalDate(req.EndDate, "endDate is invalid"); err != nil {
		return err
	}
	if err := validateOptionalInt64(req.SeasonCount, "gte=0", "seasonCount is invalid"); err != nil {
		return err
	}
	if err := validateOptionalInt64(req.EpisodeCount, "gte=0", "episodeCount is invalid"); err != nil {
		return err
	}
	return nil
}

func normalizeAltTitles(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}

	normalized := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		normalized = append(normalized, trimmed)
	}
	return normalized
}

func validateOptionalInt64(value *int64, rule string, message string) error {
	if value == nil {
		return nil
	}
	return httpx.ValidateVar(*value, rule, message)
}
