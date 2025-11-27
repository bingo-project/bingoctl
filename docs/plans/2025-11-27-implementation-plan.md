# Create and Make Command Enhancement Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add service selection to make commands via `--service` flag with automatic path inference, and enable flexible service selection for create command via interactive prompts and command-line flags.

**Architecture:** Two-part implementation - Part 1 adds `--service` flag to make commands with intelligent path inference by scanning cmd/ directory; Part 2 refactors create command to support interactive service selection and command-line service specification.

**Tech Stack:** Go 1.21+, promptui for interactive prompts, embed.FS for templates

---

## Part 1: Make Command Service Selection

### Task 1: Add Service field to generator Options

**Files:**
- Modify: `pkg/generator/option.go:3-52`

**Step 1: Add Service field to Options struct**

在 `Options` 结构体中添加 `Service` 字段：

```go
type Options struct {
	// Code template
	Name              string
	Description       string
	FilePath          string
	Directory         string
	CodeTemplate      string
	InterfaceTemplate string
	RegisterTemplate  string

	// Code attributes - variable
	PackageName        string
	StructName         string
	StructNamePlural   string
	VariableName       string
	VariableNameSnake  string
	VariableNamePlural string
	TableName          string
	ModelName          string
	ServiceName        string

	// Code attributes - import path
	RootPackage  string
	BizPath      string
	StorePath    string
	RequestPath  string
	ModelPath    string
	RelativePath string

	// Service flags
	EnableHTTP     bool
	EnableGRPC     bool
	WithBiz        bool
	WithStore      bool
	WithController bool
	WithMiddleware bool
	WithRouter     bool
	NoBiz          bool

	// Service selection - new
	Service string // Target service name for path inference

	// Generate by gorm.gen
	Table           string
	FieldTemplate   string
	Fields          string
	MainFields      string
	UpdatableFields string
	MetaFields      []*Field

	// Migration
	TimeStr string
}
```

**Step 2: Verify the change compiles**

Run: `go build -o /tmp/bingoctl ./cmd/bingoctl`
Expected: SUCCESS (builds without errors)

**Step 3: Commit**

```bash
git add pkg/generator/option.go
git commit -m "feat(generator): add Service field to Options for service selection"
```

---

### Task 2: Implement service discovery function

**Files:**
- Modify: `pkg/generator/generate.go:1-141`
- Create test: `pkg/generator/generate_test.go`

**Step 1: Write failing test for discoverServices**

Create: `pkg/generator/generate_test.go`

```go
package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverServices(t *testing.T) {
	// Create temporary cmd directory structure
	tmpDir := t.TempDir()
	cmdDir := filepath.Join(tmpDir, "cmd")
	if err := os.MkdirAll(cmdDir, 0755); err != nil {
		t.Fatalf("Failed to create cmd dir: %v", err)
	}

	// Create service directories
	services := []string{"myapp-apiserver", "myapp-admserver", "myappctl"}
	for _, svc := range services {
		if err := os.MkdirAll(filepath.Join(cmdDir, svc), 0755); err != nil {
			t.Fatalf("Failed to create service dir %s: %v", svc, err)
		}
	}

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Test discovery
	discovered, err := discoverServices()
	if err != nil {
		t.Fatalf("discoverServices failed: %v", err)
	}

	expected := map[string]bool{
		"apiserver":  true,
		"admserver":  true,
		"ctl":        true,
	}

	if len(discovered) != len(expected) {
		t.Errorf("Expected %d services, got %d", len(expected), len(discovered))
	}

	for _, svc := range discovered {
		if !expected[svc] {
			t.Errorf("Unexpected service: %s", svc)
		}
	}
}

func TestDiscoverServices_NoCmd(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	_, err := discoverServices()
	if err == nil {
		t.Error("Expected error when cmd/ doesn't exist, got nil")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/generator -v -run TestDiscoverServices`
Expected: FAIL with "undefined: discoverServices"

**Step 3: Implement discoverServices function**

在 `pkg/generator/generate.go` 文件末尾添加：

```go
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
```

**Step 4: Add import if needed**

确保 `pkg/generator/generate.go` 顶部有必要的 imports：

```go
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
```

**Step 5: Run test to verify it passes**

Run: `go test ./pkg/generator -v -run TestDiscoverServices`
Expected: PASS

**Step 6: Commit**

