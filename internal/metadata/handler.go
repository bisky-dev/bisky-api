package metadata

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keithics/devops-dashboard/api/internal/httperr"
	"github.com/keithics/devops-dashboard/api/internal/httpx"
)

// Search godoc
//
//	@Summary		Metadata search
//	@Description	Search metadata by provider type (anidb default, tvdb optional)
//	@Tags			metadata
//	@Produce		json
//	@Param			query	query		string	true	"Search query"
//	@Param			type	query		string	false	"Provider type: anidb|tvdb (default anidb)"
//	@Param			page	query		int		false	"Page"
//	@Param			limit	query		int		false	"Limit"
//	@Success		200		{array}		SearchHitResponse
//	@Failure		400		{object}	httperr.APIErrorResponse
//	@Failure		500		{object}	httperr.APIErrorResponse
//	@Router			/metadata/search [get]
func (h *Handler) Search(c *gin.Context) {
	provider, query, opts, ok := getSearchInput(c)
	if !ok {
		return
	}

	items, err := h.svc.Search(c.Request.Context(), provider, query, opts)
	if err != nil {
		abortProviderErr(c, "failed to search metadata", err)
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetShow godoc
//
//	@Summary		Metadata get show
//	@Description	Get show metadata by provider external id
//	@Tags			metadata
//	@Produce		json
//	@Param			externalId	path		string	true	"Provider external id"
//	@Param			type		query		string	false	"Provider type: anidb|tvdb (default anidb)"
//	@Success		200			{object}	ShowResponse
//	@Failure		400			{object}	httperr.APIErrorResponse
//	@Failure		500			{object}	httperr.APIErrorResponse
//	@Router			/metadata/show/{externalId} [get]
func (h *Handler) GetShow(c *gin.Context) {
	provider, externalID, ok := getProviderAndExternalID(c)
	if !ok {
		return
	}

	item, err := h.svc.GetShow(c.Request.Context(), provider, externalID)
	if err != nil {
		abortProviderErr(c, "failed to get metadata show", err)
		return
	}

	c.JSON(http.StatusOK, item)
}

// ListEpisodes godoc
//
//	@Summary		Metadata list episodes
//	@Description	List episodes metadata by provider external id
//	@Tags			metadata
//	@Produce		json
//	@Param			externalId	path		string	true	"Provider external id"
//	@Param			type		query		string	false	"Provider type: anidb|tvdb (default anidb)"
//	@Param			page		query		int		false	"Page"
//	@Param			limit		query		int		false	"Limit"
//	@Success		200			{array}		EpisodeResponse
//	@Failure		400			{object}	httperr.APIErrorResponse
//	@Failure		500			{object}	httperr.APIErrorResponse
//	@Router			/metadata/episodes/{externalId} [get]
func (h *Handler) ListEpisodes(c *gin.Context) {
	provider, externalID, opts, ok := getEpisodesInput(c)
	if !ok {
		return
	}

	items, err := h.svc.ListEpisodes(c.Request.Context(), provider, externalID, opts)
	if err != nil {
		abortProviderErr(c, "failed to list metadata episodes", err)
		return
	}

	c.JSON(http.StatusOK, items)
}

// AddShow godoc
//
//	@Summary		Add show from metadata search item
//	@Description	Create a show and enqueue a job_shows record linked to the show
//	@Tags			metadata
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		AddShowRequest	true	"Search show item payload"
//	@Success		201		{object}	AddShowResponse
//	@Failure		400		{object}	httperr.APIErrorResponse
//	@Failure		409		{object}	httperr.APIErrorResponse
//	@Failure		500		{object}	httperr.APIErrorResponse
//	@Router			/metadata/show [post]
func (h *Handler) AddShow(c *gin.Context) {
	req, ok := httpx.AbortIfMissingContext[AddShowRequest](c, ctxAddShowRequest)
	if !ok {
		return
	}

	item, err := h.svc.AddShow(c.Request.Context(), req)
	if err != nil {
		if httpx.AbortIfDBErr(c, err, "failed to add show job") {
			return
		}
		httperr.Abort(c, httperr.Internal("failed to add show job").WithCause(err))
		return
	}

	c.JSON(http.StatusCreated, item)
}
