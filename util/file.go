package util

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

// GenerateGoCode generate go source file.
func GenerateGoCode(filePath, codeTemplate string, o any) error {
	directory := GetDirectoryFromPath(filePath)
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		return err
	}

	tmpl, err := template.New("cmd").Parse(codeTemplate)
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

	fmt.Printf("Command file generated: %s\n", filePath)

	return nil
}

func GetDirectoryFromPath(filePath string) string {
	arr := strings.Split(filePath, "/")
	directorArr := arr[:len(arr)-1]

	return strings.Join(directorArr, "/")
}
