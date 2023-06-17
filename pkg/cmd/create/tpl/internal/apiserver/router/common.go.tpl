package router

import (
	"github.com/gin-gonic/gin"

	"{[.RootPackage]}/internal/apiserver/controller/v1/common"
	"{[.RootPackage]}/internal/pkg/core"
	"{[.RootPackage]}/internal/pkg/errno"
)

func MapCommonRouters(g *gin.Engine) {
	// 注册 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errno.ErrPageNotFound, nil)
	})

	// 注册 /healthz handler.
	commonController := common.NewCommonController()
	g.GET("/healthz", commonController.Healthz)
}
