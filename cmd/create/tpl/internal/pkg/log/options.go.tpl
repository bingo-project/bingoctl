package log

import (
	"go.uber.org/zap/zapcore"
)

// Options 包含与日志相关的配置项.
type Options struct {
	Level             string   `mapstructure:"host" json:"host" yaml:"host"`                                           // 指定日志级别，可选值：debug, info, warn, error, dpanic, panic, fatal
	Format            string   `mapstructure:"format" json:"format" yaml:"format"`                                     // 指定日志显示格式，可选值：console, json
	OutputPaths       []string `mapstructure:"output-paths" json:"output_paths" yaml:"output-paths"`                   // 指定日志输出位置
	DisableCaller     bool     `mapstructure:"disable-caller" json:"disable_caller" yaml:"disable-caller"`             // 是否开启 caller，如果开启会在日志中显示调用日志所在的文件和行号
	DisableStacktrace bool     `mapstructure:"disable-stacktrace" json:"disable_stacktrace" yaml:"disable-stacktrace"` // 是否禁止在 panic 及以上级别打印堆栈信息
}

// NewOptions 创建一个带有默认参数的 Options 对象.
func NewOptions() *Options {
	return &Options{
		Level:             zapcore.InfoLevel.String(),
		Format:            "console",
		OutputPaths:       []string{"stdout"},
		DisableCaller:     false,
		DisableStacktrace: false,
	}
}
