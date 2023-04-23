package {{.PackageName}}

import (
	"gorm.io/gorm"
)

type {{.StructName}}M struct {
	gorm.Model

	Name string `gorm:"column:name;not null;default:''" json:"name"`
}

func (u *{{.StructName}}M) TableName() string {
	return "{{.TableName}}"
}
