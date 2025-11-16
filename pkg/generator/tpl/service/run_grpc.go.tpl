package {{.ServiceName}}

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bingo-project/component-base/log"
)

// run 函数是实际的业务代码入口函数.
func run() error {
	// 启动 gRPC 服务
	grpcServer := NewGRPC()
	grpcServer.Run()

	// 等待中断信号优雅地关闭服务器。
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infow("Shutting down server ...")

	// 停止服务
	grpcServer.Close()

	return nil
}