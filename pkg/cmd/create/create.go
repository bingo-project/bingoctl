package create

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/bingo-project/component-base/cli/console"
	"github.com/iancoleman/strcase"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/template"
	cmdutil "github.com/bingo-project/bingoctl/pkg/util"
)

const (
	createUsageStr = "create NAME"
)

// ServiceMapping describes a service's directory structure
type ServiceMapping struct {
	Cmd      string
	Internal string
}

var (
	createUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the create command",
		createUsageStr,
	)

	defaultServices   = []string{"apiserver"}
	availableServices = []string{"apiserver", "ctl", "admserver", "bot", "scheduler"}

	// defaultServiceMapping defines the default service directory structure
	// Used when .bingo.yaml is not available
	defaultServiceMapping = map[string]ServiceMapping{
		"apiserver": {
			Cmd:      "cmd/bingo-apiserver",
			Internal: "internal/apiserver",
		},
		"admserver": {
			Cmd:      "cmd/bingo-admserver",
			Internal: "internal/admserver",
		},
		"bot": {
			Cmd:      "cmd/bingo-bot",
			Internal: "internal/bot",
		},
		"scheduler": {
			Cmd:      "cmd/bingo-scheduler",
			Internal: "internal/scheduler",
		},
		"ctl": {
			Cmd:      "cmd/bingoctl",
			Internal: "internal/bingoctl",
		},
	}
)

// CreateOptions is an option struct to support 'create' sub command.
type CreateOptions struct {
	GoVersion    string
	TemplatePath string
	RootPackage  string
	AppName      string
	AppNameCamel string

	// GitHub template
	ModuleName  string // Go module name (optional)
	TemplateRef string // Template version
	NoCache     bool   // Force re-download

	// Service selection
	All         bool     // Create all available services
	Services    []string // Explicitly specified services
	NoServices  []string // Services to exclude from defaults
	AddServices []string // Services to add to defaults
	Interactive bool     // Whether to use interactive mode (default true)

	// Git initialization
	InitGit bool // Initialize git repository (default true)

	// Build options
	Build bool // Run make build after creation (default false)

	selectedServices []string // Final computed service list (internal)
}

// NewCreateOptions returns an initialized CreateOptions instance.
func NewCreateOptions() *CreateOptions {
	return &CreateOptions{
		GoVersion: cmdutil.GetGoVersion(),
		InitGit:   true, // Initialize git by default
	}
}

// NewCmdCreate returns new initialized instance of 'create' sub command.
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

	cmd.Flags().BoolVarP(&o.All, "all", "a", false, "Create all available services (apiserver, ctl, admserver, bot, scheduler)")
	cmd.Flags().StringSliceVar(&o.Services, "services", nil, "Explicitly specify services to create (comma-separated). Use 'none' for minimal skeleton")
	cmd.Flags().StringSliceVar(&o.NoServices, "no-service", nil, "Exclude services from defaults (comma-separated)")
	cmd.Flags().StringSliceVar(&o.AddServices, "add-service", nil, "Add services to defaults (comma-separated)")

	cmd.Flags().StringVarP(&o.ModuleName, "module", "m", "",
		"Go module name (e.g., github.com/mycompany/myapp)")
	cmd.Flags().StringVarP(&o.TemplateRef, "ref", "r", "",
		"Template version (tag/branch/commit, default: recommended version)")
	cmd.Flags().BoolVar(&o.NoCache, "no-cache", false,
		"Force re-download template (for branches)")
	cmd.Flags().BoolVar(&o.InitGit, "init-git", true,
		"Initialize git repository in the created project")
	cmd.Flags().BoolVar(&o.Build, "build", false,
		"Run make build after project creation")

	return cmd
}

