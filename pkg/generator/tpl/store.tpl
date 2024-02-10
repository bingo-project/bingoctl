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
	List(ctx context.Context, req *v1.List{{.StructName}}Request) (int64, []*model.{{.StructName}}, error)
	Create(ctx context.Context, {{.VariableName}} *model.{{.StructName}}) error
	Get(ctx context.Context, ID uint) (*model.{{.StructName}}, error)
	Update(ctx context.Context, {{.VariableName}} *model.{{.StructName}}, fields ...string) error
	Delete(ctx context.Context, ID uint) error

	CreateInBatch(ctx context.Context, {{.VariableNamePlural}} []*model.{{.StructName}}) error
	CreateIfNotExist(ctx context.Context, {{.VariableName}} *model.{{.StructName}}) error
	FirstOrCreate(ctx context.Context, where any, {{.VariableName}} *model.{{.StructName}}) error
	UpdateOrCreate(ctx context.Context, where any, {{.VariableName}} *model.{{.StructName}}) error
	Upsert(ctx context.Context, {{.VariableName}} *model.{{.StructName}}, fields ...string) error
	DeleteInBatch(ctx context.Context, ids []uint) error
}

type {{.VariableNamePlural}} struct {
	db *gorm.DB
}

var _ {{.StructName}}Store = (*{{.VariableNamePlural}})(nil)

func New{{.StructNamePlural}}(db *gorm.DB) *{{.VariableNamePlural}} {
	return &{{.VariableNamePlural}}{db: db}
}

func Search{{.StructName}}(req *v1.List{{.StructName}}Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// if req.Name != "" {
		// 	db.Where("name = ?", req.Name)
		// }

		return db
	}
}

func (s *{{.VariableNamePlural}}) List(ctx context.Context, req *v1.List{{.StructName}}Request) (count int64, ret []*model.{{.StructName}}, err error) {
	db := s.db.Scopes(Search{{.StructName}}(req))
	count, err = gormutil.Paginate(db, &req.ListOptions, &ret)

	return
}

func (s *{{.VariableNamePlural}}) Create(ctx context.Context, {{.VariableName}} *model.{{.StructName}}) error {
	return s.db.Create(&{{.VariableName}}).Error
}

func (s *{{.VariableNamePlural}}) Get(ctx context.Context, ID uint) ({{.VariableName}} *model.{{.StructName}}, err error) {
	err = s.db.Where("id = ?", ID).First(&{{.VariableName}}).Error

	return
}

func (s *{{.VariableNamePlural}}) Update(ctx context.Context, {{.VariableName}} *model.{{.StructName}}, fields ...string) error {
	return s.db.Select(fields).Save(&{{.VariableName}}).Error
}

func (s *{{.VariableNamePlural}}) Delete(ctx context.Context, ID uint) error {
	return s.db.Where("id = ?", ID).Delete(&model.{{.StructName}}{}).Error
}

func (s *{{.VariableNamePlural}}) CreateInBatch(ctx context.Context, {{.VariableNamePlural}} []*model.{{.StructName}}) error {
	return s.db.CreateInBatches(&{{.VariableNamePlural}}, global.CreateBatchSize).Error
}

func (s *{{.VariableNamePlural}}) CreateIfNotExist(ctx context.Context, {{.VariableName}} *model.{{.StructName}}) error {
	return s.db.Clauses(clause.OnConflict{DoNothing: true}).
		Create(&{{.VariableName}}).
		Error
}

func (s *{{.VariableNamePlural}}) FirstOrCreate(ctx context.Context, where any, {{.VariableName}} *model.{{.StructName}}) error {
	return s.db.Where(where).
		Attrs(&{{.VariableName}}).
		FirstOrCreate(&{{.VariableName}}).
		Error
}

func (s *{{.VariableNamePlural}}) UpdateOrCreate(ctx context.Context, where any, {{.VariableName}} *model.{{.StructName}}) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var exist model.{{.StructName}}
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(where).
			First(&exist).
			Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		{{.VariableName}}.ID = exist.ID

		return tx.Omit("CreatedAt").Save(&{{.VariableName}}).Error
	})
}

func (s *{{.VariableNamePlural}}) Upsert(ctx context.Context, {{.VariableName}} *model.{{.StructName}}, fields ...string) error {
	do := clause.OnConflict{UpdateAll: true}
	if len(fields) > 0 {
		do.UpdateAll = false
		do.DoUpdates = clause.AssignmentColumns(fields)
	}

	return s.db.Clauses(do).
		Create(&{{.VariableName}}).
		Error
}

func (s *{{.VariableNamePlural}}) DeleteInBatch(ctx context.Context, ids []uint) error {
	return s.db.Where("id IN (?)", ids).
		Delete(&model.{{.StructName}}{}).
		Error
}
