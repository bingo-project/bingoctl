package {{.PackageName}}

import (
	"gorm.io/gorm"
)

type {{.StructName}}M struct {
	gorm.Model

	{{.MainFields}}
}

func (*{{.StructName}}M) TableName() string {
	return "{{.TableName}}"
}
