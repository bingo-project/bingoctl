func (b *biz) {{.StructName}}() {{.Package}}{{.StructName}}Biz {
	return {{.Package}}New{{.StructName}}(b.ds)
}
