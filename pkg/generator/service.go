package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/bingo-project/bingoctl/pkg/config"
)

// GenerateService generates a new service module with configurable HTTP/gRPC servers and business layers.
func (o *Options) GenerateService(name string) error {
	o.ServiceName = name

	// Generate cmd/<name>/main.go
	if err := o.generateCmdMain(); err != nil {
		return err
	}

	// Generate internal/<name>/app.go
	if err := o.generateApp(); err != nil {
		return err
	}

	// Generate internal/<name>/run.go
	if err := o.generateRun(); err != nil {
		return err
	}

	// Generate HTTP server if enabled
	if o.EnableHTTP {
		if err := o.generateHTTPServer(); err != nil {
			return err
		}
	}

	// Generate gRPC server if enabled
	if o.EnableGRPC {
		if err := o.generateGRPCServer(); err != nil {
			return err
		}
		// Create grpc/ directory
		if err := o.createDirectory("internal", o.ServiceName, "grpc"); err != nil {
			return err
		}
	}

	// Generate optional directories
	// Generate biz directory by default unless --no-biz is specified
	if o.WithBiz && !o.NoBiz {
		if err := o.createDirectoryWithFile("biz", "internal", o.ServiceName, "biz"); err != nil {
			return err
		}
	}
	if o.WithStore {
		if err := o.createDirectory("internal", o.ServiceName, "store"); err != nil {
			return err
		}
	}
	if o.WithController {
		if err := o.createDirectoryWithFile("controller", "internal", o.ServiceName, "controller"); err != nil {
			return err
		}
	}
	if o.WithMiddleware {
		if err := o.createDirectory("internal", o.ServiceName, "middleware"); err != nil {
			return err
		}
	}
	if o.WithRouter {
		if err := o.generateRouter(); err != nil {
			return err
		}
	}

	// Generate config file
	if err := o.generateConfig(); err != nil {
		return err
	}

	fmt.Printf("Service '%s' generated successfully!\n", o.ServiceName)
	fmt.Printf("  - cmd/%s/main.go\n", o.ServiceName)
	fmt.Printf("  - internal/%s/\n", o.ServiceName)
	fmt.Printf("  - configs/%s.yaml\n", o.ServiceName)

	return nil
}

