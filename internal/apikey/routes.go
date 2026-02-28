package apikey

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *Handler) {
	r.POST("/api-keys", h.Create)
	r.POST("/api-keys/validate", h.Validate)
	r.DELETE("/api-keys/:id", h.Delete)
}
