package {{.PackageName}}

import (
	"gorm.io/gorm"
)

type {{.StructName}}M struct {
	gorm.Model

	Name string `gorm:"type:varchar(255);not null;default:''"`
}

func (u *{{.StructName}}M) TableName() string {
	return "{{.TableName}}"
}
