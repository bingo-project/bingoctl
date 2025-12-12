package ws

import (
	"{{.RootPackage}}/internal/{{.ServiceName}}/biz"
	"{{.RootPackage}}/internal/pkg/store"
)

// Handler handles WebSocket business methods.
type Handler struct {
	b biz.IBiz
}

// NewHandler creates a new WebSocket handler.
func NewHandler(ds store.IStore) *Handler {
	return &Handler{b: biz.NewBiz(ds)}
}

// Example handler method
// func (h *Handler) Example(ctx context.Context, req *websocket.Request) (any, error) {
//     // Implementation here
//     return nil, nil
// }
