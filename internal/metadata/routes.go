package metadata

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *Handler) {
	r.GET("/metadata/search", h.BindSearch(), h.Search)
	r.POST("/metadata/show", h.BindAddShow(), h.AddShow)
}
