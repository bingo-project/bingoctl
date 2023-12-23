package {{.PackageName}}

import (
	"context"
	"errors"

	"github.com/bingo-project/component-base/util/gormutil"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"{{.RootPackage}}/internal/apiserver/global"
	"{{.RootPackage}}/{{.ModelPath}}"
	"{{.RootPackage}}/internal/pkg/util/helper"
	v1 "{{.RootPackage}}/{{.RequestPath}}"
)

type {{.StructName}}Store interface {
	List(ctx context.Context, req *v1.List{{.StructName}}Request) (int64, []*model.{{.StructName}}M, error)
	Create(ctx context.Context, {{.VariableName}} *model.{{.StructName}}M) error
	Get(ctx context.Context, ID uint) (*model.{{.StructName}}M, error)
	Update(ctx context.Context, {{.VariableName}} *model.{{.StructName}}M, fields ...string) error
	Delete(ctx context.Context, ID uint) error

	CreateInBatch(ctx context.Context, {{.VariableNamePlural}} []*model.{{.StructName}}M) error
	FirstOrCreate(ctx context.Context, where any, {{.VariableName}} *model.{{.StructName}}M) error
	UpdateOrCreate(ctx context.Context, where any, {{.VariableName}} *model.{{.StructName}}M) error
}

type {{.VariableNamePlural}} struct {
	db *gorm.DB
}

var _ {{.StructName}}Store = (*{{.VariableNamePlural}})(nil)

func New{{.StructNamePlural}}(db *gorm.DB) *{{.VariableNamePlural}} {
	return &{{.VariableNamePlural}}{db: db}
}

func (s *{{.VariableNamePlural}}) List(ctx context.Context, req *v1.List{{.StructName}}Request) (count int64, ret []*model.{{.StructName}}M, err error) {
	count, err = gormutil.Paginate(u.db, &req.ListOptions, &ret)

	return
}

func (s *{{.VariableNamePlural}}) Create(ctx context.Context, {{.VariableName}} *model.{{.StructName}}M) error {
	return s.db.Create(&{{.VariableName}}).Error
}

func (s *{{.VariableNamePlural}}) Get(ctx context.Context, ID uint) ({{.VariableName}} *model.{{.StructName}}M, err error) {
	err = s.db.Where("id = ?", ID).First(&{{.VariableName}}).Error

	return
}

func (s *{{.VariableNamePlural}}) Update(ctx context.Context, {{.VariableName}} *model.{{.StructName}}M, fields ...string) error {
	return s.db.Select(fields).Save(&{{.VariableName}}).Error
}

func (s *{{.VariableNamePlural}}) Delete(ctx context.Context, ID uint) error {
	return s.db.Where("id = ?", ID).Delete(&model.{{.StructName}}M{}).Error
}

func (s *{{.VariableNamePlural}}) CreateInBatch(ctx context.Context, {{.VariableNamePlural}} []*model.{{.StructName}}M) error {
	return s.db.CreateInBatches(&{{.VariableNamePlural}}, global.CreateBatchSize).Error
}

func (s *{{.VariableNamePlural}}) FirstOrCreate(ctx context.Context, where any, {{.VariableName}} *model.{{.StructName}}M) error {
	return s.db.Where(where).
		Attrs(&{{.VariableName}}).
		FirstOrCreate(&{{.VariableName}}).
		Error
}

func (s *{{.VariableNamePlural}}) UpdateOrCreate(ctx context.Context, where any, {{.VariableName}} *model.{{.StructName}}M) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var exist model.{{.StructName}}M
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(where).
			First(&exist).
			Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		{{.VariableName}}.ID = exist.ID

		return tx.Save(&{{.VariableName}}).Error
	})
}