package show

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keithics/devops-dashboard/api/internal/httperr"
	"github.com/keithics/devops-dashboard/api/internal/httpx"
)

// CreateShow godoc
//
//	@Summary		Create show
//	@Description	Create a new show record
//	@Tags			shows
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		createShowRequest	true	"Show payload"
//	@Success		201		{object}	showResponse
//	@Failure		400		{object}	httperr.APIErrorResponse
//	@Failure		500		{object}	httperr.APIErrorResponse
//	@Router			/shows [post]
func (h *Handler) CreateShow(c *gin.Context) {
	req, ok := httpx.AbortIfMissingContext[createShowRequest](c, ctxCreateShowRequestKey)
	if !ok {
		return
	}

	created, err := h.svc.CreateShow(c.Request.Context(), req)
	if err != nil {
		if httpx.AbortIfDBErr(c, err, "failed to create show") {
			return
		}
		httperr.Abort(c, httperr.Internal("failed to create show").WithCause(err))
		return
	}

	response, err := toShowResponse(created)
	if err != nil {
		httperr.Abort(c, httperr.Internal("failed to format show response").WithCause(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// ListShows godoc
//
//	@Summary		List shows
//	@Description	List all shows
//	@Tags			shows
//	@Produce		json
//	@Success		200	{array}		showResponse
//	@Failure		500	{object}	httperr.APIErrorResponse
//	@Router			/shows [get]
func (h *Handler) ListShows(c *gin.Context) {
	items, err := h.svc.ListShows(c.Request.Context())
	if err != nil {
		if httpx.AbortIfDBErr(c, err, "failed to list shows") {
			return
		}
		httperr.Abort(c, httperr.Internal("failed to list shows").WithCause(err))
		return
	}

	response, err := httpx.MapSliceE(items, toShowResponse)
	if err != nil {
		httperr.Abort(c, httperr.Internal("failed to format show response").WithCause(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetShow godoc
//
//	@Summary		Get show
//	@Description	Get a show by internal show id
//	@Tags			shows
//	@Produce		json
//	@Param			internalShowId	path		string	true	"Internal show UUID"
//	@Success		200			{object}	showResponse
//	@Failure		400			{object}	httperr.APIErrorResponse
//	@Failure		404			{object}	httperr.APIErrorResponse
//	@Failure		500			{object}	httperr.APIErrorResponse
//	@Router			/shows/{internalShowId} [get]
func (h *Handler) GetShow(c *gin.Context) {
	showID, ok := httpx.AbortIfMissingContext[string](c, ctxShowIDKey)
	if !ok {
		return
	}

	item, err := h.svc.GetShowByID(c.Request.Context(), showID)
	if err != nil {
		if httpx.AbortIfDBErr(c, err, "failed to get show") {
			return
		}
		httperr.Abort(c, httperr.Internal("failed to get show").WithCause(err))
		return
	}

	response, err := toShowResponse(item)
	if err != nil {
		httperr.Abort(c, httperr.Internal("failed to format show response").WithCause(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateShow godoc
//
//	@Summary		Update show
//	@Description	Update a show by internal show id
//	@Tags			shows
//	@Accept			json
//	@Produce		json
//	@Param			internalShowId	path		string			true	"Internal show UUID"
//	@Param			payload			body		updateShowRequest	true	"Show payload"
//	@Success		200				{object}	showResponse
//	@Failure		400				{object}	httperr.APIErrorResponse
//	@Failure		404				{object}	httperr.APIErrorResponse
//	@Failure		500				{object}	httperr.APIErrorResponse
//	@Router			/shows/{internalShowId} [put]
func (h *Handler) UpdateShow(c *gin.Context) {
	showID, ok := httpx.AbortIfMissingContext[string](c, ctxShowIDKey)
	if !ok {
		return
	}
	req, ok := httpx.AbortIfMissingContext[updateShowRequest](c, ctxUpdateShowRequestKey)
	if !ok {
		return
	}

	updated, err := h.svc.UpdateShow(c.Request.Context(), showID, req)
	if err != nil {
		if httpx.AbortIfDBErr(c, err, "failed to update show") {
			return
		}
		httperr.Abort(c, httperr.Internal("failed to update show").WithCause(err))
		return
	}

	response, err := toShowResponse(updated)
	if err != nil {
		httperr.Abort(c, httperr.Internal("failed to format show response").WithCause(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteShow godoc
//
//	@Summary		Delete show
//	@Description	Delete a show by internal show id
//	@Tags			shows
//	@Produce		json
//	@Param			internalShowId	path	string	true	"Internal show UUID"
//	@Success		204
//	@Failure		400	{object}	httperr.APIErrorResponse
//	@Failure		404	{object}	httperr.APIErrorResponse
//	@Failure		500	{object}	httperr.APIErrorResponse
//	@Router			/shows/{internalShowId} [delete]
func (h *Handler) DeleteShow(c *gin.Context) {
	showID, ok := httpx.AbortIfMissingContext[string](c, ctxShowIDKey)
	if !ok {
		return
	}

	err := h.svc.DeleteShow(c.Request.Context(), showID)
	if err != nil {
		if httpx.AbortIfDBErr(c, err, "failed to delete show") {
			return
		}
		httperr.Abort(c, httperr.Internal("failed to delete show").WithCause(err))
		return
	}

	c.Status(http.StatusNoContent)
}
