package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *Handler) {
	auth := r.Group("/auth")
	auth.POST("/register", h.BindRegister(), h.Register)
	auth.POST("/login", h.BindLogin(), h.Login)
	auth.POST("/forgot-password", h.BindForgotPassword(), h.ForgotPassword)
}
