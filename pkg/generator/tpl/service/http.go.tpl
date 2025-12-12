package {{.ServiceName}}

import (
	"github.com/gin-gonic/gin"

	"{{.RootPackage}}/internal/pkg/bootstrap"
)

// initGinEngine initializes the Gin engine with routes.
func initGinEngine() *gin.Engine {
	g := bootstrap.InitGin()

	// Install routes here
	// Example:
	// router.MapApiRouters(g)

	return g
}