// Validate makes sure there is no discrepancy in command options.
func (o *CreateOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, "%s", createUsageErrStr)
	}

	// Root Package & App path
	o.RootPackage = args[0]
	arr := strings.Split(o.RootPackage, "/")
	o.AppName = strings.ToLower(arr[len(arr)-1])
	o.AppNameCamel = strcase.ToCamel(o.AppName)

	// Path is exists
	exists, err := cmdutil.PathExists(o.AppName)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	// Project exists
	console.Warn("project directory already exists!")
	prompt := promptui.Prompt{
		Label:     "Overwrite",
		IsConfirm: true,
	}

	_, err = prompt.Run()
	if err != nil {
		console.Exit("skipped.")
	}

	cmdutil.Overwrite = true

	return nil
}

// Complete completes all the required options.
func (o *CreateOptions) Complete(cmd *cobra.Command, args []string) error {
	// 1. Parse template version
	if o.TemplateRef == "" {
		o.TemplateRef = template.DefaultTemplateVersion
		fmt.Printf("Using recommended version: %s\n", o.TemplateRef)
	}

	// 2. Compute service list
	o.selectedServices = o.computeServiceList()

	// Warn if no services selected
	if len(o.selectedServices) == 0 {
		console.Warn("No services selected. Creating minimal project skeleton")
		prompt := promptui.Prompt{
			Label:     "Continue",
			IsConfirm: true,
		}
		_, err := prompt.Run()
		if err != nil {
			console.Exit("Project creation cancelled")
		}
	}

	return nil
}

// handleFetchError provides user-friendly error messages for template fetch failures
func (o *CreateOptions) handleFetchError(err error) error {
	errMsg := err.Error()
	var msg string

	// Check for network connectivity issues
	if strings.Contains(errMsg, "connection refused") {
		msg = "failed to download template: network connection refused. Please check your internet connection and try again"
	} else if strings.Contains(errMsg, "no such host") {
		msg = "failed to download template: cannot resolve domain. Please check your internet connection"
	} else if strings.Contains(errMsg, "i/o timeout") {
		msg = "failed to download template: connection timeout. Try using --no-cache to force re-download"
	} else if strings.Contains(errMsg, "404") {
		// Check for HTTP 404 errors
		msg = "failed to download template: template version not found. Check the version with -r flag (e.g., -r main)"
	} else if strings.Contains(errMsg, "403") {
		// Check for HTTP 403 errors
		msg = "failed to download template: access denied. Please check your network permissions"
	} else if _, ok := err.(net.Error); ok {
		// Generic network error
		msg = "failed to download template: network error. Please check your internet connection and try again"
	} else {
		// Default: show original error
		msg = fmt.Sprintf("failed to download template: %v", err)
	}

	console.Error(msg)
	return cmdutil.ErrExit
}

