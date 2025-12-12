package router

import (
	"github.com/bingo-project/websocket"

	wshandler "{{.RootPackage}}/internal/{{.ServiceName}}/handler/ws"
	"{{.RootPackage}}/internal/pkg/store"
)

// RegisterWSHandlers registers all WebSocket handlers with the router.
func RegisterWSHandlers(router *websocket.Router) {
	h := wshandler.NewHandler(store.S)

	// Public methods (no auth required)
	public := router.Group()
	public.Handle("heartbeat", websocket.HeartbeatHandler)
	// public.Handle("example", h.Example)
	_ = h

	// Private methods (require auth)
	// private := router.Group(middleware.Auth)
	// private.Handle("user.info", h.UserInfo)
}
