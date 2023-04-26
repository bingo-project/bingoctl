package post

import (
	"{[.RootPackage]}/internal/apiserver/biz"
	"{[.RootPackage]}/internal/apiserver/store"
)

// PostController 是 post 模块在 Controller 层的实现，用来处理博客模块的请求.
type PostController struct {
	b biz.IBiz
}

// NewPostController 创建一个 post controller.
func NewPostController(ds store.IStore) *PostController {
	return &PostController{b: biz.NewBiz(ds)}
}
