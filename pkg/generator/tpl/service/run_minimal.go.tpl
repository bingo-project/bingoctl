package {{.ServiceName}}

import (
	"github.com/bingo-project/component-base/log"
)

// run 函数是实际的业务代码入口函数.
func run() error {
	log.Infow("{{.ServiceName}} service started")

	// TODO: Add your service logic here

	return nil
}