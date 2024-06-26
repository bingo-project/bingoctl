package file

import (
	"github.com/bingo-project/component-base/log"
	"github.com/gin-gonic/gin"

	"{[.RootPackage]}/internal/apiserver/biz"
	v1 "{[.RootPackage]}/internal/apiserver/http/request/v1"
	"{[.RootPackage]}/internal/apiserver/store"
	"{[.RootPackage]}/internal/pkg/core"
	"{[.RootPackage]}/internal/pkg/errno"
	"{[.RootPackage]}/pkg/auth"
)

type FileController struct {
	a *auth.Authz
	b biz.IBiz
}

func NewFileController(ds store.IStore, a *auth.Authz) *FileController {
	return &FileController{a: a, b: biz.NewBiz(ds)}
}

// Upload
// @Summary    Upload file
// @Security   Bearer
// @Tags       File
// @Accept     multipart/form-data
// @Produce    json
// @Param      file     formData    file    true  "File"
// @Success	   200		{object}	string
// @Failure	   400		{object}	core.ErrResponse
// @Failure	   500		{object}	core.ErrResponse
// @Router    /v1/file/upload [POST].
func (ctrl *FileController) Upload(c *gin.Context) {
	log.C(c).Infow("Upload file function called")

	var req v1.UploadFileRequest
	if err := c.ShouldBind(&req); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		return
	}

	// Create file
	resp, err := ctrl.b.Files().Upload(c, &req)
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, resp)
}
