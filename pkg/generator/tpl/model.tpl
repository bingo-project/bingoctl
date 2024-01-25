package {{.PackageName}}

import (
	"gorm.io/gorm"
)

type {{.StructName}} struct {
	gorm.Model

	{{.MainFields}}
}

func (*{{.StructName}}) TableName() string {
	return "{{.TableName}}"
}
