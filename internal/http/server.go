package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/keithics/devops-dashboard/api/docs/swagger"
	"github.com/keithics/devops-dashboard/api/internal/auth"
	"github.com/keithics/devops-dashboard/api/internal/config"
	"github.com/keithics/devops-dashboard/api/internal/db/sqlc"
	"github.com/keithics/devops-dashboard/api/internal/episode"
	"github.com/keithics/devops-dashboard/api/internal/httperr"
	jobshow "github.com/keithics/devops-dashboard/api/internal/job/show"
	"github.com/keithics/devops-dashboard/api/internal/metadata"
	workermeta "github.com/keithics/devops-dashboard/api/internal/metadata/provider"
	"github.com/keithics/devops-dashboard/api/internal/metadata/provider/providers/anilist"
	"github.com/keithics/devops-dashboard/api/internal/metadata/provider/providers/tvdb"
	"github.com/keithics/devops-dashboard/api/internal/show"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewServer(cfg config.Config, pool *pgxpool.Pool) *Server {
	r := gin.New()
	r.Use(httperr.Recovery())
	r.Use(gin.Logger())
	r.Use(corsMiddleware())
	r.Use(httperr.Middleware())

	q := sqlc.New(pool)
	authHandler := auth.NewHandler(q, cfg.TokenEncryptionKey)
	episodeHandler := episode.NewHandler(q)
	anilistProvider := anilist.New()
	metadataRegistry := workermeta.NewRegistry(map[workermeta.ProviderName]workermeta.Provider{
		workermeta.ProviderAniDB:   anilistProvider,
		workermeta.ProviderAniList: anilistProvider,
		workermeta.ProviderTVDB:    tvdb.New(),
	})
	jobShowService := jobshow.NewService(pool)
	metadataHandler := metadata.NewHandler(metadata.NewService(workermeta.NewService(metadataRegistry), jobShowService))
	showHandler := show.NewHandler(q)

	r.GET("/health", healthHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	auth.RegisterRoutes(r, authHandler)
	r.Use(authHandler.RequireAuth())
	episode.RegisterRoutes(r, episodeHandler)
	metadata.RegisterRoutes(r, metadataHandler)
	show.RegisterRoutes(r, showHandler)

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
