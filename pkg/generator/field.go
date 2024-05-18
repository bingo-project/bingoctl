package generator

import (
	"slices"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/copier"
	"gorm.io/gen"
	genField "gorm.io/gen/field"

	"github.com/bingo-project/bingoctl/pkg/config"
)

// Field user input structures
type Field struct {
	Name             string
	Type             string
	ColumnName       string
	ColumnComment    string
	MultilineComment bool
	Tag              genField.Tag
	GORMTag          genField.GormTag
	CustomGenType    string
	Relation         *genField.Relation
}

func (o *Options) ReadMetaFields() error {
	if len(o.MetaFields) > 0 {
		return nil
	}

	// Generate model from table.
	g := gen.NewGenerator(gen.Config{
		ModelPkgPath: o.Directory,

		// generate model global configuration
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	g.UseDB(config.DB)

	// Generate struct `StructName` based on table `Table`
	meta := g.GenerateModelAs(o.Table, o.StructName)
	if len(meta.Fields) == 0 {
		return nil
	}

	for _, item := range meta.Fields {
		var field Field
		_ = copier.Copy(&field, item)

		o.MetaFields = append(o.MetaFields, &field)
	}

	return nil
}

func (o *Options) GetFieldsFromDB() error {
	err := o.ReadMetaFields()
	if err != nil {
		return err
	}

	o.Fields = ""
	o.MainFields = ""
	o.UpdatableFields = ""
	gormFields := []string{"ID", "CreatedAt", "UpdatedAt", "DeletedAt"}
	for _, field := range o.MetaFields {
		// Comment
		var comment string
		if field.ColumnComment != "" {
			comment = " // " + field.ColumnComment
		}

		// Type
		if field.Type == "gorm.DeletedAt" && o.Name == string(TmplRequest) {
			field.Type = "*time.Time"
		}

		// Replaces
		replaces := make(map[string]string)
		replaces["{{.Name}}"] = field.Name
		replaces["{{.Type}}"] = field.Type
		replaces["{{.GORMTag}}"] = field.GORMTag.Build()
		replaces["{{.JsonTag}}"] = strcase.ToLowerCamel(field.Tag["json"])
		replaces["{{.Comment}}"] = comment

		// Replace
		fieldTemplate := o.FieldTemplate
		for search, replace := range replaces {
			fieldTemplate = strings.ReplaceAll(fieldTemplate, search, replace)
		}

		// Fields
		o.Fields += fieldTemplate + "\n"

		// Skip gorm.Model fields.
		if slices.Contains(gormFields, field.Name) {
			continue
		}

		o.MainFields += fieldTemplate + "\n"

		// Updatable
		fieldTemplateUpdatable := o.FieldTemplate
		replaces["{{.NameSnake}}"] = strcase.ToSnake(field.Name)
		replaces["{{.VariableName}}"] = o.VariableName
		for search, replace := range replaces {
			if search == "{{.Type}}" && !strings.Contains(replace, "*") {
				replace = "*" + replace
			}

			fieldTemplateUpdatable = strings.ReplaceAll(fieldTemplateUpdatable, search, replace)
		}

		o.UpdatableFields += fieldTemplateUpdatable + "\n"
	}

	// Trim space
	o.Fields = strings.TrimSpace(o.Fields)
	o.MainFields = strings.TrimSpace(o.MainFields)
	o.UpdatableFields = strings.TrimSpace(o.UpdatableFields)

	return nil
}
