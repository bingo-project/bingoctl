package bootstrap

import (
	"{[.RootPackage]}/internal/apiserver/facade"
	genericapiserver "{[.RootPackage]}/internal/pkg/server"
)

var CfgFile string

const (
	// DefaultConfigName 指定了服务的默认配置文件名.
	DefaultConfigName = "{[.AppName]}-apiserver.yaml"
)

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	genericapiserver.LoadConfig(CfgFile, DefaultConfigName, &facade.Config, Boot)
}
