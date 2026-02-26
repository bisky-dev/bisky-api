package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keithics/devops-dashboard/api/internal/httperr"
	"github.com/keithics/devops-dashboard/api/internal/httpx"
)

func (h *Handler) Register(c *gin.Context) {
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
