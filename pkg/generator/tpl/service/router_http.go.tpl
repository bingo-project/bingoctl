package router

import (
	"github.com/gin-gonic/gin"

	"{{.RootPackage}}/internal/pkg/core"
)

// MapRoutes registers HTTP routes.
func MapRoutes(g *gin.Engine) {
	// Health check
	g.GET("/healthz", func(c *gin.Context) {
		core.Response(c, map[string]string{"status": "ok"}, nil)
	})

	// API v1 routes
	v1 := g.Group("/v1")
	{
		// Add your routes here
		// Example:
		// v1.GET("/users", handler.ListUsers)
		// v1.POST("/users", handler.CreateUser)
		_ = v1
	}
}
