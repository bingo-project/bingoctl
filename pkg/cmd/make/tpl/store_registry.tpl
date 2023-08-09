func (ds *datastore) {{.StructNamePlural}}() {{.Package}}{{.StructName}}Store {
	return {{.Package}}New{{.StructNamePlural}}(ds.db)
}
