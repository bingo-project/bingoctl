package apiserver

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bingo-project/component-base/log"
	"github.com/hibiken/asynq"

	"{[.RootPackage]}/internal/apiserver/bootstrap"
	"{[.RootPackage]}/internal/apiserver/facade"
	"{[.RootPackage]}/internal/apiserver/job"
	"{[.RootPackage]}/pkg/queue"
)

// run 函数是实际的业务代码入口函数.
// kill 默认会发送 syscall.SIGTERM 信号
// kill -2 发送 syscall.SIGINT 信号，我们常用的 CTRL + C 就是触发系统 SIGINT 信号
// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它.
func run() error {
	bootstrap.Boot()

	// 启动 http 服务
	httpServer := NewHttp()
	httpServer.Run()

	// 启动 grpc 服务
	grpcServer := NewGRPC()
	grpcServer.Run()

	// 启动 queue worker.
	go runJobs()

	// 等待中断信号优雅地关闭服务器（10 秒超时)。
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infow("Shutting down server ...")

	// 停止服务
	httpServer.Close()
	grpcServer.Close()

	return nil
}

func runJobs() {
	mux := asynq.NewServeMux()
	mux.Use(queue.Logging)

	job.Register(mux)

	err := facade.Worker.Run(mux)
	if err != nil {
		log.Fatalw("run worker failed", "err", err)
	}
}
