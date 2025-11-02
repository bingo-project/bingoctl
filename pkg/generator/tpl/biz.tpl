package {{.PackageName}}

import (
	"context"
	"regexp"

	"github.com/bingo-project/component-base/log"
	"github.com/jinzhu/copier"

	"{{.RootPackage}}/{{.StorePath}}"
	"{{.RootPackage}}/internal/pkg/errno"
	model "{{.RootPackage}}/{{.ModelPath}}{{.RelativePath}}"
	v1 "{{.RootPackage}}/{{.RequestPath}}{{.RelativePath}}"
)

type {{.StructName}}Biz interface {
	List(ctx context.Context, req *v1.List{{.StructName}}Request) (*v1.List{{.StructName}}Response, error)
	Create(ctx context.Context, req *v1.Create{{.StructName}}Request) (*v1.{{.StructName}}Info, error)
	Get(ctx context.Context, ID uint) (*v1.{{.StructName}}Info, error)
	Update(ctx context.Context, ID uint, req *v1.Update{{.StructName}}Request) (*v1.{{.StructName}}Info, error)
	Delete(ctx context.Context, ID uint) error

	{{.StructName}}Expansion
}

type {{.StructName}}Expansion interface{
}

type {{.VariableName}}Biz struct {
	ds store.IStore
}

var _ {{.StructName}}Biz = (*{{.VariableName}}Biz)(nil)

func New{{.StructName}}(ds store.IStore) *{{.VariableName}}Biz {
	return &{{.VariableName}}Biz{ds: ds}
}

func (b *{{.VariableName}}Biz) List(ctx context.Context, req *v1.List{{.StructName}}Request) (*v1.List{{.StructName}}Response, error) {
	count, list, err := b.ds.{{.StructName}}().List(ctx, req)
	if err != nil {
		log.C(ctx).Errorw("Failed to list {{.VariableName}}", "err", err)

		return nil, err
	}

	data := make([]v1.{{.StructName}}Info, 0)
	for _, item := range list {
		var {{.VariableName}} v1.{{.StructName}}Info
		_ = copier.Copy(&{{.VariableName}}, item)

		data = append(data, {{.VariableName}})
	}

	return &v1.List{{.StructName}}Response{Total: count, Data: data}, nil
}

func (b *{{.VariableName}}Biz) Create(ctx context.Context, req *v1.Create{{.StructName}}Request) (*v1.{{.StructName}}Info, error) {
	var {{.VariableName}}M model.{{.StructName}}
	_ = copier.Copy(&{{.VariableName}}M, req)

	err := b.ds.{{.StructName}}().Create(ctx, &{{.VariableName}}M)
	if err != nil {
		// Check exists
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key", err.Error()); match {
			return nil, errno.ErrResourceAlreadyExists
		}

		return nil, err
	}

	var resp v1.{{.StructName}}Info
	_ = copier.Copy(&resp, {{.VariableName}}M)

	return &resp, nil
}

func (b *{{.VariableName}}Biz) Get(ctx context.Context, ID uint) (*v1.{{.StructName}}Info, error) {
	{{.VariableName}}, err := b.ds.{{.StructName}}().Get(ctx, ID)
	if err != nil {
		return nil, errno.ErrResourceNotFound
	}

	var resp v1.{{.StructName}}Info
	_ = copier.Copy(&resp, {{.VariableName}})

	return &resp, nil
}

func (b *{{.VariableName}}Biz) Update(ctx context.Context, ID uint, req *v1.Update{{.StructName}}Request) (*v1.{{.StructName}}Info, error) {
	{{.VariableName}}M, err := b.ds.{{.StructName}}().Get(ctx, ID)
	if err != nil {
		return nil, errno.ErrResourceNotFound
	}

    {{.UpdatableFields}}

	if err := b.ds.{{.StructName}}().Update(ctx, {{.VariableName}}M); err != nil {
		return nil, err
	}

	var resp v1.{{.StructName}}Info
	_ = copier.Copy(&resp, {{.VariableName}}M)

	return &resp, nil
}

func (b *{{.VariableName}}Biz) Delete(ctx context.Context, ID uint) error {
	return b.ds.{{.StructName}}().Delete(ctx, ID)
}
