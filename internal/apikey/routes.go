package apikey

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keithics/devops-dashboard/api/internal/httpx"
	"golang.org/x/time/rate"
)

func RegisterRoutes(r *gin.Engine, h *Handler) {
	r.POST("/api-keys", h.Create)
	r.POST("/api-keys/validate", httpx.RateLimitByIP(rate.Limit(5), 10, 10*time.Minute), h.Validate)
	r.DELETE("/api-keys/:id", h.Delete)
}
