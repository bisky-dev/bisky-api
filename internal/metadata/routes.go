package metadata

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *Handler) {
	r.GET("/metadata/search", h.BindSearch(), h.Search)
	r.GET("/metadata/show/:externalId", h.BindExternalID(), h.GetShow)
	r.GET("/metadata/episodes/:externalId", h.BindExternalID(), h.BindEpisodesOpts(), h.ListEpisodes)
}
