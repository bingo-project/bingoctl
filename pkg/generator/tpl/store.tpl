package {{.PackageName}}

import (
	"context"
	"errors"

	"github.com/bingo-project/component-base/util/gormutil"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"{{.RootPackage}}/internal/apiserver/global"
	v1 "{{.RootPackage}}/{{.RequestPath}}{{.RelativePath}}"
	model "{{.RootPackage}}/{{.ModelPath}}{{.RelativePath}}"
)

type {{.StructName}}Store interface {
	List(ctx context.Context, req *v1.List{{.StructName}}Request) (int64, []*model.{{.StructName}}M, error)
	Create(ctx context.Context, obj *model.{{.StructName}}M) error
	Get(ctx context.Context, ID uint) (*model.{{.StructName}}M, error)
	Update(ctx context.Context, obj *model.{{.StructName}}M, fields ...string) error
	Delete(ctx context.Context, ID uint) error

	{{.StructName}}Expansion
}

type {{.StructName}}Expansion interface{
	CreateInBatch(ctx context.Context, obj []*model.{{.StructName}}M) error
	CreateIfNotExist(ctx context.Context, obj *model.{{.StructName}}M) error
	FirstOrCreate(ctx context.Context, where any, obj *model.{{.StructName}}M) error
	UpdateOrCreate(ctx context.Context, where any, obj *model.{{.StructName}}M) error
	Upsert(ctx context.Context, obj *model.{{.StructName}}M, fields ...string) error
	DeleteInBatch(ctx context.Context, ids []uint) error
}

type {{.VariableName}}Store struct {
	store *datastore
}

var _ {{.StructName}}Store = (*{{.VariableName}}Store)(nil)

func New{{.StructName}}Store (store *datastore) *{{.VariableName}}Store {
	return &{{.VariableName}}Store{store: store}
}

func Search{{.StructName}}(req *v1.List{{.StructName}}Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
        {{.UpdatableFields}}

		return db
	}
}

func (s *{{.VariableName}}Store) List(ctx context.Context, req *v1.List{{.StructName}}Request) (count int64, ret []*model.{{.StructName}}, err error) {
	db := s.db.WithContext(ctx).Scopes(Search{{.StructName}}(req))
	count, err = gormutil.Paginate(db, &req.ListOptions, &ret)

	return
}

func (s *{{.VariableName}}Store) Create(ctx context.Context, obj *model.{{.StructName}}M) error {
	return s.db.WithContext(ctx).Create(&obj).Error
}

func (s *{{.VariableName}}Store) Get(ctx context.Context, ID uint) (obj *model.{{.StructName}}M, err error) {
	err = s.db.WithContext(ctx).Where("id = ?", ID).First(&obj).Error

	return
}

func (s *{{.VariableName}}Store) Update(ctx context.Context, obj *model.{{.StructName}}M, fields ...string) error {
	return s.db.WithContext(ctx).Select(fields).Save(&obj).Error
}

func (s *{{.VariableName}}Store) Delete(ctx context.Context, ID uint) error {
	return s.db.WithContext(ctx).Where("id = ?", ID).Delete(&model.{{.StructName}}{}).Error
}

func (s *{{.VariableName}}Store) CreateInBatch(ctx context.Context, obj []*model.{{.StructName}}M) error {
	return s.db.WithContext(ctx).CreateInBatches(&obj, global.CreateBatchSize).Error
}

func (s *{{.VariableName}}Store) CreateIfNotExist(ctx context.Context, obj *model.{{.StructName}}M) error {
	return s.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&obj).
		Error
}

func (s *{{.VariableName}}Store) FirstOrCreate(ctx context.Context, where any, obj *model.{{.StructName}}M) error {
	return s.db.WithContext(ctx).
		Where(where).
		Attrs(&obj).
		FirstOrCreate(&obj).
		Error
}

func (s *{{.VariableName}}Store) UpdateOrCreate(ctx context.Context, where any, obj *model.{{.StructName}}M) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var exist model.{{.StructName}}
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(where).
			First(&exist).
			Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		obj.ID = exist.ID

		return tx.Omit("CreatedAt").Save(&obj).Error
	})
}

func (s *{{.VariableName}}Store) Upsert(ctx context.Context, obj *model.{{.StructName}}M, fields ...string) error {
	do := clause.OnConflict{UpdateAll: true}
	if len(fields) > 0 {
		do.UpdateAll = false
		do.DoUpdates = clause.AssignmentColumns(fields)
	}

	return s.db.WithContext(ctx).
		Clauses(do).
		Create(&obj).
		Error
}

func (s *{{.VariableName}}Store) DeleteInBatch(ctx context.Context, ids []uint) error {
	return s.db.WithContext(ctx).
		Where("id IN (?)", ids).
		Delete(&model.{{.StructName}}{}).
		Error
}
