package bootstrap

import "{[.RootPackage]}/internal/apiserver/facade"

func InitTimezone() {
	facade.Config.Server.SetTimezone()
}