func (o *Options) generateCmdMain() error {
	tplContent, err := ReadServiceTemplate("cmd_main.go.tpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("cmd_main").Parse(string(tplContent))
	if err != nil {
		return err
	}

	cmdDir := filepath.Join("cmd", o.ServiceName)
	if err := os.MkdirAll(cmdDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(cmdDir, "main.go")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := map[string]string{
		"RootPackage": config.Cfg.RootPackage,
		"ServiceName": o.ServiceName,
	}

	return tmpl.Execute(file, data)
}

func (o *Options) generateApp() error {
	tplContent, err := ReadServiceTemplate("app.go.tpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("app").Parse(string(tplContent))
	if err != nil {
		return err
	}

	internalDir := filepath.Join("internal", o.ServiceName)
	if err := os.MkdirAll(internalDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(internalDir, "app.go")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := map[string]string{
		"RootPackage": config.Cfg.RootPackage,
		"ServiceName": o.ServiceName,
	}

	return tmpl.Execute(file, data)
}

func (o *Options) generateRun() error {
	// Select template based on server flags
	var tplName string
	if o.EnableHTTP && o.EnableGRPC {
		tplName = "run_both.go.tpl"
	} else if o.EnableHTTP {
		tplName = "run_http.go.tpl"
	} else if o.EnableGRPC {
		tplName = "run_grpc.go.tpl"
	} else {
		tplName = "run_minimal.go.tpl"
	}

	tplContent, err := ReadServiceTemplate(tplName)
	if err != nil {
		return err
	}

	tmpl, err := template.New("run").Parse(string(tplContent))
	if err != nil {
		return err
	}

	internalDir := filepath.Join("internal", o.ServiceName)
	filePath := filepath.Join(internalDir, "run.go")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := map[string]string{
		"ServiceName": o.ServiceName,
	}

	return tmpl.Execute(file, data)
}

func (o *Options) generateHTTPServer() error {
	tplContent, err := ReadServiceTemplate("server.go.tpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("server").Parse(string(tplContent))
	if err != nil {
		return err
	}

	internalDir := filepath.Join("internal", o.ServiceName)
	filePath := filepath.Join(internalDir, "server.go")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := map[string]string{
		"RootPackage": config.Cfg.RootPackage,
		"ServiceName": o.ServiceName,
	}

	return tmpl.Execute(file, data)
}

func (o *Options) generateGRPCServer() error {
	tplContent, err := ReadServiceTemplate("grpc.go.tpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("grpc").Parse(string(tplContent))
	if err != nil {
		return err
	}

	internalDir := filepath.Join("internal", o.ServiceName)
	filePath := filepath.Join(internalDir, "grpc.go")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := map[string]string{
		"RootPackage": config.Cfg.RootPackage,
		"ServiceName": o.ServiceName,
	}

	return tmpl.Execute(file, data)
}

func (o *Options) generateRouter() error {
	routerDir := filepath.Join("internal", o.ServiceName, "router")
	if err := os.MkdirAll(routerDir, 0755); err != nil {
		return err
	}

	// Generate HTTP router if HTTP is enabled
	if o.EnableHTTP {
		tplContent, err := ReadServiceTemplate("router_http.go.tpl")
		if err != nil {
			return err
		}

		tmpl, err := template.New("router_http").Parse(string(tplContent))
		if err != nil {
			return err
		}

		filePath := filepath.Join(routerDir, "http.go")
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		if err := tmpl.Execute(file, nil); err != nil {
			return err
		}
	}

	// Generate gRPC router if gRPC is enabled
	if o.EnableGRPC {
		tplContent, err := ReadServiceTemplate("router_grpc.go.tpl")
		if err != nil {
			return err
		}

		tmpl, err := template.New("router_grpc").Parse(string(tplContent))
		if err != nil {
			return err
		}

		filePath := filepath.Join(routerDir, "grpc.go")
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		if err := tmpl.Execute(file, nil); err != nil {
			return err
		}
	}

	return nil
}

func (o *Options) generateConfig() error {
	tplContent, err := ReadServiceTemplate("config.yaml.tpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("config").Parse(string(tplContent))
	if err != nil {
		return err
	}

	configsDir := "configs"
	if err := os.MkdirAll(configsDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(configsDir, o.ServiceName+".yaml")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := map[string]bool{
		"EnableHTTP": o.EnableHTTP,
		"EnableGRPC": o.EnableGRPC,
	}

	return tmpl.Execute(file, data)
}

func (o *Options) createDirectory(parts ...string) error {
	dir := filepath.Join(parts...)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create .gitkeep file
	gitkeepPath := filepath.Join(dir, ".gitkeep")
	file, err := os.Create(gitkeepPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

// createDirectoryWithFile creates a directory and a basic file for specific directories.
func (o *Options) createDirectoryWithFile(dirType string, parts ...string) error {
	dir := filepath.Join(parts...)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create a basic file based on directory type
	var fileName, tplName string
	switch dirType {
	case "biz":
		fileName = "biz.go"
		tplName = "biz.go.tpl"
	case "controller":
		fileName = "controller.go"
		tplName = "controller.go.tpl"
	default:
		// For other directories, just create .gitkeep
		gitkeepPath := filepath.Join(dir, ".gitkeep")
		file, err := os.Create(gitkeepPath)
		if err != nil {
			return err
		}
		defer file.Close()
		return nil
	}

	// Read template
	tplContent, err := ReadServiceTemplate(tplName)
	if err != nil {
		return err
	}

	tmpl, err := template.New(tplName).Parse(string(tplContent))
	if err != nil {
		return err
	}

	// Create file
	filePath := filepath.Join(dir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Execute template
	data := map[string]string{
		"RootPackage": config.Cfg.RootPackage,
		"ServiceName": o.ServiceName,
	}

	return tmpl.Execute(file, data)
}
