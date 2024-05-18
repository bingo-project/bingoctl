package router

import (
	"google.golang.org/grpc"

	"{[.RootPackage]}/internal/apiserver/grpc/controller/v1/apiserver"
	v1 "{[.RootPackage]}/internal/apiserver/grpc/proto/v1/pb"
	"{[.RootPackage]}/internal/apiserver/store"
)

func GRPC(g *grpc.Server) {
	// ApiServer
	v1.RegisterApiServerServer(g, apiserver.New(store.S))
}