// Run executes a new sub command using the specified options.
func (o *CreateOptions) Run(args []string) error {
	fmt.Printf("Creating project '%s'...\n", o.AppName)

	// 1. Fetch template (download or use cache)
	fetcher, err := template.NewFetcher()
	if err != nil {
		return fmt.Errorf("failed to create fetcher: %w", err)
	}

	templatePath, err := fetcher.FetchTemplate(o.TemplateRef, o.NoCache)
	if err != nil {
		return o.handleFetchError(err)
	}

	// 2. Create temporary directory
	tmpDir := filepath.Join(os.TempDir(), fmt.Sprintf("bingo-%d", time.Now().Unix()))
	defer os.RemoveAll(tmpDir)

	// 3. Copy to temporary directory
	if err := cmdutil.CopyDir(templatePath, tmpDir); err != nil {
		return fmt.Errorf("failed to copy template: %w", err)
	}

	// 4. Filter services (before renaming, using original directory names)
	if len(o.selectedServices) > 0 {
		if err := o.filterServices(tmpDir); err != nil {
			return err
		}
	}

	// 5. Rename directories (always execute, only for remaining directories)
	replacer := template.NewReplacer(tmpDir, "bingo", o.ModuleName, o.AppName)
	if err := replacer.RenameDirs(); err != nil {
		return fmt.Errorf("重命名目录失败: %w", err)
	}

	// 5.1. Rename config files
	if err := replacer.RenameConfigFiles(); err != nil {
		return fmt.Errorf("重命名配置文件失败: %w", err)
	}

	// 6. Replace module name (only if -m specified)
	// Only replace if user explicitly provided a new module name
	// Otherwise, keep original module name from go.mod
	if o.ModuleName != "" {
		if err := replacer.ReplaceModuleName(); err != nil {
			return fmt.Errorf("替换模块名失败: %w", err)
		}
	}

	// 6.1. Replace app name in files (service names, paths, etc.)
	if err := replacer.ReplaceAppName(); err != nil {
		return fmt.Errorf("替换应用名失败: %w", err)
	}

	// 7. Copy .bingo.example.yaml to .bingo.yaml
	exampleConfigPath := filepath.Join(tmpDir, ".bingo.example.yaml")
	targetConfigPath := filepath.Join(tmpDir, ".bingo.yaml")

	if cmdutil.Exists(exampleConfigPath) {
		if err := cmdutil.CopyFile(exampleConfigPath, targetConfigPath); err != nil {
			return fmt.Errorf("复制 .bingo.yaml 失败: %w", err)
		}
	}

	// 8. Atomically move to target location
	// If project already exists (from Validate overwrite confirmation), remove it first
	if cmdutil.Overwrite && cmdutil.Exists(o.AppName) {
		if err := os.RemoveAll(o.AppName); err != nil {
			return fmt.Errorf("failed to remove existing project: %w", err)
		}
	}

	if err := os.Rename(tmpDir, o.AppName); err != nil {
		return fmt.Errorf("failed to move project: %w", err)
	}

	// 9. Cleanup template files (remove bingo docs and create new README)
	if err := o.cleanupTemplateFiles(o.AppName); err != nil {
		return err
	}

	// 10. Setup configuration files
	projectPath := o.AppName

	// Copy .air.example.toml to .air.toml
	airExamplePath := filepath.Join(projectPath, ".air.example.toml")
	airPath := filepath.Join(projectPath, ".air.toml")
	if cmdutil.Exists(airExamplePath) {
		cmdutil.CopyFile(airExamplePath, airPath)
	}

	// Copy config file for apiserver
	configsSrcPath := filepath.Join(projectPath, "configs", fmt.Sprintf("%s-apiserver.yaml", o.AppName))
	configsDstPath := filepath.Join(projectPath, "configs", "app-apiserver.yaml")
	if cmdutil.Exists(configsSrcPath) {
		cmdutil.CopyFile(configsSrcPath, configsDstPath)
	}

	// Copy example configs to project root for easy configuration
	if err := o.copyExampleConfigs(projectPath); err != nil {
		return err
	}

	// 11. Generate protobuf files (if applicable)
	o.runMakeProtoc(projectPath)

	// 12. Run go mod tidy
	o.runGoModTidy(projectPath)

	// 13. Initialize git repository if requested
	if o.InitGit {
		if err := o.initializeGit(projectPath); err != nil {
			console.Warn(fmt.Sprintf("Failed to initialize git repository: %v", err))
		}
	}

	// 14. Run make build if requested (after git init so Makefile can access git info)
	if o.Build {
		o.runMakeBuild(projectPath)
	}

	// Success message - show in green
	console.Info(fmt.Sprintf("Project '%s' created successfully", o.AppName))
	if len(o.selectedServices) == 0 {
		console.Warn("All services were removed. Consider running 'go mod tidy' to clean up unused dependencies")
	}

	return nil
}

