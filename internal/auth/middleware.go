package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/keithics/devops-dashboard/api/internal/httpx"
	normalizeutil "github.com/keithics/devops-dashboard/api/internal/utils/normalize"
)

func (h *Handler) BindRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req registerRequest
		if httpx.AbortIfErr(c, c.ShouldBindJSON(&req)) {
			return
		}

		req.Email = normalizeutil.LowerString(req.Email)
		if httpx.AbortIfErr(c, validateRegisterRequest(req)) {
			return
		}

		c.Set(ctxRegisterRequestKey, req)
		c.Next()
	}
}

func (h *Handler) BindLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest
		if httpx.AbortIfErr(c, c.ShouldBindJSON(&req)) {
			return
		}

		req.Email = normalizeutil.LowerString(req.Email)
		if httpx.AbortIfErr(c, validateLoginRequest(req)) {
			return
		}

		c.Set(ctxLoginRequestKey, req)
		c.Next()
	}
}

func (h *Handler) BindForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req forgotPasswordRequest
		if httpx.AbortIfErr(c, c.ShouldBindJSON(&req)) {
			return
		}

		req.Email = normalizeutil.LowerString(req.Email)
		if httpx.AbortIfErr(c, validateForgotPasswordRequest(req)) {
			return
		}

		c.Set(ctxForgotPasswordRequestKey, req)
		c.Next()
	}
}
