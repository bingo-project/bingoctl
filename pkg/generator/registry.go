package generator

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/mgutz/ansi"

	"github.com/bingo-project/bingoctl/pkg/config"
)

func (o *Options) Register(registry config.Registry, interfaceTemplate, registerTemplate string) error {
	if registry.Filepath == "" {
		return nil
	}

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
	newContent, err := RegisterInterface(registry.Interface, string(content), interfaceTemplate, registerTemplate)
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

func RegisterInterface(name, content, interfaceTemplate, registerTemplate string) (data string, err error) {
	// Check if already registered
	if strings.Contains(content, interfaceTemplate) {
		return "", errors.New("interface already registered: " + interfaceTemplate)
	}

	// Register
	seek := fmt.Sprintf("type %s interface {", name)
	rule := seek + `([^}]*)}`
	reg := regexp.MustCompile(rule)

	// 根据规则提取关键信息
	results := reg.FindAllString(content, -1)
	if len(results) == 0 {
		err = errors.New("not matched")

		return
	}

	old := results[0]
	str := strings.TrimRight(old, "}")
	str = str + "\t" + interfaceTemplate + "\n}"

	newContent := strings.Replace(content, old, str, 1)
	newContent = newContent + "\n" + registerTemplate

	return newContent, nil
}
