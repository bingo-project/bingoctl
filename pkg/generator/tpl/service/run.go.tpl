package {{.ServiceName}}

import (
	"context"
	"os/signal"
	"syscall"

	"{{.RootPackage}}/internal/pkg/facade"
	"{{.RootPackage}}/internal/pkg/server"
)

// run starts all enabled servers based on configuration.
func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
{{if .EnableHTTP}}
	ginEngine := initGinEngine()
{{- end}}
{{- if .EnableGRPC}}
	grpcServer := initGRPCServer(facade.Config.GRPC)
{{- end}}
{{- if .EnableWS}}
	wsEngine, wsHub := initWebSocket()
{{- end}}

	runner := server.Assemble(
		&facade.Config,
{{- if .EnableHTTP}}
		server.WithGinEngine(ginEngine),
{{- end}}
{{- if .EnableGRPC}}
		server.WithGRPCServer(grpcServer),
{{- end}}
{{- if .EnableWS}}
		server.WithWebSocket(wsEngine, wsHub),
{{- end}}
	)

	return runner.Run(ctx)
}