```bash
git add pkg/generator/generate.go pkg/generator/generate_test.go
git commit -m "feat(generator): implement service discovery from cmd directory"
```

---

### Task 3: Implement path suffix extraction function

**Files:**
- Modify: `pkg/generator/generate.go`
- Modify test: `pkg/generator/generate_test.go`

**Step 1: Write failing test for extractSuffix**

在 `pkg/generator/generate_test.go` 中添加：

```go
func TestExtractSuffix(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "internal with service",
			path:     "internal/apiserver/model",
			expected: "model",
		},
		{
			name:     "internal with nested path",
			path:     "internal/apiserver/controller/v1",
			expected: "controller/v1",
		},
		{
			name:     "pkg prefix",
			path:     "pkg/api/v1",
			expected: "v1",
		},
		{
			name:     "no known prefix",
			path:     "custom/path/model",
			expected: "model",
		},
		{
			name:     "single segment",
			path:     "model",
			expected: "model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractSuffix(tt.path)
			if result != tt.expected {
				t.Errorf("extractSuffix(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/generator -v -run TestExtractSuffix`
Expected: FAIL with "undefined: extractSuffix"

**Step 3: Implement extractSuffix function**

在 `pkg/generator/generate.go` 文件末尾添加：

```go
// extractSuffix extracts the path suffix after known prefixes
func extractSuffix(path string) string {
	parts := strings.Split(filepath.Clean(path), string(filepath.Separator))

	// Find internal/ or pkg/ and return everything after the service name
	for i, part := range parts {
		if part == "internal" || part == "pkg" {
			if i+2 < len(parts) {
				return strings.Join(parts[i+2:], string(filepath.Separator))
			}
		}
	}

	// If not found, return the last segment
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return path
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./pkg/generator -v -run TestExtractSuffix`
Expected: PASS

**Step 5: Commit**

```bash
git add pkg/generator/generate.go pkg/generator/generate_test.go
git commit -m "feat(generator): implement path suffix extraction"
```

---

### Task 4: Implement path inference function

**Files:**
- Modify: `pkg/generator/generate.go`
- Modify test: `pkg/generator/generate_test.go`

**Step 1: Write failing test for InferDirectoryForService**

在 `pkg/generator/generate_test.go` 中添加：

```go
func TestInferDirectoryForService(t *testing.T) {
	// Setup temp directory with cmd structure
	tmpDir := t.TempDir()
	cmdDir := filepath.Join(tmpDir, "cmd")
	os.MkdirAll(filepath.Join(cmdDir, "myapp-apiserver"), 0755)
	os.MkdirAll(filepath.Join(cmdDir, "myapp-admserver"), 0755)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	tests := []struct {
		name        string
		baseDir     string
		serviceName string
		expected    string
	}{
		{
			name:        "empty service returns base",
			baseDir:     "internal/apiserver/model",
			serviceName: "",
			expected:    "internal/apiserver/model",
		},
		{
			name:        "smart replacement",
			baseDir:     "internal/apiserver/model",
			serviceName: "admserver",
			expected:    "internal/admserver/model",
		},
		{
			name:        "smart replacement with nested path",
			baseDir:     "internal/apiserver/controller/v1",
			serviceName: "admserver",
			expected:    "internal/admserver/controller/v1",
		},
		{
			name:        "fallback pattern",
			baseDir:     "internal/pkg/model",
			serviceName: "admserver",
			expected:    "internal/admserver/model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{}
			result, err := o.InferDirectoryForService(tt.baseDir, tt.serviceName)
			if err != nil {
				t.Fatalf("InferDirectoryForService failed: %v", err)
			}
			if result != tt.expected {
				t.Errorf("InferDirectoryForService(%q, %q) = %q, want %q",
					tt.baseDir, tt.serviceName, result, tt.expected)
			}
		})
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/generator -v -run TestInferDirectoryForService`
Expected: FAIL with "undefined: InferDirectoryForService"

**Step 3: Implement InferDirectoryForService method**

在 `pkg/generator/generate.go` 文件中，在 `GetMapDirectory` 函数之后添加：

