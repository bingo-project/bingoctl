package bootstrap

import (
	"github.com/bingo-project/component-base/log"

	"{[.RootPackage]}/internal/pkg/facade"
)

func InitLog() {
	log.Init(facade.Config.Log)
}
