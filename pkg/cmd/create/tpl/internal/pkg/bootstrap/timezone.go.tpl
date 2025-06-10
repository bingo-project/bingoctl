package bootstrap

import (
	"{[.RootPackage]}/internal/pkg/facade"
)

func InitTimezone() {
	facade.Config.Server.SetTimezone()
}