// computeServiceList computes the final service list based on flags
func (o *CreateOptions) computeServiceList() []string {
	// Priority 1: --all flag
	if o.All {
		return availableServices
	}

	// Priority 2: --services flag explicitly specified
	if len(o.Services) > 0 {
		if len(o.Services) == 1 && o.Services[0] == "none" {
			return []string{}
		}
		return o.Services
	}

	// Priority 3: Start with defaults and apply modifications
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

// runMakeProtoc attempts to generate protobuf files before go mod tidy
// Returns true if protoc was run (regardless of success), false if skipped
func (o *CreateOptions) runMakeProtoc(projectPath string) bool {
	// Check if make is available
	if _, err := exec.LookPath("make"); err != nil {
		return false
	}

	// Check if Makefile exists
	makefilePath := filepath.Join(projectPath, "Makefile")
	if !cmdutil.Exists(makefilePath) {
		return false
	}

	// Check if protoc target exists in Makefile
	content, err := os.ReadFile(makefilePath)
	if err != nil {
		return false
	}
	if !strings.Contains(string(content), "protoc:") {
		return false
	}

	// Try to run make protoc
	fmt.Println("Generating protobuf files...")
	cmd := exec.Command("make", "protoc")
	cmd.Dir = projectPath
	// Suppress output - we'll handle success/failure messaging
	if err := cmd.Run(); err != nil {
		// Silently fail - user may not have protoc installed
		return true
	}

	console.Info("Protobuf files generated")
	return true
}

// runGoModTidy executes go mod tidy to clean up dependencies
// Failures are warnings only, not blocking errors
func (o *CreateOptions) runGoModTidy(projectPath string) {
	fmt.Println("Running go mod tidy...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath

	// Capture stderr for better error messaging
	var stderr strings.Builder
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errOutput := stderr.String()
		fmt.Fprint(os.Stderr, errOutput)

		// Check if it's a proto-related error
		if strings.Contains(errOutput, "proto") || strings.Contains(errOutput, "/pb") {
			console.Warn("go mod tidy failed due to missing protobuf files")
			fmt.Println()
			fmt.Println("Next steps:")
			fmt.Println("  1. Install protoc: https://grpc.io/docs/protoc-installation/")
			fmt.Println("  2. Run: cd " + o.AppName + " && make protoc")
			fmt.Println("  3. Run: go mod tidy")
		} else {
			console.Warn(fmt.Sprintf("go mod tidy failed: %v", err))
		}
	} else {
		console.Info("go mod tidy completed")
	}
}

// runMakeBuild executes make build to compile the project
// Failures are warnings only, not blocking errors
func (o *CreateOptions) runMakeBuild(projectPath string) {
	// Check if make is available
	if _, err := exec.LookPath("make"); err != nil {
		console.Warn("make not found, skipping make build")
		return
	}

	fmt.Println("Running make build...")
	cmd := exec.Command("make", "build")
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		console.Warn(fmt.Sprintf("make build failed: %v", err))
	} else {
		console.Info("make build completed")
	}
}

// initializeGit initializes a git repository in the project directory
func (o *CreateOptions) initializeGit(projectPath string) error {
	// Initialize git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize git: %w", err)
	}

	// Add all files
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add files to git: %w", err)
	}

	// Create initial commit
	cmd = exec.Command("git", "commit", "-m", "Initial commit from bingo")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create initial commit: %w", err)
	}

	return nil
}

// filterServices deletes unselected service directories
// Uses service mapping from .bingo.yaml if available, otherwise uses default mapping
func (o *CreateOptions) filterServices(targetDir string) error {
	// Try to load .bingo.yaml
	configPath := filepath.Join(targetDir, ".bingo.yaml")
	config, err := template.LoadBingoConfig(configPath)
	if err != nil {
		// If config file doesn't exist, use default mapping
		return o.filterServicesWithMapping(targetDir, defaultServiceMapping)
	}

	// Convert config services to ServiceMapping format
	mapping := make(map[string]ServiceMapping)
	for svc, info := range config.Services {
		mapping[svc] = ServiceMapping{
			Cmd:      info.Cmd,
			Internal: info.Internal,
		}
	}

	return o.filterServicesWithMapping(targetDir, mapping)
}

