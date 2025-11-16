package {{.ServiceName}}

import (
	"net"

	"github.com/bingo-project/component-base/log"
	"google.golang.org/grpc"

	genericserver "{{.RootPackage}}/internal/pkg/server"
)

// GRPCServer represents the gRPC server.
type GRPCServer struct {
	*grpc.Server
	address string
}

// NewGRPC creates a new gRPC server instance.
func NewGRPC() *GRPCServer {
	// Create gRPC server with options.
	grpcServer := grpc.NewServer()

	// Register your gRPC services here.
	// Example:
	// pb.RegisterYourServiceServer(grpcServer, &yourServiceImpl{})

	return &GRPCServer{
		Server:  grpcServer,
		address: genericserver.Config.GRPCServer.Addr,
	}
}

// Run starts the gRPC server.
func (s *GRPCServer) Run() {
	listen, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatalw("Failed to listen", "err", err)
	}

	log.Infow("Start to listening the incoming requests on grpc address", "addr", s.address)

	go func() {
		if err := s.Serve(listen); err != nil {
			log.Fatalw("Failed to start grpc server", "err", err)
		}
	}()
}

// Close gracefully shuts down the gRPC server.
func (s *GRPCServer) Close() {
	s.GracefulStop()
	log.Infow("gRPC server stopped")
}