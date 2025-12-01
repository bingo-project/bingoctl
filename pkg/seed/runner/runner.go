// ABOUTME: Dynamic seed runner that compiles and executes user seeders
// ABOUTME: Generates temporary Go program in project dir, compiles it, and runs seed commands
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

// Runner handles dynamic compilation and execution of seeders.
type Runner struct {
	projectPath string
	projectName string
	userModule  string
	seederDir   string
	seederPath  string

	dbOptions *db.MySQLOptions

	cacheDir string
	tmpDir   string

	verbose bool
	rebuild bool
}

// NewRunner creates a new seed runner.
func NewRunner(verbose, rebuild bool) (*Runner, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	userModule, err := readModuleName(pwd)
	if err != nil {
		return nil, fmt.Errorf("failed to read go.mod: %w", err)
	}

	projectName := filepath.Base(userModule)

	seederDir := config.Cfg.Directory.Seeder
	if seederDir == "" {
		seederDir = "internal/pkg/database/seeder"
	}

	seederPath := filepath.Join(pwd, seederDir)

	pathHash := CalculatePathHash(pwd)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	cacheDir := filepath.Join(homeDir, ".bingo", "seeder", fmt.Sprintf("%s_%s", projectName, pathHash))

	tmpDir := filepath.Join(pwd, ".bingo_tmp")

	return &Runner{
		projectPath: pwd,
		projectName: projectName,
		userModule:  userModule,
		seederDir:   seederDir,
		seederPath:  seederPath,
		dbOptions:   config.Cfg.MysqlOptions,
		cacheDir:    cacheDir,
		tmpDir:      tmpDir,
		verbose:     verbose,
		rebuild:     rebuild,
	}, nil
}

// Run executes the seed command.
func (r *Runner) Run(seederName string) error {
	if err := r.validate(); err != nil {
		return err
	}

	needsBuild, err := r.needsBuild()
	if err != nil {
		return fmt.Errorf("failed to check build status: %w", err)
	}

	if needsBuild || r.rebuild {
		if err := r.build(); err != nil {
			return err
		}
	}

	return r.execute(seederName)
}

func (r *Runner) validate() error {
	if _, err := os.Stat(r.seederPath); os.IsNotExist(err) {
		return fmt.Errorf("seeder directory not found: %s", r.seederPath)
	}

	if r.dbOptions == nil {
		return fmt.Errorf("database configuration not found in .bingo.yaml")
	}

	return nil
}

func (r *Runner) needsBuild() (bool, error) {
	binaryPath := r.binaryPath()

	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return true, nil
	}

	checksumFile := filepath.Join(r.cacheDir, ".checksum")
	oldChecksum, err := os.ReadFile(checksumFile)
	if err != nil {
		return true, nil
	}

	newChecksum, err := CalculateChecksum(r.seederPath)
	if err != nil {
		return true, nil
	}

	return string(oldChecksum) != newChecksum, nil
}

func (r *Runner) build() error {
	fmt.Println("Compiling seeders...")

	if err := os.MkdirAll(r.cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}
	if err := os.MkdirAll(r.tmpDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	defer os.RemoveAll(r.tmpDir)

	if err := r.generateMainGo(); err != nil {
		return fmt.Errorf("failed to generate main.go: %w", err)
	}

	binaryPath := r.binaryPath()
	if err := r.runBuildCommand(binaryPath); err != nil {
		return fmt.Errorf("build failed (use --verbose for details): %w", err)
	}

	checksum, err := CalculateChecksum(r.seederPath)
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
		"SeederImport": r.userModule + "/" + r.seederDir,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	mainPath := filepath.Join(r.tmpDir, "main.go")
	return os.WriteFile(mainPath, buf.Bytes(), 0644)
}

func (r *Runner) runBuildCommand(outputPath string) error {
	cmd := exec.Command("go", "build", "-o", outputPath, "./.bingo_tmp")
	cmd.Dir = r.projectPath

	if r.verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

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

func (r *Runner) execute(seederName string) error {
	binaryPath := r.binaryPath()

	args := []string{
		"--host", r.dbOptions.Host,
		"--username", r.dbOptions.Username,
		"--password", r.dbOptions.Password,
		"--database", r.dbOptions.Database,
	}

	if seederName != "" {
		args = append(args, "--seeder", seederName)
	}

	cmd := exec.Command(binaryPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (r *Runner) binaryPath() string {
	name := "seeder"
	if runtime.GOOS == "windows" {
		name = "seeder.exe"
	}
	return filepath.Join(r.cacheDir, name)
}

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
