package bootstrap

import (
	"github.com/bingo-project/component-base/crypt"

	"{[.RootPackage]}/internal/pkg/facade"
)

func InitAES() {
	facade.AES = crypt.NewAES(facade.Config.Server.Key)
}
