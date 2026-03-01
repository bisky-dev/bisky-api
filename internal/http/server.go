package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/keithics/devops-dashboard/api/docs/swagger"
	"github.com/keithics/devops-dashboard/api/internal/apikey"
	"github.com/keithics/devops-dashboard/api/internal/auth"
	"github.com/keithics/devops-dashboard/api/internal/config"
	"github.com/keithics/devops-dashboard/api/internal/db/sqlc"
	"github.com/keithics/devops-dashboard/api/internal/episode"
	"github.com/keithics/devops-dashboard/api/internal/hooks"
	"github.com/keithics/devops-dashboard/api/internal/hooksettings"
	"github.com/keithics/devops-dashboard/api/internal/httperr"
	"github.com/keithics/devops-dashboard/api/internal/httpx"
	"github.com/keithics/devops-dashboard/api/internal/metadata"
	workermeta "github.com/keithics/devops-dashboard/api/internal/metadata/provider"
	"github.com/keithics/devops-dashboard/api/internal/metadata/provider/providers/anilist"
	"github.com/keithics/devops-dashboard/api/internal/metadata/provider/providers/tvdb"
	"github.com/keithics/devops-dashboard/api/internal/show"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

func NewServer(cfg config.Config, pool *pgxpool.Pool) *Server {
	r := gin.New()
	r.Use(httperr.Recovery())
	r.Use(gin.Logger())
	r.Use(corsMiddleware())
	r.Use(httperr.Middleware())

	q := sqlc.New(pool)
	apiKeyHandler := apikey.NewHandler(pool)
	authHandler := auth.NewHandler(q, cfg.TokenEncryptionKey)
	var hookDispatcher hooks.Dispatcher = hooks.NoopDispatcher{}
	httpHookDispatcher, err := hooks.NewHTTPDispatcher(pool)
	if err != nil {
		log.Printf("failed to initialize hook dispatcher, using noop: %v", err)
	} else {
		hookDispatcher = httpHookDispatcher
	}
	episodeHandler := episode.NewHandlerWithHooks(q, hookDispatcher)
	anilistProvider := anilist.New()
	metadataRegistry := workermeta.NewRegistry(map[workermeta.ProviderName]workermeta.Provider{
		workermeta.ProviderAniDB:   anilistProvider,
		workermeta.ProviderAniList: anilistProvider,
		workermeta.ProviderTVDB:    tvdb.New(),
	})
	metadataHandler := metadata.NewHandler(metadata.NewService(workermeta.NewService(metadataRegistry)))
	showHandler := show.NewHandlerWithHooks(q, hookDispatcher)
	hookSettingsHandler, err := hooksettings.NewHandler(pool)
	if err != nil {
		log.Printf("failed to initialize hook settings handler: %v", err)
	}

	r.GET("/health", httpx.RateLimitByIP(rate.Limit(10), 20, 5*time.Minute), healthHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	auth.RegisterRoutes(r, authHandler)
	r.Use(authHandler.RequireAuth())
	episode.RegisterRoutes(r, episodeHandler)
	apikey.RegisterRoutes(r, apiKeyHandler)
	metadata.RegisterRoutes(r, metadataHandler)
	show.RegisterRoutes(r, showHandler)
	if hookSettingsHandler != nil {
		hooksettings.RegisterRoutes(r, hookSettingsHandler)
	}

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
