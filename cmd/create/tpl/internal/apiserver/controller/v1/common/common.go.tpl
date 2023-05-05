package common

import (
	"github.com/gin-gonic/gin"

	"{[.RootPackage]}/internal/apiserver/biz"
	"{[.RootPackage]}/internal/apiserver/store"
	"{[.RootPackage]}/internal/pkg/core"
	"{[.RootPackage]}/internal/pkg/log"
	v1 "{[.RootPackage]}/pkg/api/goer/v1"
	"{[.RootPackage]}/pkg/auth"
)

// CommonController 是 common 模块在 Controller 层的实现，用来处理用户模块的请求.
type CommonController struct {
	a *auth.Authz
	b biz.IBiz
}

// NewCommonController 创建一个 common controller.
func NewCommonController(ds store.IStore, a *auth.Authz) *CommonController {
	return &CommonController{a: a, b: biz.NewBiz(ds)}
}

// Healthz
// @Summary    Heath check
// @Tags       Common
// @Accept     application/json
// @Produce    json
// @Success	   200		{object}	v1.HealthzResponse
// @Failure	   400		{object}	core.ErrResponse
// @Failure	   500		{object}	core.ErrResponse
// @Router    /healthz  [GET].
func (ctrl *CommonController) Healthz(c *gin.Context) {
	log.C(c).Infow("Healthz function called")

	data := &v1.HealthzResponse{Status: "ok"}

	core.WriteResponse(c, nil, data)
}
