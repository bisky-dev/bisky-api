package apikey

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/keithics/devops-dashboard/api/internal/httperr"
	"github.com/keithics/devops-dashboard/api/internal/httpx"
)

// Create godoc
//
//	@Summary		Create API key
//	@Description	Create an API key and return plaintext value once
//	@Tags			api-keys
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		createAPIKeyRequest	true	"API key payload"
//	@Success		201		{object}	createAPIKeyResponse
//	@Failure		400		{object}	httperr.APIErrorResponse
//	@Failure		500		{object}	httperr.APIErrorResponse
//	@Router			/api-keys [post]
func (h *Handler) Create(c *gin.Context) {
	var req createAPIKeyRequest
	if httpx.AbortIfErr(c, c.ShouldBindJSON(&req)) {
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if httpx.AbortIfErr(c, httpx.ValidateVar(req.Name, "required,max=120", "name is invalid")) {
		return
	}

	created, err := h.svc.Create(c.Request.Context(), req.Name)
	if err != nil {
		httperr.Abort(c, httperr.Internal("failed to create api key").WithCause(err))
		return
	}

	c.JSON(http.StatusCreated, created)
}

// Validate godoc
//
//	@Summary		Validate API key
//	@Description	Validate an API key value
//	@Tags			api-keys
//	@Produce		json
//	@Success		200
//	@Failure		401	{object}	httperr.APIErrorResponse
//	@Failure		500	{object}	httperr.APIErrorResponse
//	@Router			/api-keys/validate [post]
func (h *Handler) Validate(c *gin.Context) {
	rawKey := extractAPIKey(c)
	if rawKey == "" {
		httperr.Abort(c, httperr.Unauthorized("missing api key"))
		return
	}

	ok, err := h.svc.Validate(c.Request.Context(), rawKey)
	if err != nil {
		httperr.Abort(c, httperr.Internal("failed to validate api key").WithCause(err))
		return
	}
	if !ok {
		httperr.Abort(c, httperr.Unauthorized("invalid api key"))
		return
	}

	c.Status(http.StatusOK)
}

// Delete godoc
//
//	@Summary		Delete API key
//	@Description	Delete an API key by id
//	@Tags			api-keys
//	@Produce		json
//	@Param			id	path	string	true	"API key id"
//	@Success		204
//	@Failure		400	{object}	httperr.APIErrorResponse
//	@Failure		500	{object}	httperr.APIErrorResponse
//	@Router			/api-keys/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if httpx.AbortIfErr(c, httpx.ValidateVar(id, "required,uuid4", "id is invalid")) {
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		httperr.Abort(c, httperr.Internal("failed to delete api key").WithCause(err))
		return
	}

	c.Status(http.StatusNoContent)
}

func extractAPIKey(c *gin.Context) string {
	value := strings.TrimSpace(c.GetHeader("X-API-Key"))
	if value != "" {
		return value
	}

	authorization := strings.TrimSpace(c.GetHeader("Authorization"))
	const apiKeyPrefix = "ApiKey "
	if strings.HasPrefix(authorization, apiKeyPrefix) {
		return strings.TrimSpace(strings.TrimPrefix(authorization, apiKeyPrefix))
	}

	return ""
}
