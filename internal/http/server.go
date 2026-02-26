package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/keithics/devops-dashboard/api/internal/auth"
	"github.com/keithics/devops-dashboard/api/internal/config"
	"github.com/keithics/devops-dashboard/api/internal/db/sqlc"
	"github.com/keithics/devops-dashboard/api/internal/httperr"
)

func NewServer(cfg config.Config, pool *pgxpool.Pool) *Server {
	r := gin.New()
	r.Use(httperr.Recovery())
	r.Use(gin.Logger())
	r.Use(corsMiddleware())
	r.Use(httperr.Middleware())

	q := sqlc.New(pool)
	authHandler := auth.NewHandler(q, cfg.TokenEncryptionKey)

	r.GET("/health", healthHandler)
	r.POST("/auth/register", authHandler.BindRegister(), authHandler.Register)
	r.POST("/auth/login", authHandler.BindLogin(), authHandler.Login)
	r.POST("/auth/forgot-password", authHandler.BindForgotPassword(), authHandler.ForgotPassword)

	return &Server{
		cfg:    cfg,
		pool:   pool,
		router: r,
	}
}

func (s *Server) Router() http.Handler {
	return s.router
}

// healthHandler godoc
//
//	@Summary		Health check
//	@Description	Service liveness endpoint
//	@Tags			system
//	@Produce		json
//	@Success		200	{object}	healthResponse
//	@Router			/health [get]
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, healthResponse{OK: true})
}
