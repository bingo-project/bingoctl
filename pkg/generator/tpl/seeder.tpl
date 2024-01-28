package {{.PackageName}}

type {{.StructName}} struct {
}

// Signature The name and signature of the seeder.
func ({{.StructName}}) Signature() string {
	return "{{.StructName}}"
}

// Run seed the application's database.
func ({{.StructName}}) Run() error {
	//

	return nil
}