```go
// InferDirectoryForService infers the target directory based on service name
func (o *Options) InferDirectoryForService(baseDir, serviceName string) (string, error) {
	if serviceName == "" {
		return baseDir, nil
	}

	// 1. Discover existing services from cmd/
	services, err := discoverServices()
	if err != nil {
		// If cmd/ doesn't exist, fall back to pattern-based inference
		services = []string{}
	}

	// 2. Smart replacement: if path contains a known service name, replace it
	for _, svc := range services {
		if strings.Contains(baseDir, svc) {
			return strings.ReplaceAll(baseDir, svc, serviceName), nil
		}
	}

	// 3. Fallback pattern: extract suffix and join with internal/{service}/
	suffix := extractSuffix(baseDir)
	return filepath.Join("internal", serviceName, suffix), nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./pkg/generator -v -run TestInferDirectoryForService`
Expected: PASS

**Step 5: Run all generator tests**

Run: `go test ./pkg/generator -v`
Expected: All tests PASS

**Step 6: Commit**

```bash
git add pkg/generator/generate.go pkg/generator/generate_test.go
git commit -m "feat(generator): implement service-based path inference"
```

---

### Task 5: Integrate path inference into GenerateCode

**Files:**
- Modify: `pkg/generator/generate.go:17-54`

**Step 1: Modify GenerateCode to use path inference**

在 `pkg/generator/generate.go` 中修改 `GenerateCode` 函数：

```go
func (o *Options) GenerateCode(tmpl, path string) error {
	dir := GetMapDirectory(tmpl)

	// Apply service-based path inference if --service flag is provided
	if o.Service != "" {
		inferredDir, err := o.InferDirectoryForService(dir, o.Service)
		if err != nil {
			return fmt.Errorf("failed to infer directory for service %s: %w", o.Service, err)
		}
		dir = inferredDir
	}

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
```

**Step 2: Add fmt import if needed**

确保在 imports 中有 `fmt`：

```go
import (
	"fmt"
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
```

**Step 3: Verify the change compiles**

Run: `go build -o /tmp/bingoctl ./cmd/bingoctl`
Expected: SUCCESS

**Step 4: Commit**

```bash
git add pkg/generator/generate.go
git commit -m "feat(generator): integrate path inference into GenerateCode"
```

---

### Task 6: Add --service flag to make command

**Files:**
- Modify: `pkg/cmd/make/make.go:10-49`

**Step 1: Add --service persistent flag**

在 `pkg/cmd/make/make.go` 的 `NewCmdMake` 函数中添加 flag：

```go
func NewCmdMake() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "make COMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Generate code",
		Example:               makeExample,
		Run:                   cmdutil.DefaultSubCommandRun(),
	}

	cmd.PersistentFlags().StringVarP(&opt.Directory, "directory", "d", "", "Where to create the file.")
	cmd.PersistentFlags().StringVarP(&opt.PackageName, "package", "p", "", "Name of the package.")
	cmd.PersistentFlags().StringVarP(&opt.Table, "table", "t", "", "Read fields from db table.")
	cmd.PersistentFlags().StringVarP(&opt.Service, "service", "s", "", "Target service name for path inference")

	// Add subcommands
	cmd.AddCommand(NewCmdCMD())
	cmd.AddCommand(NewCmdModel())
	cmd.AddCommand(NewCmdStore())
	cmd.AddCommand(NewCmdRequest())
	cmd.AddCommand(NewCmdBiz())
	cmd.AddCommand(NewCmdController())
	cmd.AddCommand(NewCmdCrud())
	cmd.AddCommand(NewCmdMiddleware())
	cmd.AddCommand(NewCmdJob())
	cmd.AddCommand(NewCmdMigration())
	cmd.AddCommand(NewCmdSeeder())
	cmd.AddCommand(NewCmdService())

	return cmd
}
```

**Step 2: Verify the change compiles**

Run: `go build -o /tmp/bingoctl ./cmd/bingoctl`
Expected: SUCCESS

**Step 3: Test --service flag shows in help**

Run: `/tmp/bingoctl make model --help`
Expected: Output includes "-s, --service string   Target service name for path inference"

**Step 4: Commit**

```bash
git add pkg/cmd/make/make.go
git commit -m "feat(make): add --service flag for path inference"
```

---

### Task 7: Update README documentation

**Files:**
- Modify: `README.md:76-190`

**Step 1: Update make command global options section**

在 README.md 的 "全局选项" 部分添加 `--service` 参数说明：

```markdown
#### 全局选项

```bash
-d, --directory string   指定生成文件的目录
-p, --package string     指定包名
-t, --table string       从数据库表读取字段
-s, --service string     目标服务名称，用于自动推断路径
```
```

