package router

import (
	"google.golang.org/grpc"
)

// InstallGRPCServices registers gRPC services.
func InstallGRPCServices(s *grpc.Server) {
	// Register your gRPC services here
	// Example:
	// pb.RegisterYourServiceServer(s, &service.YourServiceImpl{})
}