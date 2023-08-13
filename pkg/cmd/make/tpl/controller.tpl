package {{.PackageName}}

import (
	"github.com/asaskevich/govalidator"
	"github.com/bingo-project/component-base/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"{{.RootPackage}}/{{.BizPath}}"
	"{{.RootPackage}}/{{.StorePath}}"
	"{{.RootPackage}}/internal/pkg/core"
	"{{.RootPackage}}/internal/pkg/errno"
	v1 "{{.RootPackage}}/{{.RequestPath}}"
	"{{.RootPackage}}/pkg/auth"
)

type {{.StructName}}Controller struct {
	a *auth.Authz
	b biz.IBiz
}

func New{{.StructName}}Controller(ds store.IStore, a *auth.Authz) *{{.StructName}}Controller {
	return &{{.StructName}}Controller{a: a, b: biz.NewBiz(ds)}
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
func (ctrl *{{.StructName}}Controller) List(c *gin.Context) {
	log.C(c).Infow("List {{.VariableName}} function called")

	var r v1.List{{.StructName}}Request
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		return
	}

	resp, err := ctrl.b.{{.StructNamePlural}}().List(c, r.Offset, r.Limit)
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, resp)
}

// Create
// @Summary    Create a {{.VariableName}}
// @Security   Bearer
// @Tags       {{.StructName}}
// @Accept     application/json
// @Produce    json
// @Param      request	 body	    v1.Create{{.StructName}}Request	 true  "Param"
// @Success	   200		{object}	v1.Get{{.StructName}}Response
// @Failure	   400		{object}	core.ErrResponse
// @Failure	   500		{object}	core.ErrResponse
// @Router    /v1/{{.VariableNamePlural}} [POST]
func (ctrl *{{.StructName}}Controller) Create(c *gin.Context) {
	log.C(c).Infow("Create {{.VariableName}} function called")

	var r v1.Create{{.StructName}}Request
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		return
	}

	// Validator
	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)

		return
	}

	// Create {{.VariableName}}
	resp, err := ctrl.b.{{.StructNamePlural}}().Create(c, &r)
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
// @Success	   200		{object}	v1.Get{{.StructName}}Response
// @Failure	   400		{object}	core.ErrResponse
// @Failure	   500		{object}	core.ErrResponse
// @Router    /v1/{{.VariableNamePlural}}/{id} [GET]
func (ctrl *{{.StructName}}Controller) Get(c *gin.Context) {
	log.C(c).Infow("Get {{.VariableName}} function called")

	ID := cast.ToUint(c.Param("id"))
	{{.VariableName}}, err := ctrl.b.{{.StructNamePlural}}().Get(c, ID)
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
// @Success	   200		{object}	v1.Get{{.StructName}}Response
// @Failure	   400		{object}	core.ErrResponse
// @Failure	   500		{object}	core.ErrResponse
// @Router    /v1/{{.VariableNamePlural}}/{id} [PUT]
func (ctrl *{{.StructName}}Controller) Update(c *gin.Context) {
	log.C(c).Infow("Update {{.VariableName}} function called")

	var r v1.Update{{.StructName}}Request
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)

		return
	}

	ID := cast.ToUint(c.Param("id"))
	resp, err := ctrl.b.{{.StructNamePlural}}().Update(c, ID, &r)
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, resp)
}

// Delete
// @Summary    Delete a {{.VariableName}}
// @Security   Bearer
// @Tags       {{.StructName}}
// @Accept     application/json
// @Produce    json
// @Param      id	    path	    string            true  "ID"
// @Success	   200		{object}	nil
// @Failure	   400		{object}	core.ErrResponse
// @Failure	   500		{object}	core.ErrResponse
// @Router    /v1/{{.VariableNamePlural}}/{id} [DELETE]
func (ctrl *{{.StructName}}Controller) Delete(c *gin.Context) {
	log.C(c).Infow("Delete {{.VariableName}} function called")

	ID := cast.ToUint(c.Param("id"))
	if err := ctrl.b.{{.StructNamePlural}}().Delete(c, ID); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
