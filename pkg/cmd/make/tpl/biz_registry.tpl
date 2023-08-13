func (b *biz) {{.StructNamePlural}}() {{.Package}}{{.StructName}}Biz {
	return {{.Package}}New{{.StructName}}(b.ds)
}
