package generator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mgutz/ansi"

	"github.com/bingo-project/bingoctl/pkg/config"
)

func (o *Options) Register(registry config.Registry, interfaceTemplate, registerTemplate, importPath string) error {
	if registry.Filepath == "" {
		return nil
	}

	// Check if the registry file is in the same package as the generated file
	// If so, skip import to avoid self-referencing import
	registryDir := filepath.Dir(registry.Filepath)
	generatedDir := strings.TrimSuffix(o.Directory, "/")
	samePackage := registryDir == generatedDir

	// Package
	pkg := ""
	if o.PackageName != o.Name {
		pkg = o.PackageName + "."
	}

	// Replace
	replaces := make(map[string]string)
	replaces["{{.Package}}"] = pkg
	replaces["{{.StructName}}"] = o.StructName
	replaces["{{.StructNamePlural}}"] = o.StructNamePlural
	replaces["{{.VariableName}}"] = o.VariableName
	replaces["{{.VariableNameSnake}}"] = o.VariableNameSnake

	for search, replace := range replaces {
		interfaceTemplate = strings.ReplaceAll(interfaceTemplate, search, replace)
		registerTemplate = strings.ReplaceAll(registerTemplate, search, replace)
	}

	content, err := os.ReadFile(registry.Filepath)
	if err != nil {
		return err
	}

	// 注册 interface
	newContent, err := RegisterInterface(registry.Interface, string(content), interfaceTemplate, registerTemplate, importPath, samePackage)
	if err != nil {
		return err
	}

	err = os.WriteFile(registry.Filepath, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	fmt.Printf("%s %s\n", ansi.Color("Registered:", "green"), registry.Filepath)

	return nil
}

func RegisterInterface(name, content, interfaceTemplate, registerTemplate, importPath string, samePackage bool) (data string, err error) {
	// Check if already registered
	if strings.Contains(content, interfaceTemplate) {
		return "", errors.New("interface already registered: " + interfaceTemplate)
	}

	// Register
	pattern := fmt.Sprintf("type %s interface {([^}]*)}", name)
	match, err := Match(pattern, content)
	if err != nil {
		return "", err
	}

	str := strings.TrimRight(match, "}")
	str = str + "\t" + interfaceTemplate + "\n}"

	newContent := strings.Replace(content, match, str, 1)
	newContent = newContent + "\n" + registerTemplate

	// Skip import if same package to avoid self-referencing import
	if samePackage {
		return newContent, nil
	}

	// Import path
	if strings.Contains(newContent, importPath) {
		return newContent, nil
	}

	pattern = `import\s?\(([^}]*?)\)`
	match, err = Match(pattern, content)
	if err != nil {
		return "", err
	}

	str = strings.TrimRight(match, ")")
	str = fmt.Sprintf("%s\t\"%s\"\n)", str, importPath)

	newContent = strings.Replace(newContent, match, str, 1)

	return newContent, nil
}

func Match(pattern, content string) (string, error) {
	reg := regexp.MustCompile(pattern)

	// 根据规则提取关键信息
	results := reg.FindAllString(content, -1)
	if len(results) == 0 {
		return "", errors.New("not matched")
	}

	match := results[0]

	return match, nil
}
