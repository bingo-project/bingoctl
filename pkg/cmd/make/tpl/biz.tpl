package {{.PackageName}}

import (
	"context"
	"regexp"

    "github.com/bingo-project/component-base/log"
	"github.com/jinzhu/copier"

	"{{.RootPackage}}/{{.StorePath}}"
	"{{.RootPackage}}/internal/pkg/errno"
	"{{.RootPackage}}/{{.ModelPath}}"
	v1 "{{.RootPackage}}/{{.RequestPath}}"
)

type {{.StructName}}Biz interface {
	List(ctx context.Context, offset, limit int) (*v1.List{{.StructName}}Response, error)
	Create(ctx context.Context, r *v1.Create{{.StructName}}Request) (*v1.Get{{.StructName}}Response, error)
	Get(ctx context.Context, ID uint) (*v1.Get{{.StructName}}Response, error)
	Update(ctx context.Context, ID uint, r *v1.Update{{.StructName}}Request) (*v1.Get{{.StructName}}Response, error)
	Delete(ctx context.Context, ID uint) error
}

type {{.VariableName}}Biz struct {
	ds store.IStore
}

// 确保 {{.VariableName}}Biz 实现了 {{.StructName}}Biz 接口.
var _ {{.StructName}}Biz = (*{{.VariableName}}Biz)(nil)

func New{{.StructName}}(ds store.IStore) *{{.VariableName}}Biz {
	return &{{.VariableName}}Biz{ds: ds}
}

func (b *{{.VariableName}}Biz) List(ctx context.Context, offset, limit int) (*v1.List{{.StructName}}Response, error) {
	count, list, err := b.ds.{{.StructNamePlural}}().List(ctx, offset, limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list {{.VariableNamePlural}}", "err", err)

		return nil, err
	}

	{{.VariableNamePlural}} := make([]*v1.{{.StructName}}Info, 0, len(list))
	for _, item := range list {
		var {{.VariableName}} v1.{{.StructName}}Info
		_ = copier.Copy(&{{.VariableName}}, item)

		{{.VariableNamePlural}} = append({{.VariableNamePlural}}, &{{.VariableName}})
	}

	return &v1.List{{.StructName}}Response{TotalCount: count, Data: {{.VariableNamePlural}}}, nil
}

func (b *{{.VariableName}}Biz) Create(ctx context.Context, request *v1.Create{{.StructName}}Request) (*v1.Get{{.StructName}}Response, error) {
	var {{.VariableName}}M model.{{.StructName}}M
	_ = copier.Copy(&{{.VariableName}}M, request)

	err := b.ds.{{.StructNamePlural}}().Create(ctx, &{{.VariableName}}M)
	if err != nil {
		// Check exists
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key", err.Error()); match {
			return nil, errno.ErrResourceAlreadyExists
		}

		return nil, err
	}

	var resp v1.Get{{.StructName}}Response
	_ = copier.Copy(&resp, {{.VariableName}}M)

	return &resp, nil
}

func (b *{{.VariableName}}Biz) Get(ctx context.Context, ID uint) (*v1.Get{{.StructName}}Response, error) {
	{{.VariableName}}, err := b.ds.{{.StructNamePlural}}().Get(ctx, ID)
	if err != nil {
		return nil, errno.ErrResourceNotFound
	}

	var resp v1.Get{{.StructName}}Response
	_ = copier.Copy(&resp, {{.VariableName}})

	return &resp, nil
}

func (b *{{.VariableName}}Biz) Update(ctx context.Context, ID uint, request *v1.Update{{.StructName}}Request) (*v1.Get{{.StructName}}Response, error) {
	{{.VariableName}}M, err := b.ds.{{.StructNamePlural}}().Get(ctx, ID)
	if err != nil {
		return nil, errno.ErrResourceNotFound
	}

	// if request.Name != nil {
	// 	{{.VariableName}}M.Name = *request.Name
	// }

	if err := b.ds.{{.StructNamePlural}}().Update(ctx, {{.VariableName}}M); err != nil {
		return nil, err
	}

	var resp v1.Get{{.StructName}}Response
	_ = copier.Copy(&resp, request)

	return &resp, nil
}

func (b *{{.VariableName}}Biz) Delete(ctx context.Context, ID uint) error {
	return b.ds.{{.StructNamePlural}}().Delete(ctx, ID)
}
