package v1

import (
	"time"

	"github.com/bingo-project/component-base/util/gormutil"
)

type {{.StructName}}Info struct {
	{{.Fields}}
}

type List{{.StructName}}Request struct {
	gormutil.ListOptions
}

type List{{.StructName}}Response struct {
	Total int64 `json:"total"`
	Data  []{{.StructName}}Info   `json:"data"`
}

type Create{{.StructName}}Request struct {
	{{.MainFields}}
}

type Update{{.StructName}}Request struct {
	{{.MainFields}}
}
