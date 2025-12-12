package grpc

import (
	"{{.RootPackage}}/internal/{{.ServiceName}}/biz"
	"{{.RootPackage}}/internal/pkg/store"
)

type Handler struct {
	b biz.IBiz
	// Add UnimplementedXxxServer here after generating proto
}

func NewHandler(ds store.IStore) *Handler {
	return &Handler{b: biz.NewBiz(ds)}
}

// Example handler method
// func (h *Handler) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
//     // Implementation here
//     return &pb.GetResponse{}, nil
// }
