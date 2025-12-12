package {{.ServiceName}}

import (
	"context"
	"net/http"

	"github.com/bingo-project/component-base/web/token"
	"github.com/bingo-project/websocket"
	"github.com/gin-gonic/gin"
	gorillaWS "github.com/gorilla/websocket"

	"{{.RootPackage}}/internal/{{.ServiceName}}/router"
	"{{.RootPackage}}/internal/pkg/bootstrap"
	"{{.RootPackage}}/internal/pkg/config"
	"{{.RootPackage}}/internal/pkg/facade"
)

// initWebSocket initializes the WebSocket engine and hub.
func initWebSocket() (*gin.Engine, *websocket.Hub) {
	// Create hub
	hub := websocket.NewHub()

	// Create router and register handlers
	wsRouter := websocket.NewRouter()
	router.RegisterWSHandlers(wsRouter)

	// Create Gin engine for WebSocket
	engine := bootstrap.InitGinForWebSocket()

	// Configure WebSocket upgrader
	upgrader := gorillaWS.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin:     checkOrigin(facade.Config.WebSocket),
	}

	// Register WebSocket route
	engine.GET("/ws", func(c *gin.Context) {
		serveWS(c, hub, wsRouter, upgrader)
	})

	return engine, hub
}

// checkOrigin returns an origin checker function based on config.
func checkOrigin(cfg *config.WebSocket) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		if cfg == nil || cfg.AllowAllOrigins() {
			return true
		}

		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}

		return cfg.IsOriginAllowed(origin)
	}
}

// serveWS handles WebSocket upgrade requests.
func serveWS(c *gin.Context, hub *websocket.Hub, router *websocket.Router, upgrader gorillaWS.Upgrader) {
	ctx := context.Background()
	ctx = websocket.WithRequestID(ctx, c.GetHeader("X-Request-ID"))

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := websocket.NewClient(hub, conn, ctx,
		websocket.WithRouter(router),
		websocket.WithTokenParser(tokenParser),
		websocket.WithContextUpdater(contextUpdater),
	)

	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

// tokenParser parses JWT token and returns user info.
func tokenParser(tokenStr string) (*websocket.TokenInfo, error) {
	payload, err := token.Parse(tokenStr)
	if err != nil {
		return nil, err
	}

	return &websocket.TokenInfo{
		UserID:    payload.Subject,
		ExpiresAt: payload.ExpiresAt.Unix(),
	}, nil
}

// contextUpdater updates context with user ID after login.
func contextUpdater(ctx context.Context, userID string) context.Context {
	return websocket.WithUserID(ctx, userID)
}
