package apiserver

import (
	"github.com/gin-gonic/gin"

	"{[.RootPackage]}/internal/apiserver/controller/v1/post"
	"{[.RootPackage]}/internal/apiserver/controller/v1/user"
	"{[.RootPackage]}/internal/apiserver/store"
	"{[.RootPackage]}/internal/pkg/core"
	"{[.RootPackage]}/internal/pkg/errno"
	"{[.RootPackage]}/internal/pkg/log"
	"{[.RootPackage]}/internal/pkg/middleware"
	"{[.RootPackage]}/pkg/auth"
)

func installRouters(g *gin.Engine) error {
	// 注册 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errno.ErrPageNotFound, nil)
	})

	// 注册 /healthz handler.
	g.GET("/healthz", func(c *gin.Context) {
		log.C(c).Infow("Healthz function called")

		core.WriteResponse(c, nil, map[string]string{"status": "ok"})
	})

	// Authz
	authz, err := auth.NewAuthz(store.S.DB())
	if err != nil {
		return err
	}

	// Controller
	userController := user.NewUserController(store.S, authz)
	postController := post.NewPostController(store.S)

	// v1 group
	v1 := g.Group("/v1")

	// Login
	v1.POST("login", userController.Login)

	// User
	userV1 := v1.Group("users")
	userV1.POST("", userController.Create)                             // 创建用户
	userV1.PUT(":name/change-password", userController.ChangePassword) // 修改用户密码
	userV1.Use(middleware.Authn(), middleware.Authz(authz))
	userV1.GET("", userController.List)           // 列出用户列表，只有 root 用户才能访问
	userV1.GET(":name", userController.Get)       // 获取用户详情
	userV1.PUT(":name", userController.Update)    // 更新用户
	userV1.DELETE(":name", userController.Delete) // 删除用户

	// Post
	postV1 := v1.Group("posts", middleware.Authn())
	postV1.GET("", postController.List)                // 获取博客列表
	postV1.POST("", postController.Create)             // 创建博客
	postV1.GET(":postID", postController.Get)          // 获取博客详情
	postV1.PUT(":postID", postController.Update)       // 更新博客
	postV1.DELETE(":postID", postController.Delete)    // 删除博客
	postV1.DELETE("", postController.DeleteCollection) // 批量删除博客

	return nil
}
