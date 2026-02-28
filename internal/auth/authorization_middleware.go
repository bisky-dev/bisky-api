package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/keithics/devops-dashboard/api/internal/httperr"
)

func (h *Handler) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isPublicRoute(c) {
			c.Next()
			return
		}

		authorization := strings.TrimSpace(c.GetHeader("Authorization"))
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authorization, bearerPrefix) {
			httperr.Abort(c, httperr.Unauthorized("missing or invalid authorization header"))
			return
		}

		accessToken := strings.TrimSpace(strings.TrimPrefix(authorization, bearerPrefix))
		userID, err := h.svc.verifyAccessToken(accessToken)
		if err != nil {
			httperr.Abort(c, httperr.Unauthorized("invalid or expired access token"))
			return
		}

		c.Set(ctxUserIDKey, userID)
		c.Next()
	}
}

func isPublicRoute(c *gin.Context) bool {
	if c.Request.Method == "OPTIONS" {
		return true
	}

	path := c.FullPath()
	switch path {
	case "/health", "/swagger/*any", "/auth/register", "/auth/login", "/auth/forgot-password", "/api-keys/validate":
		return true
	default:
		return false
	}
}
