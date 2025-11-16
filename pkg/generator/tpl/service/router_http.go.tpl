package router

import (
	"github.com/bingo-project/component-base/web"
	"github.com/gin-gonic/gin"
)

// InstallHTTPRoutes registers HTTP routes.
func InstallHTTPRoutes(g *gin.Engine) {
	// Health check
	g.GET("/healthz", func(c *gin.Context) {
		web.WriteResponse(c, nil, map[string]string{"status": "ok"})
	})

	// API v1 routes
	v1 := g.Group("/v1")
	{
		// Add your routes here
		// Example:
		// v1.GET("/users", controller.ListUsers)
		// v1.POST("/users", controller.CreateUser)
		_ = v1
	}
}