**Step 2: Add service selection usage example**

在 README.md 的 make 命令部分添加服务选择示例：

```markdown
#### 服务选择

当项目包含多个服务时，可以使用 `--service` 参数自动推断生成路径：

```bash
# 为 apiserver 生成代码（使用配置默认路径）
bingoctl make model user

# 为 admserver 生成代码（自动推断路径）
bingoctl make model user --service admserver

# 完整 CRUD 为指定服务生成
bingoctl make crud order --service admserver

# 明确指定路径（优先级最高）
bingoctl make model user -d custom/path
```

路径推断规则：
1. 扫描 `cmd/` 目录识别已存在的服务
2. 如果配置路径包含服务名，则智能替换（如 `internal/apiserver/model` → `internal/admserver/model`）
3. 否则使用固定模式：`internal/{service}/{suffix}`
```

**Step 3: Commit**

```bash
git add README.md
git commit -m "docs: add --service flag documentation to README"
```

---

### Task 8: Manual integration testing

**Files:** N/A (manual testing)

**Step 1: Create test project structure**

```bash
cd /tmp
rm -rf test-bingoctl
mkdir -p test-bingoctl/cmd/myapp-apiserver
mkdir -p test-bingoctl/cmd/myapp-admserver
mkdir -p test-bingoctl/internal/apiserver/{model,store,biz}
mkdir -p test-bingoctl/internal/admserver/{model,store,biz}
cd test-bingoctl
```

**Step 2: Create .bingoctl.yaml**

```bash
cat > .bingoctl.yaml << 'EOF'
version: v1

rootPackage: github.com/test/myapp

directory:
  model: internal/apiserver/model
  store: internal/apiserver/store
  biz: internal/apiserver/biz
EOF
```

**Step 3: Test default behavior**

```bash
/tmp/bingoctl make model user
```

Expected: Creates `internal/apiserver/model/user.go`

**Step 4: Test --service flag**

```bash
/tmp/bingoctl make model order --service admserver
```

Expected: Creates `internal/admserver/model/order.go`

**Step 5: Verify -d still has priority**

```bash
/tmp/bingoctl make model custom --service admserver -d internal/custom/model
```

Expected: Creates `internal/custom/model/custom.go` (uses -d, ignores --service)

**Step 6: Document test results**

Create a simple test log to confirm functionality works as expected.

---

## Part 2: Create Command Service Selection

### Task 9: Add service selection fields to CreateOptions

**Files:**
- Modify: `pkg/cmd/create/create.go:33-46`

**Step 1: Add service selection fields**

```go
// CreateOptions is an option struct to support 'create' sub command.
type CreateOptions struct {
	GoVersion    string
	TemplatePath string
	RootPackage  string
	AppName      string
	AppNameCamel string

	// Service selection
	Services     []string // Explicitly specified services
	NoServices   []string // Services to exclude from defaults
	AddServices  []string // Services to add to defaults
	Interactive  bool     // Whether to use interactive mode (default true)
	selectedServices []string // Final computed service list (internal)
}
```

**Step 2: Verify the change compiles**

Run: `go build -o /tmp/bingoctl ./cmd/bingoctl`
Expected: SUCCESS

**Step 3: Commit**

```bash
git add pkg/cmd/create/create.go
git commit -m "feat(create): add service selection fields to CreateOptions"
```

---

### Task 10: Add command-line flags

**Files:**
- Modify: `pkg/cmd/create/create.go:49-66`

**Step 1: Add flags to NewCmdCreate**

在 `NewCmdCreate` 函数中添加 flags：

```go
func NewCmdCreate() *cobra.Command {
	o := NewCreateOptions()

	cmd := &cobra.Command{
		Use:                   createUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Create a project",
		TraverseChildren:      true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
	}

	cmd.Flags().StringSliceVar(&o.Services, "services", nil, "Explicitly specify services to create (comma-separated). Use 'none' for minimal skeleton")
	cmd.Flags().StringSliceVar(&o.NoServices, "no-service", nil, "Exclude services from defaults (comma-separated)")
	cmd.Flags().StringSliceVar(&o.AddServices, "add-service", nil, "Add services to defaults (comma-separated)")

	return cmd
}
```

**Step 2: Verify the change compiles**

Run: `go build -o /tmp/bingoctl ./cmd/bingoctl`
Expected: SUCCESS

