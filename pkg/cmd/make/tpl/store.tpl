package {{.PackageName}}

import (
	"context"

	"gorm.io/gorm"

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
}

type {{.VariableNamePlural}} struct {
	db *gorm.DB
}

var _ {{.StructName}}Store = (*{{.VariableNamePlural}})(nil)

func New{{.StructNamePlural}}(db *gorm.DB) *{{.VariableNamePlural}} {
	return &{{.VariableNamePlural}}{db: db}
}

func (u *{{.VariableNamePlural}}) List(ctx context.Context, req *v1.List{{.StructName}}Request) (count int64, ret []*model.{{.StructName}}M, err error) {
	// Order
	if req.Order == "" {
		req.Order = "id"
	}

	// Sort
	if req.Sort == "" {
		req.Sort = "desc"
	}

	err = u.db.Offset(req.Offset).
		Limit(helper.DefaultLimit(req.Limit)).
		Order(req.Order + " " + req.Sort).
		Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error

	return
}

func (u *{{.VariableNamePlural}}) Create(ctx context.Context, {{.VariableName}} *model.{{.StructName}}M) error {
	return u.db.Create(&{{.VariableName}}).Error
}

func (u *{{.VariableNamePlural}}) Get(ctx context.Context, ID uint) ({{.VariableName}} *model.{{.StructName}}M, err error) {
	err = u.db.Where("id = ?", ID).First(&{{.VariableName}}).Error

	return
}

func (u *{{.VariableNamePlural}}) Update(ctx context.Context, {{.VariableName}} *model.{{.StructName}}M, fields ...string) error {
	return u.db.Select(fields).Save(&{{.VariableName}}).Error
}

func (u *{{.VariableNamePlural}}) Delete(ctx context.Context, ID uint) error {
	return u.db.Where("id = ?", ID).Delete(&model.{{.StructName}}M{}).Error
}
