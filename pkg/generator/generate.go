package generator

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"

	"github.com/bingo-project/bingoctl/pkg/config"
	cmdutil "github.com/bingo-project/bingoctl/pkg/util"
)

func (o *Options) GenerateCode(tmpl, path string) error {
	dir := GetMapDirectory(tmpl)

	o.SetName(tmpl)
	o.ReadCodeTemplates()
	o.GenerateAttributes(dir, path)

	// Generate from db table.
	dbTemplates := []Tmpl{TmplModel, TmplRequest, TmplStore, TmplBiz}
	if slices.Contains(dbTemplates, Tmpl(o.Name)) && o.Table != "" {
		_ = o.GetFieldsFromDB()
	}

	err := cmdutil.GenerateCode(o.FilePath, o.CodeTemplate, o.Name, o)
	if err != nil {
		return err
	}

	if o.Name == string(TmplStore) {
		err = o.Register(config.Cfg.Registries.Store, o.InterfaceTemplate, o.RegisterTemplate, o.RootPackage+"/"+o.StorePath+o.RelativePath)
		if err != nil {
			return err
		}
	}

	if o.Name == string(TmplBiz) {
		err = o.Register(config.Cfg.Registries.Biz, o.InterfaceTemplate, o.RegisterTemplate, o.RootPackage+"/"+o.BizPath+o.RelativePath)
		if err != nil {
			return err
		}
	}

	// Format code
	cmd := exec.Command("gofmt", "-w", o.FilePath)
	_ = cmd.Run()

	return nil
}

func GetMapDirectory(tmpl string) (dir string) {
	dir = config.Cfg.Directory.CMD
	if tmpl == string(TmplModel) {
		dir = config.Cfg.Directory.Model
	}
	if tmpl == string(TmplStore) {
		dir = config.Cfg.Directory.Store
	}
	if tmpl == string(TmplRequest) {
		dir = config.Cfg.Directory.Request
	}
	if tmpl == string(TmplBiz) {
		dir = config.Cfg.Directory.Biz
	}
	if tmpl == string(TmplController) {
		dir = config.Cfg.Directory.Controller
	}
	if tmpl == string(TmplMiddleware) {
		dir = config.Cfg.Directory.Middleware
	}
	if tmpl == string(TmplJob) {
		dir = config.Cfg.Directory.Job
	}
	if tmpl == string(TmplMigration) {
		dir = config.Cfg.Directory.Migration
	}
	if tmpl == string(TmplSeeder) {
		dir = config.Cfg.Directory.Seeder
	}

	return
}

func (o *Options) GenerateAttributes(directory string, path string) *Options {
	// Set code attributes
	o.RootPackage = config.Cfg.RootPackage
	o.BizPath = config.Cfg.Directory.Biz
	o.StorePath = config.Cfg.Directory.Store
	o.RequestPath = config.Cfg.Directory.Request
	o.ModelPath = config.Cfg.Directory.Model
	if filepath.Dir(path) != "." {
		o.RelativePath = "/" + filepath.Dir(path)
	}

	if o.Directory == "" {
		o.Directory = directory
	}

	arr := strings.Split(filepath.Join(o.Directory, path), "/")
	name := arr[len(arr)-1]

	o.StructName = strcase.ToCamel(name)
	o.StructNamePlural = pluralize.NewClient().Plural(o.StructName)
	o.VariableName = strcase.ToLowerCamel(o.StructName)
	o.VariableNameSnake = strcase.ToSnake(o.StructName)
	o.VariableNamePlural = pluralize.NewClient().Plural(o.VariableName)
	o.TableName = strcase.ToSnake(o.StructName)
	if o.Table != "" {
		o.TableName = o.Table
	}

	// Flags: Model name
	if o.ModelName == "" {
		o.ModelName = o.StructName
	}

	// Directory
	directoryArr := arr[:len(arr)-1]
	o.Directory = strings.Join(directoryArr, "/")
	o.Directory = o.Directory + "/"
	if o.PackageName == "" && len(directoryArr) > 0 {
		o.PackageName = strings.ToLower(directoryArr[len(directoryArr)-1])
	}

	// File path
	o.FilePath = filepath.Join(o.Directory, o.VariableNameSnake+".go")

	// Migration
	o.TimeStr = time.Now().Format("2006_01_02_150405")
	if o.Name == string(TmplMigration) {
		o.FilePath = filepath.Join(o.Directory, o.TimeStr+"_"+o.VariableNameSnake+".go")
	}

	return o
}

// discoverServices scans cmd/ directory to find existing services
func discoverServices() ([]string, error) {
	entries, err := os.ReadDir("cmd")
	if err != nil {
		return nil, err
	}

	var services []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Extract service name: myapp-apiserver → apiserver, myappctl → ctl
		name := entry.Name()
		parts := strings.Split(name, "-")
		if len(parts) > 1 {
			services = append(services, parts[len(parts)-1])
		} else if strings.HasSuffix(name, "ctl") {
			services = append(services, "ctl")
		}
	}

	return services, nil
}
