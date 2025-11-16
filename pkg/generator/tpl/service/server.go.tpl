package {{.ServiceName}}

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bingo-project/component-base/log"
	"github.com/bingo-project/component-base/web"
	"github.com/gin-gonic/gin"

	genericserver "{{.RootPackage}}/internal/pkg/server"
)

// HTTPServer represents the HTTP server.
type HTTPServer struct {
	*http.Server
	engine *gin.Engine
}

// NewHTTP creates a new HTTP server instance.
func NewHTTP() *HTTPServer {
	// Set Gin mode.
	gin.SetMode(genericserver.Config.Server.Mode)

	// Create Gin engine.
	g := gin.New()

	// Install middlewares.
	installMiddlewares(g)

	// Install routes.
	installRoutes(g)

	// Create HTTP server.
	httpsrv := &http.Server{
		Addr:    genericserver.Config.Server.Addr,
		Handler: g,
	}

	return &HTTPServer{Server: httpsrv, engine: g}
}

// Run starts the HTTP server.
func (s *HTTPServer) Run() {
	log.Infow("Start to listening the incoming requests on http address", "addr", s.Addr)

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalw("Failed to start http server", "err", err)
		}
	}()
}

// Close gracefully shuts down the HTTP server.
func (s *HTTPServer) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Errorw("Failed to gracefully shutdown http server", "err", err)
	}

	log.Infow("HTTP server stopped")
}

func installMiddlewares(g *gin.Engine) {
	g.Use(gin.Recovery())
	g.Use(web.RequestID())
	g.Use(web.Context())
	g.Use(web.Logger())
}

func installRoutes(g *gin.Engine) {
	// Health check endpoint.
	g.GET("/healthz", func(c *gin.Context) {
		web.WriteResponse(c, nil, map[string]string{"status": "ok"})
	})

	// Install your routes here.
	// Example:
	// v1 := g.Group("/v1")
	// {
	//     v1.GET("/example", exampleHandler)
	// }
}