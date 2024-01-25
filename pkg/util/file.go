package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/manifoldco/promptui"
	"github.com/mgutz/ansi"
)

var Overwrite bool

// GenerateCode generate go source file.
func GenerateCode(filePath, codeTemplate, name string, o any) error {
	if Exists(filePath) && !Overwrite {
		prompt := promptui.Prompt{
			Label:     "Overwrite " + ansi.Color(filePath, "yellow"),
			IsConfirm: true,
		}

		_, err := prompt.Run()
		if err != nil {
			return err
		}
	}

	directory := GetDirectoryFromPath(filePath)
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		return err
	}

	tmpl := template.New(name)
	if name == "init" {
		tmpl.Delims("{[", "]}")
	}
	tmpl, err = tmpl.Parse(codeTemplate)
	if err != nil {
		return err
	}

	fd, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer fd.Close()

	err = tmpl.Execute(fd, o)
	if err != nil {
		return err
	}

	fmt.Printf("%s %s\n", ansi.Color("Generated:", "green"), filePath)

	return nil
}

func Exists(fileToCheck string) bool {
	if _, err := os.Stat(fileToCheck); os.IsNotExist(err) {
		return false
	}

	return true
}

func GetDirectoryFromPath(filePath string) string {
	arr := strings.Split(filePath, "/")
	directorArr := arr[:len(arr)-1]

	return strings.Join(directorArr, "/")
}

func GetFileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
