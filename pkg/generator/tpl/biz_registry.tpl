func (b *biz) {{.StructName}}V1() {{.Package}}{{.StructName}}Biz {
	return {{.Package}}New{{.StructName}}(b.store)
}
