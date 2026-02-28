package hooksettings

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keithics/devops-dashboard/api/internal/httperr"
	"github.com/keithics/devops-dashboard/api/internal/httpx"
)

// ListHooks godoc
//
//	@Summary		List hook settings
//	@Description	View predefined hooks and configured target URLs
//	@Tags			settings
//	@Produce		json
//	@Success		200	{array}	hooks.Config
//	@Failure		500	{object}	httperr.APIErrorResponse
//	@Router			/settings/hooks [get]
func (h *Handler) ListHooks(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		httperr.Abort(c, httperr.Internal("failed to list hook settings").WithCause(err))
		return
	}
	c.JSON(http.StatusOK, items)
}

// ListHookKeys godoc
//
//	@Summary		List hook keys
//	@Description	View all predefined hook event keys
//	@Tags			settings
//	@Produce		json
//	@Success		200	{array}	string
//	@Router			/settings/hooks/keys [get]
func (h *Handler) ListHookKeys(c *gin.Context) {
	c.JSON(http.StatusOK, h.svc.ListKeys())
}

// UpsertHooks godoc
//
//	@Summary		Upsert hook settings
//	@Description	Upsert target URL by predefined hook event name
//	@Tags			settings
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		upsertHooksRequest	true	"Hooks settings payload"
//	@Success		200		{array}		hooks.Config
//	@Failure		400		{object}	httperr.APIErrorResponse
//	@Failure		500		{object}	httperr.APIErrorResponse
//	@Router			/settings/hooks [put]
func (h *Handler) UpsertHooks(c *gin.Context) {
	var req upsertHooksRequest
	if httpx.AbortIfErr(c, c.ShouldBindJSON(&req)) {
		return
	}

	items, err := h.svc.Upsert(c.Request.Context(), req)
	if err != nil {
		httperr.Abort(c, httperr.BadRequest(err.Error()))
		return
	}
	c.JSON(http.StatusOK, items)
}
