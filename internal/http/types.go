package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/keithics/devops-dashboard/api/internal/config"
)

type Server struct {
	cfg    config.Config
	pool   *pgxpool.Pool
	router *gin.Engine
}

type healthResponse struct {
	OK bool `json:"ok"`
}
