package v1

import (
	"time"
)

type {{.StructName}}Info struct {
	ID uint `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name  string    `json:"name"`
}

type List{{.StructName}}Request struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

type List{{.StructName}}Response struct {
	TotalCount int64       `json:"totalCount"`
	Data      []*{{.StructName}}Info `json:"data"`
}

type Create{{.StructName}}Request struct {
	Name string `json:"name" valid:"required,alphanum,stringlength(1|255)"`
}

type Get{{.StructName}}Response {{.StructName}}Info

type Update{{.StructName}}Request struct {
	Name *string `json:"name" valid:"required,stringlength(1|255)"`
}