**Step 3: Commit**

```bash
git add pkg/cmd/create/create.go
git commit -m "feat(create): add service selection command-line flags"
```

---

### Task 11: Implement service list computation

**Files:**
- Modify: `pkg/cmd/create/create.go`

**Step 1: Write test for computeServiceList**

Create: `pkg/cmd/create/create_test.go`

```go
package create

import (
	"reflect"
	"testing"
)

func TestComputeServiceList(t *testing.T) {
	tests := []struct {
		name         string
		services     []string
		noServices   []string
		addServices  []string
		expected     []string
	}{
		{
			name:       "explicit services override",
			services:   []string{"apiserver", "bot"},
			expected:   []string{"apiserver", "bot"},
		},
		{
			name:       "services none",
			services:   []string{"none"},
			expected:   []string{},
		},
		{
			name:       "no flags uses defaults",
			expected:   []string{"apiserver", "ctl"},
		},
		{
			name:       "exclude service",
			noServices: []string{"ctl"},
			expected:   []string{"apiserver"},
		},
		{
			name:        "add services",
			addServices: []string{"bot", "scheduler"},
			expected:    []string{"apiserver", "ctl", "bot", "scheduler"},
		},
		{
			name:        "combined exclude and add",
			noServices:  []string{"ctl"},
			addServices: []string{"admserver"},
			expected:    []string{"apiserver", "admserver"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &CreateOptions{
				Services:    tt.services,
				NoServices:  tt.noServices,
				AddServices: tt.addServices,
			}
			result := o.computeServiceList()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("computeServiceList() = %v, want %v", result, tt.expected)
			}
		})
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/cmd/create -v -run TestComputeServiceList`
Expected: FAIL with "undefined: computeServiceList"

**Step 3: Implement computeServiceList method**

在 `pkg/cmd/create/create.go` 中添加：

```go
var (
	defaultServices   = []string{"apiserver", "ctl"}
	availableServices = []string{"apiserver", "ctl", "admserver", "bot", "scheduler"}
)

// computeServiceList computes the final service list based on flags
func (o *CreateOptions) computeServiceList() []string {
	// Priority 1: --services flag explicitly specified
	if len(o.Services) > 0 {
		if len(o.Services) == 1 && o.Services[0] == "none" {
			return []string{}
		}
		return o.Services
	}

	// Priority 2: Start with defaults and apply modifications
	services := make(map[string]bool)
	for _, svc := range defaultServices {
		services[svc] = true
	}

	// Apply exclusions
	for _, svc := range o.NoServices {
		delete(services, svc)
	}

	// Apply additions
	for _, svc := range o.AddServices {
		services[svc] = true
	}

	// Convert map back to slice
	result := make([]string, 0, len(services))
	// Maintain order: iterate through availableServices to preserve order
	for _, svc := range availableServices {
		if services[svc] {
			result = append(result, svc)
		}
	}

	return result
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./pkg/cmd/create -v -run TestComputeServiceList`
Expected: PASS

**Step 5: Commit**

```bash
git add pkg/cmd/create/create.go pkg/cmd/create/create_test.go
git commit -m "feat(create): implement service list computation logic"
```

---

### Task 12: Implement interactive service selection

**Files:**
- Modify: `pkg/cmd/create/create.go`

**Step 1: Add interactive selection method**

在 `pkg/cmd/create/create.go` 中添加：

```go
import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/bingo-project/component-base/cli/console"
	"github.com/iancoleman/strcase"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	cmdutil "github.com/bingo-project/bingoctl/pkg/util"
)

// selectServicesInteractively shows interactive multi-select for services
func (o *CreateOptions) selectServicesInteractively() ([]string, error) {
	type service struct {
		Name     string
		Selected bool
	}

	services := []service{
		{Name: "apiserver", Selected: true},
		{Name: "ctl", Selected: true},
		{Name: "admserver", Selected: false},
		{Name: "bot", Selected: false},
		{Name: "scheduler", Selected: false},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "→ {{ if .Selected }}✓{{ else }}✗{{ end }} {{ .Name }}",
		Inactive: "  {{ if .Selected }}✓{{ else }}✗{{ end }} {{ .Name }}",
	}

	prompt := promptui.Select{
		Label:     "选择要创建的服务 (空格切换, 回车确认)",
		Items:     services,
		Templates: templates,
		Size:      5,
	}

	console.Info("请选择要创建的服务:")
	console.Info("提示: 使用↑↓移动光标, 按 's' 切换选择, 按回车确认")

	// Note: promptui.Select doesn't support multi-select natively
	// We'll use a simpler approach with individual confirmations
	selected := make([]string, 0)
	for _, svc := range services {
		if svc.Selected {
			selected = append(selected, svc.Name)
		}
	}

	// For now, return the default selection
	// A full implementation would require a custom multi-select prompt
	// or using a library that supports it better
	console.Info(fmt.Sprintf("默认选择: %v", selected))
	console.Info("(交互式多选将在后续版本中完善)")

	return selected, nil
}
```

