package episode

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keithics/devops-dashboard/api/internal/httperr"
	"github.com/keithics/devops-dashboard/api/internal/httpx"
)

// CreateEpisode godoc
//
//	@Summary		Create episode
//	@Description	Create a new episode record
//	@Tags			episodes
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		createEpisodeRequest	true	"Episode payload"
//	@Success		201		{object}	episodeResponse
//	@Failure		400		{object}	httperr.APIErrorResponse
//	@Failure		500		{object}	httperr.APIErrorResponse
//	@Router			/episodes [post]
func (h *Handler) CreateEpisode(c *gin.Context) {
	req, ok := httpx.AbortIfMissingContext[createEpisodeRequest](c, ctxCreateEpisodeRequestKey)
	if !ok {
		return
	}

	created, err := h.svc.CreateEpisode(c.Request.Context(), req)
	if err != nil {
		if httpx.AbortIfDBErr(c, err, "failed to create episode") {
			return
		}
		httperr.Abort(c, httperr.Internal("failed to create episode").WithCause(err))
		return
	}

	response, err := toEpisodeResponse(created)
	if err != nil {
		httperr.Abort(c, httperr.Internal("failed to format episode response").WithCause(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// ListEpisodes godoc
//
//	@Summary		List episodes
//	@Description	List all episodes
//	@Tags			episodes
//	@Produce		json
//	@Success		200	{array}		episodeResponse
//	@Failure		500	{object}	httperr.APIErrorResponse
//	@Router			/episodes [get]
func (h *Handler) ListEpisodes(c *gin.Context) {
	items, err := h.svc.ListEpisodes(c.Request.Context())
	if err != nil {
		if httpx.AbortIfDBErr(c, err, "failed to list episodes") {
			return
		}
		httperr.Abort(c, httperr.Internal("failed to list episodes").WithCause(err))
		return
	}

	response, err := httpx.MapSliceE(items, toEpisodeResponse)
	if err != nil {
		httperr.Abort(c, httperr.Internal("failed to format episode response").WithCause(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetEpisode godoc
//
//	@Summary		Get episode
//	@Description	Get an episode by internal episode id
//	@Tags			episodes
//	@Produce		json
//	@Param			internalEpisodeId	path		string	true	"Internal episode UUID"
//	@Success		200				{object}	episodeResponse
//	@Failure		400				{object}	httperr.APIErrorResponse
//	@Failure		404				{object}	httperr.APIErrorResponse
//	@Failure		500				{object}	httperr.APIErrorResponse
//	@Router			/episodes/{internalEpisodeId} [get]
func (h *Handler) GetEpisode(c *gin.Context) {
	episodeID, ok := httpx.AbortIfMissingContext[string](c, ctxEpisodeIDKey)
	if !ok {
		return
	}

	item, err := h.svc.GetEpisodeByID(c.Request.Context(), episodeID)
	if err != nil {
		if httpx.AbortIfDBErr(c, err, "failed to get episode") {
			return
		}
		httperr.Abort(c, httperr.Internal("failed to get episode").WithCause(err))
		return
	}

	response, err := toEpisodeResponse(item)
	if err != nil {
		httperr.Abort(c, httperr.Internal("failed to format episode response").WithCause(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateEpisode godoc
//
//	@Summary		Update episode
//	@Description	Update an episode by internal episode id
//	@Tags			episodes
//	@Accept			json
//	@Produce		json
//	@Param			internalEpisodeId	path		string			 true	"Internal episode UUID"
//	@Param			payload			body		updateEpisodeRequest	true	"Episode payload"
//	@Success		200				{object}	episodeResponse
//	@Failure		400				{object}	httperr.APIErrorResponse
//	@Failure		404				{object}	httperr.APIErrorResponse
//	@Failure		500				{object}	httperr.APIErrorResponse
//	@Router			/episodes/{internalEpisodeId} [put]
func (h *Handler) UpdateEpisode(c *gin.Context) {
	episodeID, ok := httpx.AbortIfMissingContext[string](c, ctxEpisodeIDKey)
	if !ok {
		return
	}
	req, ok := httpx.AbortIfMissingContext[updateEpisodeRequest](c, ctxUpdateEpisodeRequestKey)
	if !ok {
		return
	}

	updated, err := h.svc.UpdateEpisode(c.Request.Context(), episodeID, req)
	if err != nil {
		if httpx.AbortIfDBErr(c, err, "failed to update episode") {
			return
		}
		httperr.Abort(c, httperr.Internal("failed to update episode").WithCause(err))
		return
	}

	response, err := toEpisodeResponse(updated)
	if err != nil {
		httperr.Abort(c, httperr.Internal("failed to format episode response").WithCause(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteEpisode godoc
//
//	@Summary		Delete episode
//	@Description	Delete an episode by internal episode id
//	@Tags			episodes
//	@Produce		json
//	@Param			internalEpisodeId	path	string	true	"Internal episode UUID"
//	@Success		204
//	@Failure		400	{object}	httperr.APIErrorResponse
//	@Failure		404	{object}	httperr.APIErrorResponse
//	@Failure		500	{object}	httperr.APIErrorResponse
//	@Router			/episodes/{internalEpisodeId} [delete]
func (h *Handler) DeleteEpisode(c *gin.Context) {
	episodeID, ok := httpx.AbortIfMissingContext[string](c, ctxEpisodeIDKey)
	if !ok {
		return
	}

	if err := h.svc.DeleteEpisode(c.Request.Context(), episodeID); err != nil {
		if httpx.AbortIfDBErr(c, err, "failed to delete episode") {
			return
		}
		httperr.Abort(c, httperr.Internal("failed to delete episode").WithCause(err))
		return
	}

	c.Status(http.StatusNoContent)
}
