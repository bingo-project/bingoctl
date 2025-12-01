// ABOUTME: Dynamic migration runner that compiles and executes user migrations
// ABOUTME: Generates temporary Go program, compiles it, and runs migration commands
package runner

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/bingo-project/bingoctl/pkg/config"
	"github.com/bingo-project/bingoctl/pkg/db"
)

//go:embed tpl/*.tpl
var tplFS embed.FS

// Runner handles dynamic compilation and execution of migrations.
type Runner struct {
	// Project info
	projectPath   string
	projectName   string
	userModule    string
	migrationDir  string
	migrationPath string

	// Database config
	dbOptions *db.MySQLOptions

	// Cache directory
	cacheDir string

	// Options
	verbose bool
	rebuild bool
}

// NewRunner creates a new migration runner.
func NewRunner(verbose, rebuild bool) (*Runner, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Read user's go.mod to get module name
	userModule, err := readModuleName(pwd)
	if err != nil {
		return nil, fmt.Errorf("failed to read go.mod: %w", err)
	}

	// Get project name from module
	projectName := filepath.Base(userModule)

	// Get migration directory from config
	migrationDir := config.Cfg.Directory.Migration
	if migrationDir == "" {
		migrationDir = "internal/apiserver/database/migration"
	}

	migrationPath := filepath.Join(pwd, migrationDir)

	// Calculate cache directory
	pathHash := CalculatePathHash(pwd)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	cacheDir := filepath.Join(homeDir, ".bingoctl", "migrator", fmt.Sprintf("%s_%s", projectName, pathHash))

	return &Runner{
		projectPath:   pwd,
		projectName:   projectName,
		userModule:    userModule,
		migrationDir:  migrationDir,
		migrationPath: migrationPath,
		dbOptions:     config.Cfg.MysqlOptions,
		cacheDir:      cacheDir,
		verbose:       verbose,
		rebuild:       rebuild,
	}, nil
}

// Run executes the migration command.
func (r *Runner) Run(command string) error {
	// Validate
	if err := r.validate(); err != nil {
		return err
	}

	// Check if rebuild is needed
	needsBuild, err := r.needsBuild()
	if err != nil {
		return fmt.Errorf("failed to check build status: %w", err)
	}

	if needsBuild || r.rebuild {
		if err := r.build(); err != nil {
			return err
		}
	}

	// Execute
	return r.execute(command)
}

func (r *Runner) validate() error {
	// Check migration directory exists
	if _, err := os.Stat(r.migrationPath); os.IsNotExist(err) {
		return fmt.Errorf("migration directory not found: %s", r.migrationPath)
	}

	// Check database config
	if r.dbOptions == nil {
		return fmt.Errorf("database configuration not found in .bingoctl.yaml")
	}

	return nil
}

func (r *Runner) needsBuild() (bool, error) {
	binaryPath := r.binaryPath()

	// Check if binary exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return true, nil
	}

	// Check checksum
	checksumFile := filepath.Join(r.cacheDir, ".checksum")
	oldChecksum, err := os.ReadFile(checksumFile)
	if err != nil {
		return true, nil
	}

	newChecksum, err := CalculateChecksum(r.migrationPath)
	if err != nil {
		return true, nil
	}

	return string(oldChecksum) != newChecksum, nil
}

func (r *Runner) build() error {
	fmt.Println("Compiling migrations...")

	// Create cache directory
	if err := os.MkdirAll(r.cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Generate main.go
	if err := r.generateMainGo(); err != nil {
		return fmt.Errorf("failed to generate main.go: %w", err)
	}

	// Generate go.mod
	if err := r.generateGoMod(); err != nil {
		return fmt.Errorf("failed to generate go.mod: %w", err)
	}

	// Run go mod tidy
	if err := r.runCommand("go", "mod", "tidy"); err != nil {
		return fmt.Errorf("go mod tidy failed: %w", err)
	}

	// Build
	binaryPath := r.binaryPath()
	if err := r.runCommand("go", "build", "-o", binaryPath, "."); err != nil {
		return fmt.Errorf("build failed (use --verbose for details): %w", err)
	}

	// Save checksum
	checksum, err := CalculateChecksum(r.migrationPath)
	if err != nil {
		return fmt.Errorf("failed to calculate checksum: %w", err)
	}
	checksumFile := filepath.Join(r.cacheDir, ".checksum")
	if err := os.WriteFile(checksumFile, []byte(checksum), 0644); err != nil {
		return fmt.Errorf("failed to save checksum: %w", err)
	}

	fmt.Println("Compilation successful.")
	return nil
}

func (r *Runner) generateMainGo() error {
	tplContent, err := tplFS.ReadFile("tpl/main.go.tpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("main").Parse(string(tplContent))
	if err != nil {
		return err
	}

	data := map[string]string{
		"MigrationImport": r.userModule + "/" + r.migrationDir,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	mainPath := filepath.Join(r.cacheDir, "main.go")
	return os.WriteFile(mainPath, buf.Bytes(), 0644)
}

func (r *Runner) generateGoMod() error {
	tplContent, err := tplFS.ReadFile("tpl/go.mod.tpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("gomod").Parse(string(tplContent))
	if err != nil {
		return err
	}

	// Get Go version
	goVersion := strings.TrimPrefix(runtime.Version(), "go")
	// Extract major.minor only (e.g., "1.21" from "1.21.5")
	parts := strings.Split(goVersion, ".")
	if len(parts) >= 2 {
		goVersion = parts[0] + "." + parts[1]
	}

	// Get bingoctl version from go.mod
	bingoctlVersion := getBingoctlVersion()

	data := map[string]string{
		"GoVersion":       goVersion,
		"UserModule":      r.userModule,
		"UserModulePath":  r.projectPath,
		"BingoctlVersion": bingoctlVersion,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	modPath := filepath.Join(r.cacheDir, "go.mod")
	return os.WriteFile(modPath, buf.Bytes(), 0644)
}

func (r *Runner) runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = r.cacheDir

	if r.verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			if stderr.Len() > 0 {
				return fmt.Errorf("%w\n%s", err, stderr.String())
			}
			return err
		}
		return nil
	}

	return cmd.Run()
}

func (r *Runner) execute(command string) error {
	binaryPath := r.binaryPath()

	cmd := exec.Command(binaryPath, command,
		"--host", r.dbOptions.Host,
		"--username", r.dbOptions.Username,
		"--password", r.dbOptions.Password,
		"--database", r.dbOptions.Database,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (r *Runner) binaryPath() string {
	name := "migrator"
	if runtime.GOOS == "windows" {
		name = "migrator.exe"
	}
	return filepath.Join(r.cacheDir, name)
}

// readModuleName reads the module name from go.mod in the given directory.
func readModuleName(dir string) (string, error) {
	modPath := filepath.Join(dir, "go.mod")
	content, err := os.ReadFile(modPath)
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	return "", fmt.Errorf("module name not found in go.mod")
}

// getBingoctlVersion returns the current bingoctl version.
// In production, this would read from a version constant or go.mod.
func getBingoctlVersion() string {
	// Try to get from go.mod of bingoctl itself
	// For now, return a version that will work with replace directive
	return "v1.5.0"
}