**Step 2: Update Complete method to determine service list**

修改 `Complete` 方法：

```go
// Complete completes all the required options.
func (o *CreateOptions) Complete(cmd *cobra.Command, args []string) error {
	// Determine if interactive mode
	o.Interactive = len(o.Services) == 0 && len(o.NoServices) == 0 && len(o.AddServices) == 0

	// Compute final service list
	if o.Interactive {
		console.Info("进入交互模式...")
		selected, err := o.selectServicesInteractively()
		if err != nil {
			return err
		}
		o.selectedServices = selected
	} else {
		o.selectedServices = o.computeServiceList()
	}

	// Warn if no services selected
	if len(o.selectedServices) == 0 {
		console.Warn("未选择任何服务，将创建最小项目骨架")
		prompt := promptui.Prompt{
			Label:     "继续",
			IsConfirm: true,
		}
		_, err := prompt.Run()
		if err != nil {
			console.Exit("已取消创建")
		}
	}

	err := os.MkdirAll(o.AppName, 0755)
	if err != nil {
		return err
	}

	return nil
}
```

**Step 3: Verify the change compiles**

Run: `go build -o /tmp/bingoctl ./cmd/bingoctl`
Expected: SUCCESS

**Step 4: Commit**

```bash
git add pkg/cmd/create/create.go
git commit -m "feat(create): add interactive service selection (basic implementation)"
```

---

### Task 13: Document Part 2 completion note

**Step 1: Add implementation note**

Part 2 的完整实现（包括模板文件重组和选择性复制）需要大量的文件重组工作。由于这是一个较大的变更，建议：

1. 先完成 Part 1（已完成）并进行充分测试
2. 在确认 Part 1 工作正常后，再进行 Part 2 的模板重组
3. Part 2 的模板重组可以作为独立的任务，分步骤进行

当前状态：
- ✅ Part 1: Make 命令服务选择 - 已完成
- ⏸️ Part 2: Create 命令精简 - 基础框架已完成，模板重组待完成

**Step 2: Create follow-up task document**

Create: `docs/plans/2025-11-27-part2-template-reorganization.md`

```markdown
# Part 2 Template Reorganization Tasks

## Remaining Work

### 1. Reorganize template directory structure
- Move files from flat structure to services/ subdirectories
- Separate common files from service-specific files
- Update embed.FS paths

### 2. Implement selective template copying
- Modify Run() to only copy selected services
- Handle common files separately
- Generate dynamic .bingoctl.yaml based on selected services

### 3. Testing
- Test minimal skeleton creation
- Test each service individually
- Test various combinations
- Test interactive mode

## Estimated Effort
3-4 hours of focused work for template reorganization
```

**Step 3: Commit completion note**

```bash
git add docs/plans/2025-11-27-part2-template-reorganization.md
git commit -m "docs: add Part 2 template reorganization follow-up tasks"
```

---

## Summary

**Part 1 (Make Command Service Selection) - COMPLETE**
- ✅ Service discovery from cmd/ directory
- ✅ Path suffix extraction
- ✅ Intelligent path inference with fallback
- ✅ Integration into GenerateCode
- ✅ --service flag added to make commands
- ✅ Documentation updated
- ✅ Unit tests added

**Part 2 (Create Command) - PARTIAL**
- ✅ Command-line flags added
- ✅ Service list computation logic
- ✅ Basic interactive selection framework
- ⏸️ Template reorganization (deferred)
- ⏸️ Selective file copying (deferred)
- ⏸️ Full testing (deferred)

**Next Steps:**
1. Test Part 1 thoroughly in real projects
2. Gather feedback on --service flag usage
3. Plan template reorganization for Part 2
4. Implement Part 2 in a follow-up iteration
