package {{.ServiceName}}

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"{{.RootPackage}}/internal/pkg/config"
)

// initGRPCServer initializes the gRPC server with services.
func initGRPCServer(cfg *config.GRPC) *grpc.Server {
	opts := []grpc.ServerOption{
		// Add interceptors here
		// grpc.ChainUnaryInterceptor(...),
	}

	srv := grpc.NewServer(opts...)

	// Register gRPC services here
	// Example:
	// router.GRPC(srv)

	// Enable reflection for grpcurl debugging
	reflection.Register(srv)

	return srv
}
