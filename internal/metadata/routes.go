package metadata

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *Handler) {
	r.GET("/metadata/search", h.BindSearch(), h.Search)
	r.GET("/metadata/discover", h.BindDiscover(), h.Discover)
	r.GET("/metadata/show/:externalId", h.BindExternalID(), h.GetShow)
	r.POST("/metadata/show/:externalId", h.BindExternalID(), h.AddShow)
	r.GET("/metadata/episodes/:externalId", h.BindExternalID(), h.BindEpisodesOpts(), h.ListEpisodes)
}
