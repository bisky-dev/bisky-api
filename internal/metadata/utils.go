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
		return "", errors.New("type must be one of anidb|tvdb")
	}
}

func validateQuery(query string) error {
	if query == "" {
		return errors.New("query is required")
	}
	return nil
}

func validateExternalID(externalID string) error {
	if externalID == "" {
		return errors.New("externalId is required")
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

func getProviderAndExternalID(c *gin.Context) (worker.ProviderName, string, bool) {
	provider, ok := httpx.AbortIfMissingContext[worker.ProviderName](c, ctxProviderTypeKey)
	if !ok {
		return "", "", false
	}
	externalID, ok := httpx.AbortIfMissingContext[string](c, ctxExternalIDKey)
	if !ok {
		return "", "", false
	}
	return provider, externalID, true
}

func getEpisodesInput(c *gin.Context) (worker.ProviderName, string, worker.ListEpisodesOpts, bool) {
	provider, externalID, ok := getProviderAndExternalID(c)
	if !ok {
		return "", "", worker.ListEpisodesOpts{}, false
	}
	opts, ok := httpx.AbortIfMissingContext[worker.ListEpisodesOpts](c, ctxEpisodesOptsKey)
	if !ok {
		return "", "", worker.ListEpisodesOpts{}, false
	}
	return provider, externalID, opts, true
}
