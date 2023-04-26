package user

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"

	"{[.RootPackage]}/internal/pkg/core"
	"{[.RootPackage]}/internal/pkg/errno"
	"{[.RootPackage]}/internal/pkg/log"
	v1 "{[.RootPackage]}/pkg/api/{[.AppName]}/v1"
)

// ChangePassword 修改指定用户的密码.
func (ctrl *UserController) ChangePassword(c *gin.Context) {
	log.C(c).Infow("Change password function called")

	var r v1.ChangePasswordRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)

		return
	}

	username := c.Param("name")
	err := ctrl.b.Users().ChangePassword(c, username, &r)
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
