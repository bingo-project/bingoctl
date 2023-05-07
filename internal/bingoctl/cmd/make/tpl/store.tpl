package {{.PackageName}}

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"{{.RootPackage}}/{{.ModelPath}}"
)

type {{.StructName}}Store interface {
	List(ctx context.Context, offset, limit int) (int64, []*model.{{.StructName}}M, error)
	Create(ctx context.Context, {{.VariableName}} *model.{{.StructName}}M) error
	Get(ctx context.Context, ID uint) (*model.{{.StructName}}M, error)
	Update(ctx context.Context, {{.VariableName}} *model.{{.StructName}}M) error
	Delete(ctx context.Context, ID uint) error
}

type {{.VariableNamePlural}} struct {
	db *gorm.DB
}

// 确保 {{.VariableNamePlural}} 实现了 {{.StructName}}Store 接口.
var _ {{.StructName}}Store = (*{{.VariableNamePlural}})(nil)

func new{{.StructNamePlural}}(db *gorm.DB) *{{.VariableNamePlural}} {
	return &{{.VariableNamePlural}}{db: db}
}

func (u *{{.VariableNamePlural}}) Create(ctx context.Context, {{.VariableName}} *model.{{.StructName}}M) error {
	return u.db.Create(&{{.VariableName}}).Error
}

func (u *{{.VariableNamePlural}}) Get(ctx context.Context, ID uint) ({{.VariableName}} *model.{{.StructName}}M, err error) {
	err = u.db.Where("id = ?", ID).First(&{{.VariableName}}).Error

	return
}

func (u *{{.VariableNamePlural}}) Update(ctx context.Context, {{.VariableName}} *model.{{.StructName}}M) error {
	return u.db.Save(&{{.VariableName}}).Error
}

func (u *{{.VariableNamePlural}}) List(ctx context.Context, offset, limit int) (count int64, ret []*model.{{.StructName}}M, err error) {
	err = u.db.Offset(offset).Limit(limit).Order("id desc").Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error

	return
}

func (u *{{.VariableNamePlural}}) Delete(ctx context.Context, ID uint) error {
	err := u.db.Where("id = ?", ID).Delete(&model.{{.StructName}}M{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}
