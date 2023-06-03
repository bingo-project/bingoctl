package util

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/goer-project/goer-utils/console"
)

// GenerateGoCode generate go source file.
func GenerateGoCode(filePath, codeTemplate, name string, o any) error {
	if Exists(filePath) {
		console.Warn(filePath + " already exists!")

		return nil
	}

	// Log
	fmt.Printf(" - Generating: ")
	console.Info(filePath)

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
