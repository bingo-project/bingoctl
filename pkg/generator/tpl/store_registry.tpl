func (store *datastore) {{.StructName}}() {{.Package}}{{.StructName}}Store {
	return {{.Package}}New{{.StructName}}Store(store)
}
