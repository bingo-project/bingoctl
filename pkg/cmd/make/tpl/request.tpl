package v1

import (
	"time"

	"github.com/bingo-project/component-base/util/gormutil"
)

type {{.StructName}}Info struct {
	ID uint `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name  string    `json:"name"`
}

type List{{.StructName}}Request struct {
	gormutil.ListOptions
}

type List{{.StructName}}Response struct {
	Total int64 `json:"total"`
	Data  []{{.StructName}}Info   `json:"data"`
}

type Create{{.StructName}}Request struct {
	Name string `json:"name" valid:"required,alphanum,stringlength(1|255)"`
}

type Update{{.StructName}}Request struct {
	Name *string `json:"name" valid:"required,stringlength(1|255)"`
}
