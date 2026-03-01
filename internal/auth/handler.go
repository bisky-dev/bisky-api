package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keithics/devops-dashboard/api/internal/httperr"
	"github.com/keithics/devops-dashboard/api/internal/httpx"
)

// Register godoc
//
//	@Summary		Register user
//	@Description	Create a new user account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		registerRequest	true	"Register payload"
//	@Success		201		{object}	registerResponse
//	@Failure		400		{object}	httperr.APIErrorResponse
//	@Failure		409		{object}	httperr.APIErrorResponse
//	@Failure		500		{object}	httperr.APIErrorResponse
//	@Router			/auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	const registrationDisabled = true
	if registrationDisabled {
		httperr.Abort(c, httperr.New(http.StatusForbidden, "forbidden", "registration is disabled"))
		return
	}

	req, ok := httpx.AbortIfMissingContext[registerRequest](c, ctxRegisterRequestKey)
	if !ok {
		return
	}

	user, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		if httpx.AbortIfDBErr(c, err, "failed to create user") {
			return
		}
		httperr.Abort(c, httperr.Internal("failed to create user").WithCause(err))
		return
	}

	c.JSON(http.StatusCreated, registerResponse{
		User: userResponse{ID: user.ID, Email: user.Email, CreatedAt: user.CreatedAt},
	})
}

// Login godoc
//
//	@Summary		Login user
//	@Description	Authenticate a user and return an access token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		loginRequest	true	"Login payload"
//	@Success		200		{object}	loginResponse
//	@Failure		400		{object}	httperr.APIErrorResponse
//	@Failure		401		{object}	httperr.APIErrorResponse
//	@Failure		500		{object}	httperr.APIErrorResponse
//	@Router			/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	req, ok := httpx.AbortIfMissingContext[loginRequest](c, ctxLoginRequestKey)
	if !ok {
		return
	}

	user, accessToken, expiresIn, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, errInvalidCredentials) {
			httperr.Abort(c, httperr.Unauthorized("invalid email or password"))
			return
		}
		if httpx.AbortIfDBErr(c, err, "failed to login") {
			return
		}
		httperr.Abort(c, httperr.Internal("failed to login").WithCause(err))
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		User:        userResponse{ID: user.ID, Email: user.Email, CreatedAt: user.CreatedAt},
	})
}

// ForgotPassword godoc
//
//	@Summary		Forgot password
//	@Description	Request a password reset token for an email address
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		forgotPasswordRequest	true	"Forgot password payload"
//	@Success		202		{object}	forgotPasswordResponse
//	@Failure		400		{object}	httperr.APIErrorResponse
//	@Failure		500		{object}	httperr.APIErrorResponse
//	@Router			/auth/forgot-password [post]
func (h *Handler) ForgotPassword(c *gin.Context) {
	req, ok := httpx.AbortIfMissingContext[forgotPasswordRequest](c, ctxForgotPasswordRequestKey)
	if !ok {
		return
	}

	err := h.svc.ForgotPassword(c.Request.Context(), req)
	if err != nil {
		if httpx.AbortIfDBErr(c, err, "failed to process forgot password") {
			return
		}
		httperr.Abort(c, httperr.Internal("failed to process forgot password").WithCause(err))
		return
	}

	c.JSON(http.StatusAccepted, forgotPasswordResponse{
		Message: "If an account exists for that email, a reset link will be sent.",
	})
}
