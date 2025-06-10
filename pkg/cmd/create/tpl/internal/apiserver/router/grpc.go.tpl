package router

import (
	"google.golang.org/grpc"

	"{[.RootPackage]}/internal/apiserver/grpc/v1/apiserver"
	"{[.RootPackage]}/internal/apiserver/store"
	v1 "{[.RootPackage]}/pkg/proto/v1/pb"
)

func GRPC(g *grpc.Server) {
	// ApiServer
	v1.RegisterApiServerServer(g, apiserver.New(store.S))
}
