package {{.PackageName}}

import (
	"context"

	"github.com/bingo-project/component-base/util/gormutil"

	model "{{.RootPackage}}/{{.ModelPath}}{{.RelativePath}}"
	v1 "{{.RootPackage}}/{{.RequestPath}}{{.RelativePath}}"
	genericstore "{{.RootPackage}}/pkg/store"
	"{{.RootPackage}}/pkg/store/where"
)

type {{.StructName}}Store interface {
	Create(ctx context.Context, obj *model.{{.StructName}}M) error
	Update(ctx context.Context, obj *model.{{.StructName}}M, fields ...string) error
	Delete(ctx context.Context, opts *where.Options) error
	Get(ctx context.Context, opts *where.Options) (*model.{{.StructName}}M, error)
	List(ctx context.Context, opts *where.Options) (int64, []*model.{{.StructName}}M, error)

	{{.StructName}}Expansion
}

type {{.StructName}}Expansion interface {
	ListWithRequest(ctx context.Context, req *v1.List{{.StructName}}Request) (int64, []*model.{{.StructName}}M, error)
	GetByID(ctx context.Context, id uint) (*model.{{.StructName}}M, error)
	DeleteByID(ctx context.Context, id uint) error
}

type {{.VariableName}}Store struct {
	*genericstore.Store[model.{{.StructName}}M]
}

var _ {{.StructName}}Store = (*{{.VariableName}}Store)(nil)

func New{{.StructName}}Store(store *datastore) *{{.VariableName}}Store {
	return &{{.VariableName}}Store{
		Store: genericstore.NewStore[model.{{.StructName}}M](store, NewLogger()),
	}
}

func (s *{{.VariableName}}Store) ListWithRequest(ctx context.Context, req *v1.List{{.StructName}}Request) (int64, []*model.{{.StructName}}M, error) {
	opts := where.NewWhere()
	{{.UpdatableFields}}

	db := s.DB(ctx, opts)
	var ret []*model.{{.StructName}}M
	count, err := gormutil.Paginate(db, &req.ListOptions, &ret)

	return count, ret, err
}

func (s *{{.VariableName}}Store) GetByID(ctx context.Context, id uint) (*model.{{.StructName}}M, error) {
	return s.Get(ctx, where.F("id", id))
}

func (s *{{.VariableName}}Store) DeleteByID(ctx context.Context, id uint) error {
	return s.Delete(ctx, where.F("id", id))
}
