package http

import (
	"{{.RootPackage}}/internal/{{.ServiceName}}/biz"
	"{{.RootPackage}}/internal/pkg/store"
)

type Handler struct {
	b biz.IBiz
}

func NewHandler(ds store.IStore) *Handler {
	return &Handler{b: biz.NewBiz(ds)}
}

// Example handler method
// func (h *Handler) List(c *gin.Context) {
//     // Implementation here
//     core.Response(c, nil, nil)
// }
