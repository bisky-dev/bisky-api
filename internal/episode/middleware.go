package episode

import (
	"github.com/gin-gonic/gin"
	"github.com/keithics/devops-dashboard/api/internal/httpx"
)

func (h *Handler) BindCreateEpisode() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createEpisodeRequest
		if httpx.AbortIfErr(c, c.ShouldBindJSON(&req)) {
			return
		}
		normalizeCreateEpisodeRequest(&req)
		if httpx.AbortIfErr(c, validateCreateEpisodeRequest(req)) {
			return
		}
		c.Set(ctxCreateEpisodeRequestKey, req)
		c.Next()
	}
}

func (h *Handler) BindUpdateEpisode() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req updateEpisodeRequest
		if httpx.AbortIfErr(c, c.ShouldBindJSON(&req)) {
			return
		}
		normalizeUpdateEpisodeRequest(&req)
		if httpx.AbortIfErr(c, validateUpdateEpisodeRequest(req)) {
			return
		}
		c.Set(ctxUpdateEpisodeRequestKey, req)
		c.Next()
	}
}

func (h *Handler) BindEpisodeID() gin.HandlerFunc {
	return func(c *gin.Context) {
		episodeID := c.Param("internalEpisodeId")
		if httpx.AbortIfErr(c, validateEpisodeID(episodeID)) {
			return
		}
		c.Set(ctxEpisodeIDKey, episodeID)
		c.Next()
	}
}
