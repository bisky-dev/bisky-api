package episode

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *Handler) {
	r.GET("/episodes", h.ListEpisodes)
	r.GET("/episodes/:internalEpisodeId", h.BindEpisodeID(), h.GetEpisode)
	r.POST("/episodes", h.BindCreateEpisode(), h.CreateEpisode)
	r.PUT("/episodes/:internalEpisodeId", h.BindEpisodeID(), h.BindUpdateEpisode(), h.UpdateEpisode)
	r.DELETE("/episodes/:internalEpisodeId", h.BindEpisodeID(), h.DeleteEpisode)
}
