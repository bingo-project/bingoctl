package config

import (
	"gorm.io/gorm"

	"github.com/bingo-project/bingoctl/pkg/db"
)

var (
	Cfg *Config
	DB  *gorm.DB
)

type Config struct {
	Version      string           `mapstructure:"version" json:"version" yaml:"version"`
	RootPackage  string           `mapstructure:"rootPackage" json:"rootPackage" yaml:"rootPackage"`
	Directory    Directory        `mapstructure:"directory" json:"directory" yaml:"directory"`
	MysqlOptions *db.MySQLOptions `mapstructure:"mysql" json:"mysql" yaml:"mysql"`

	Registries Registries `mapstructure:"registries" json:"registries" yaml:"registries"`
}

type Directory struct {
	CMD        string `mapstructure:"cmd" json:"cmd" yaml:"cmd"`
	Model      string `mapstructure:"model" json:"model" yaml:"model"`
	Store      string `mapstructure:"store" json:"store" yaml:"store"`
	Request    string `mapstructure:"request" json:"request" yaml:"request"`
	Biz        string `mapstructure:"biz" json:"biz" yaml:"biz"`
	Controller string `mapstructure:"controller" json:"controller" yaml:"controller"`
	Middleware string `mapstructure:"middleware" json:"middleware" yaml:"middleware"`
	Job        string `mapstructure:"job" json:"job" yaml:"job"`
	Migration  string `mapstructure:"migration" json:"migration" yaml:"migration"`
	Seeder     string `mapstructure:"seeder" json:"seeder" yaml:"seeder"`
}

type Registries struct {
	Router string   `mapstructure:"router" json:"router" yaml:"router"`
	Store  Registry `mapstructure:"store" json:"store" yaml:"store"`
	Biz    Registry `mapstructure:"biz" json:"biz" yaml:"biz"`
}

type Registry struct {
	Filepath  string `mapstructure:"filepath" json:"filepath" yaml:"filepath"`
	Interface string `mapstructure:"interface" json:"interface" yaml:"interface"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Version:     "v1",
		RootPackage: "bingoctl",
		Directory: Directory{
			CMD:        "internal/bingoctl/cmd",
			Model:      "internal/apiserver/model",
			Store:      "internal/apiserver/store",
			Biz:        "internal/apiserver/biz",
			Controller: "internal/apiserver/http/controller/v1",
			Middleware: "internal/apiserver/http/middleware",
			Request:    "internal/apiserver/http/request/v1",
			Job:        "internal/watcher/watcher",
			Migration:  "internal/apiserver/database/migration",
			Seeder:     "internal/apiserver/database/seeder",
		},
	}
}
