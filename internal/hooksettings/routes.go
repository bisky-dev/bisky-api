package hooksettings

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *Handler) {
	r.GET("/settings/hooks", h.ListHooks)
	r.GET("/settings/hooks/keys", h.ListHookKeys)
	r.PUT("/settings/hooks", h.UpsertHooks)
}
