package auth

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keithics/devops-dashboard/api/internal/httpx"
	"golang.org/x/time/rate"
)

func RegisterRoutes(r *gin.Engine, h *Handler) {
	auth := r.Group("/auth")
	auth.POST("/register", httpx.RateLimitByIP(rate.Limit(1), 3, 10*time.Minute), h.BindRegister(), h.Register)
	auth.POST("/login", httpx.RateLimitByIP(rate.Limit(2), 6, 10*time.Minute), h.BindLogin(), h.Login)
	auth.POST("/forgot-password", httpx.RateLimitByIP(rate.Limit(1), 3, 10*time.Minute), h.BindForgotPassword(), h.ForgotPassword)
}
