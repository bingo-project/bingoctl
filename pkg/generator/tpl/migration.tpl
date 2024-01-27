package {{.PackageName}}

import (
	"gorm.io/gorm"

	"github.com/bingo-project/bingoctl/pkg/migrate"
)

type {{.StructName}} struct {
	gorm.Model
}

func ({{.StructName}}) TableName() string {
	return "{{.TableName}}"
}

func ({{.StructName}}) Up(migrator gorm.Migrator) {
	_ = migrator.AutoMigrate(&{{.StructName}}{})
}

func ({{.StructName}}) Down(migrator gorm.Migrator) {
	_ = migrator.DropTable(&{{.StructName}}{})
}

func init() {
	migrate.Add("{{.TimeStr}}_{{.VariableNameSnake}}", {{.StructName}}{}.Up, {{.StructName}}{}.Down)
}
