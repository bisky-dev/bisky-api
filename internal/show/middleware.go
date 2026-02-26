package show

import (
	"github.com/gin-gonic/gin"
	"github.com/keithics/devops-dashboard/api/internal/httpx"
)

func (h *Handler) BindCreateShow() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createShowRequest
		if httpx.AbortIfErr(c, c.ShouldBindJSON(&req)) {
			return
		}

		normalizeCreateShowRequest(&req)
		if httpx.AbortIfErr(c, validateCreateShowRequest(req)) {
			return
		}

		c.Set(ctxCreateShowRequestKey, req)
		c.Next()
	}
}

func (h *Handler) BindUpdateShow() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req updateShowRequest
		if httpx.AbortIfErr(c, c.ShouldBindJSON(&req)) {
			return
		}

		normalizeUpdateShowRequest(&req)
		if httpx.AbortIfErr(c, validateUpdateShowRequest(req)) {
			return
		}

		c.Set(ctxUpdateShowRequestKey, req)
		c.Next()
	}
}

func (h *Handler) BindShowID() gin.HandlerFunc {
	return func(c *gin.Context) {
		showID := c.Param("internalShowId")
		if httpx.AbortIfErr(c, validateShowID(showID)) {
			return
		}
		c.Set(ctxShowIDKey, showID)
		c.Next()
	}
}