// copyExampleConfigs copies configs/*.example.yaml to project root as *.yaml
func (o *CreateOptions) copyExampleConfigs(projectPath string) error {
	configsDir := filepath.Join(projectPath, "configs")
	if !cmdutil.Exists(configsDir) {
		return nil
	}

	entries, err := os.ReadDir(configsDir)
	if err != nil {
		return fmt.Errorf("failed to read configs directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".example.yaml") {
			continue
		}

		// Remove .example from filename
		targetName := strings.Replace(name, ".example.yaml", ".yaml", 1)
		srcPath := filepath.Join(configsDir, name)
		dstPath := filepath.Join(projectPath, targetName)

		if err := cmdutil.CopyFile(srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to copy %s: %w", name, err)
		}
	}

	return nil
}

// cleanupTemplateFiles removes bingo template docs and creates a new README
func (o *CreateOptions) cleanupTemplateFiles(projectPath string) error {
	// Remove docs directory
	docsPath := filepath.Join(projectPath, "docs")
	if cmdutil.Exists(docsPath) {
		if err := os.RemoveAll(docsPath); err != nil {
			return fmt.Errorf("failed to remove docs directory: %w", err)
		}
	}

	// Remove CHANGELOG.md
	changelogPath := filepath.Join(projectPath, "CHANGELOG.md")
	if cmdutil.Exists(changelogPath) {
		if err := os.Remove(changelogPath); err != nil {
			return fmt.Errorf("failed to remove CHANGELOG.md: %w", err)
		}
	}

	// Remove README.zh-CN.md
	readmeZhPath := filepath.Join(projectPath, "README.zh-CN.md")
	if cmdutil.Exists(readmeZhPath) {
		if err := os.Remove(readmeZhPath); err != nil {
			return fmt.Errorf("failed to remove README.zh-CN.md: %w", err)
		}
	}

	// Create new README.md
	readmePath := filepath.Join(projectPath, "README.md")
	readmeContent := fmt.Sprintf(`# %s

Project created with [bingo](https://github.com/bingo-project/bingoctl) based on the [bingo](https://github.com/bingo-project/bingo) scaffold.

## Getting Started

1. Configure MySQL and Redis in `+"`*.yaml`"+` (copied from `+"`configs/*.example.yaml`"+`)
2. Build the project:

`+"```bash"+`
make build
`+"```"+`

3. Run the server:

`+"```bash"+`
./_output/platforms/<os>/<arch>/%s-apiserver
`+"```"+`

## Documentation

- Online documentation: https://bingoctl.dev
- GitHub: https://github.com/bingo-project/bingo
`, o.AppName, o.AppName)

	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	return nil
}

// filterServicesWithMapping deletes unselected service directories using the provided mapping
func (o *CreateOptions) filterServicesWithMapping(targetDir string, mapping map[string]ServiceMapping) error {
	// Mark selected services
	selected := make(map[string]bool)
	for _, svc := range o.selectedServices {
		selected[svc] = true
	}

	// Delete unselected service directories
	for svc, serviceMapping := range mapping {
		if !selected[svc] {
			// Delete cmd directory
			cmdPath := filepath.Join(targetDir, serviceMapping.Cmd)
			if cmdutil.Exists(cmdPath) {
				if err := os.RemoveAll(cmdPath); err != nil {
					return fmt.Errorf("删除 %s 失败: %w", serviceMapping.Cmd, err)
				}
			}

			// Delete internal directory
			internalPath := filepath.Join(targetDir, serviceMapping.Internal)
			if cmdutil.Exists(internalPath) {
				if err := os.RemoveAll(internalPath); err != nil {
					return fmt.Errorf("删除 %s 失败: %w", serviceMapping.Internal, err)
				}
			}

			// Delete bootstrap file for this service
			bootstrapPath := filepath.Join(targetDir, "internal/pkg/bootstrap", svc+".go")
			if cmdutil.Exists(bootstrapPath) {
				if err := os.Remove(bootstrapPath); err != nil {
					return fmt.Errorf("删除 %s 失败: %w", bootstrapPath, err)
				}
			}
		}
	}

	return nil
}
