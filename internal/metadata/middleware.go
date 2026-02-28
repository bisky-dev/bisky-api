package metadata

import (
	"github.com/gin-gonic/gin"
	"github.com/keithics/devops-dashboard/api/internal/httpx"
	worker "github.com/keithics/devops-dashboard/api/internal/metadata/provider"
	normalizeutil "github.com/keithics/devops-dashboard/api/internal/utils/normalize"
)

func (h *Handler) BindSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		provider, err := parseProvider(c.Query("type"))
		if httpx.AbortIfErr(c, err) {
			return
		}

		query := normalizeutil.String(c.Query("query"))
		if httpx.AbortIfErr(c, validateQuery(query)) {
			return
		}

		opts := worker.SearchOpts{
			Page:   httpx.ParsePositiveInt(c.Query("page"), 1),
			Limit:  httpx.ParsePositiveInt(c.Query("limit"), 10),
			Strict: httpx.ParseBool(c.Query("strict"), true),
		}

		c.Set(ctxProviderTypeKey, provider)
		c.Set(ctxQueryKey, query)
		c.Set(ctxSearchOptsKey, opts)
		c.Next()
	}
}

func (h *Handler) BindAddShow() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AddShowRequest
		if httpx.AbortIfErr(c, c.ShouldBindJSON(&req)) {
			return
		}

		normalizeAddShowRequest(&req)
		if httpx.AbortIfErr(c, validateAddShowRequest(req)) {
			return
		}

		c.Set(ctxAddShowRequest, req)
		c.Next()
	}
}

func (h *Handler) BindExternalID() gin.HandlerFunc {
	return func(c *gin.Context) {
		provider, err := parseProvider(c.Query("type"))
		if httpx.AbortIfErr(c, err) {
			return
		}

		externalID := normalizeutil.String(c.Param("externalId"))
		if httpx.AbortIfErr(c, validateExternalID(externalID)) {
			return
		}

		c.Set(ctxProviderTypeKey, provider)
		c.Set(ctxExternalIDKey, externalID)
		c.Next()
	}
}

func (h *Handler) BindEpisodesOpts() gin.HandlerFunc {
	return func(c *gin.Context) {
		opts := worker.ListEpisodesOpts{
			Page:  httpx.ParsePositiveInt(c.Query("page"), 1),
			Limit: httpx.ParsePositiveInt(c.Query("limit"), 25),
		}
		c.Set(ctxEpisodesOptsKey, opts)
		c.Next()
	}
}
