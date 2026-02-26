package show

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *Handler) {
	r.GET("/shows", h.ListShows)
	r.GET("/shows/:internalShowId", h.BindShowID(), h.GetShow)
	r.POST("/shows", h.BindCreateShow(), h.CreateShow)
	r.PUT("/shows/:internalShowId", h.BindShowID(), h.BindUpdateShow(), h.UpdateShow)
	r.DELETE("/shows/:internalShowId", h.BindShowID(), h.DeleteShow)
}
