package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/bingo-project/bingoctl/pkg/config"
)

// getAppName extracts the application name from the root package.
// e.g., "github.com/xxx/demo" -> "demo"
func getAppName() string {
	parts := strings.Split(config.Cfg.RootPackage, "/")
	return parts[len(parts)-1]
}

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
		if err := o.generateHTTP(); err != nil {
			return err
		}
	}

	// Generate gRPC server if enabled
	if o.EnableGRPC {
		if err := o.generateGRPCServer(); err != nil {
			return err
		}
	}

	// Generate WebSocket server if enabled
	if o.EnableWS {
		if err := o.generateWS(); err != nil {
			return err
		}
	}

	// Determine if we should generate router and handler
	// When --http, --grpc or --ws is set, generate by default unless --no-router or --no-handler is specified
	hasServer := o.EnableHTTP || o.EnableGRPC || o.EnableWS
	shouldGenerateRouter := hasServer && !o.NoRouter
	shouldGenerateHandler := hasServer && !o.NoHandler

	// Generate optional directories
	// Generate biz directory by default unless --no-biz is specified
	if !o.NoBiz {
		if err := o.createDirectoryWithFile("biz", "internal", o.ServiceName, "biz"); err != nil {
			return err
		}
	}
	if o.WithStore {
		if err := o.createDirectory("internal", o.ServiceName, "store"); err != nil {
			return err
		}
	}
	if shouldGenerateHandler {
		if err := o.generateHandler(); err != nil {
			return err
		}
	}
	if o.WithMiddleware {
		if err := o.createDirectory("internal", o.ServiceName, "middleware"); err != nil {
			return err
		}
	}
	if shouldGenerateRouter {
		if err := o.generateRouter(); err != nil {
			return err
		}
	}

	// Generate config file
	if err := o.generateConfig(); err != nil {
		return err
	}

	appName := getAppName()
	cmdDirName := appName + "-" + o.ServiceName
	fmt.Printf("Service '%s' generated successfully!\n", o.ServiceName)
	fmt.Printf("  - cmd/%s/main.go\n", cmdDirName)
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

	// cmd directory uses app-service format, e.g., demo-admin
	appName := getAppName()
	cmdDirName := appName + "-" + o.ServiceName
	cmdDir := filepath.Join("cmd", cmdDirName)
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
	tplContent, err := ReadServiceTemplate("run.go.tpl")
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

	data := map[string]any{
		"RootPackage": config.Cfg.RootPackage,
		"ServiceName": o.ServiceName,
		"EnableHTTP":  o.EnableHTTP,
		"EnableGRPC":  o.EnableGRPC,
		"EnableWS":    o.EnableWS,
	}

	return tmpl.Execute(file, data)
}

func (o *Options) generateHTTP() error {
	tplContent, err := ReadServiceTemplate("http.go.tpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("http").Parse(string(tplContent))
	if err != nil {
		return err
	}

	internalDir := filepath.Join("internal", o.ServiceName)
	filePath := filepath.Join(internalDir, "http.go")
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

func (o *Options) generateWS() error {
	tplContent, err := ReadServiceTemplate("ws.go.tpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("ws").Parse(string(tplContent))
	if err != nil {
		return err
	}

	internalDir := filepath.Join("internal", o.ServiceName)
	filePath := filepath.Join(internalDir, "ws.go")
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

func (o *Options) generateHandler() error {
	data := map[string]string{
		"RootPackage": config.Cfg.RootPackage,
		"ServiceName": o.ServiceName,
	}

	// Generate HTTP handler if HTTP is enabled
	if o.EnableHTTP {
		handlerDir := filepath.Join("internal", o.ServiceName, "handler", "http")
		if err := os.MkdirAll(handlerDir, 0755); err != nil {
			return err
		}

		tplContent, err := ReadServiceTemplate("handler_http.go.tpl")
		if err != nil {
			return err
		}

		tmpl, err := template.New("handler_http").Parse(string(tplContent))
		if err != nil {
			return err
		}

		filePath := filepath.Join(handlerDir, "handler.go")
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		if err := tmpl.Execute(file, data); err != nil {
			return err
		}
	}

	// Generate gRPC handler if gRPC is enabled
	if o.EnableGRPC {
		handlerDir := filepath.Join("internal", o.ServiceName, "handler", "grpc")
		if err := os.MkdirAll(handlerDir, 0755); err != nil {
			return err
		}

		tplContent, err := ReadServiceTemplate("handler_grpc.go.tpl")
		if err != nil {
			return err
		}

		tmpl, err := template.New("handler_grpc").Parse(string(tplContent))
		if err != nil {
			return err
		}

		filePath := filepath.Join(handlerDir, "handler.go")
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		if err := tmpl.Execute(file, data); err != nil {
			return err
		}
	}

	// Generate WebSocket handler if WebSocket is enabled
	if o.EnableWS {
		handlerDir := filepath.Join("internal", o.ServiceName, "handler", "ws")
		if err := os.MkdirAll(handlerDir, 0755); err != nil {
			return err
		}

		tplContent, err := ReadServiceTemplate("handler_ws.go.tpl")
		if err != nil {
			return err
		}

		tmpl, err := template.New("handler_ws").Parse(string(tplContent))
		if err != nil {
			return err
		}

		filePath := filepath.Join(handlerDir, "handler.go")
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		if err := tmpl.Execute(file, data); err != nil {
			return err
		}
	}

	return nil
}

func (o *Options) generateRouter() error {
	routerDir := filepath.Join("internal", o.ServiceName, "router")
	if err := os.MkdirAll(routerDir, 0755); err != nil {
		return err
	}

	data := map[string]string{
		"RootPackage": config.Cfg.RootPackage,
		"ServiceName": o.ServiceName,
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

		if err := tmpl.Execute(file, data); err != nil {
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

		if err := tmpl.Execute(file, data); err != nil {
			return err
		}
	}

	// Generate WebSocket router if WebSocket is enabled
	if o.EnableWS {
		tplContent, err := ReadServiceTemplate("router_ws.go.tpl")
		if err != nil {
			return err
		}

		tmpl, err := template.New("router_ws").Parse(string(tplContent))
		if err != nil {
			return err
		}

		filePath := filepath.Join(routerDir, "ws.go")
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		if err := tmpl.Execute(file, data); err != nil {
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
		"EnableWS":   o.EnableWS,
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
	case "handler":
		fileName = "handler.go"
		tplName = "handler.go.tpl"
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
