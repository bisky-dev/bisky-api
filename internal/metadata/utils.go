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
	req.Provider = strings.ToLower(strings.TrimSpace(req.Provider))
	req.ExternalID = strings.TrimSpace(req.ExternalID)
	req.TitlePreferred = strings.TrimSpace(req.TitlePreferred)
	req.TitleOriginal = httpx.TrimmedOrNil(req.TitleOriginal)
	req.Type = httpx.TrimmedOrNil(req.Type)
	req.Description = httpx.TrimmedOrNil(req.Description)
	req.BannerURL = httpx.TrimmedOrNil(req.BannerURL)
}

func validateAddShowRequest(req AddShowRequest) error {
	if err := httpx.ValidateVar(req.Provider, "required,oneof=anidb anilist tvdb", "provider is invalid"); err != nil {
		return err
	}
	if err := httpx.ValidateVar(req.ExternalID, "required,max=128", "externalId is invalid"); err != nil {
		return err
	}
	if err := httpx.ValidateVar(req.TitlePreferred, "required,max=500", "titlePreferred is invalid"); err != nil {
		return err
	}
	if req.Type != nil {
		if err := httpx.ValidateVar(*req.Type, "oneof=anime tv movie ova special", "type is invalid"); err != nil {
			return err
		}
	}
	if req.Score != nil {
		if err := httpx.ValidateVar(*req.Score, "gte=0", "score is invalid"); err != nil {
			return err
		}
	}
	if req.BannerURL != nil {
		if err := httpx.ValidateVar(*req.BannerURL, "url", "bannerUrl is invalid"); err != nil {
			return err
		}
	}
	return nil
}
