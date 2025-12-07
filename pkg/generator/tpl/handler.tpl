package {{.PackageName}}

import (
	"github.com/bingo-project/component-base/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"{{.RootPackage}}/{{.BizPath}}"
	"{{.RootPackage}}/{{.StorePath}}"
	"{{.RootPackage}}/internal/pkg/core"
	"{{.RootPackage}}/internal/pkg/errno"
	v1 "{{.RootPackage}}/{{.RequestPath}}{{.RelativePath}}"
	"{{.RootPackage}}/pkg/auth"
)

type {{.StructName}}Handler struct {
	a *auth.Authz
	b biz.IBiz
}

func New{{.StructName}}Handler(ds store.IStore, a *auth.Authz) *{{.StructName}}Handler {
	return &{{.StructName}}Handler{a: a, b: biz.NewBiz(ds)}
}

// List
// @Summary    List {{.VariableNamePlural}}
// @Security   Bearer
// @Tags       {{.StructName}}
// @Accept     application/json
// @Produce    json
// @Param      request	 query	    v1.List{{.StructName}}Request	 true  "Param"
// @Success	   200		{object}	v1.List{{.StructName}}Response
// @Failure	   400		{object}	core.ErrResponse
// @Failure	   500		{object}	core.ErrResponse
// @Router    /v1/{{.VariableNamePlural}} [GET]
func (h *{{.StructName}}Handler) List(c *gin.Context) {
	log.C(c).Infow("List {{.VariableName}} function called")

	var req v1.List{{.StructName}}Request
	if err := c.ShouldBindQuery(&req); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)

		return
	}

	resp, err := h.b.{{.StructNamePlural}}().List(c, &req)
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, resp)
}

// Create
// @Summary    Create {{.VariableName}}
// @Security   Bearer
// @Tags       {{.StructName}}
// @Accept     application/json
// @Produce    json
// @Param      request	 body	    v1.Create{{.StructName}}Request	 true  "Param"
// @Success	   200		{object}	v1.{{.StructName}}Info
// @Failure	   400		{object}	core.ErrResponse
// @Failure	   500		{object}	core.ErrResponse
// @Router    /v1/{{.VariableNamePlural}} [POST]
func (h *{{.StructName}}Handler) Create(c *gin.Context) {
	log.C(c).Infow("Create {{.VariableName}} function called")

	var req v1.Create{{.StructName}}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)

		return
	}

	// Create {{.VariableName}}
	resp, err := h.b.{{.StructNamePlural}}().Create(c, &req)
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, resp)
}

// Get
// @Summary    Get {{.VariableName}} info
// @Security   Bearer
// @Tags       {{.StructName}}
// @Accept     application/json
// @Produce    json
// @Param      id	     path	    string            		 true  "ID"
// @Success	   200		{object}	v1.{{.StructName}}Info
// @Failure	   400		{object}	core.ErrResponse
// @Failure	   500		{object}	core.ErrResponse
// @Router    /v1/{{.VariableNamePlural}}/{id} [GET]
func (h *{{.StructName}}Handler) Get(c *gin.Context) {
	log.C(c).Infow("Get {{.VariableName}} function called")

	ID := cast.ToUint(c.Param("id"))
	{{.VariableName}}, err := h.b.{{.StructNamePlural}}().Get(c, ID)
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, {{.VariableName}})
}

// Update
// @Summary    Update {{.VariableName}} info
// @Security   Bearer
// @Tags       {{.StructName}}
// @Accept     application/json
// @Produce    json
// @Param      id	     path	    string            		 true  "ID"
// @Param      request	 body	    v1.Update{{.StructName}}Request	 true  "Param"
// @Success	   200		{object}	v1.{{.StructName}}Info
// @Failure	   400		{object}	core.ErrResponse
// @Failure	   500		{object}	core.ErrResponse
// @Router    /v1/{{.VariableNamePlural}}/{id} [PUT]
func (h *{{.StructName}}Handler) Update(c *gin.Context) {
	log.C(c).Infow("Update {{.VariableName}} function called")

	var req v1.Update{{.StructName}}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)

		return
	}

	ID := cast.ToUint(c.Param("id"))
	resp, err := h.b.{{.StructNamePlural}}().Update(c, ID, &req)
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, resp)
}

// Delete
// @Summary    Delete {{.VariableName}}
// @Security   Bearer
// @Tags       {{.StructName}}
// @Accept     application/json
// @Produce    json
// @Param      id	    path	    string            true  "ID"
// @Success	   200		{object}	nil
// @Failure	   400		{object}	core.ErrResponse
// @Failure	   500		{object}	core.ErrResponse
// @Router    /v1/{{.VariableNamePlural}}/{id} [DELETE]
func (h *{{.StructName}}Handler) Delete(c *gin.Context) {
	log.C(c).Infow("Delete {{.VariableName}} function called")

	ID := cast.ToUint(c.Param("id"))
	if err := h.b.{{.StructNamePlural}}().Delete(c, ID); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
