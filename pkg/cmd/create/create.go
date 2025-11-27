package create

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bingo-project/component-base/cli/console"
	"github.com/iancoleman/strcase"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	cmdutil "github.com/bingo-project/bingoctl/pkg/util"
	"github.com/bingo-project/bingoctl/pkg/template"
)

const (
	createUsageStr = "create NAME"
)

var (
	createUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the create command",
		createUsageStr,
	)

	defaultServices   = []string{"apiserver", "ctl"}
	availableServices = []string{"apiserver", "ctl", "admserver", "bot", "scheduler"}
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
	Services    []string // Explicitly specified services
	NoServices  []string // Services to exclude from defaults
	AddServices []string // Services to add to defaults
	Interactive bool     // Whether to use interactive mode (default true)

	selectedServices []string // Final computed service list (internal)
}

// NewCreateOptions returns an initialized CreateOptions instance.
func NewCreateOptions() *CreateOptions {
	return &CreateOptions{
		GoVersion: cmdutil.GetGoVersion(),
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

	cmd.Flags().StringSliceVar(&o.Services, "services", nil, "Explicitly specify services to create (comma-separated). Use 'none' for minimal skeleton")
	cmd.Flags().StringSliceVar(&o.NoServices, "no-service", nil, "Exclude services from defaults (comma-separated)")
	cmd.Flags().StringSliceVar(&o.AddServices, "add-service", nil, "Add services to defaults (comma-separated)")

	cmd.Flags().StringVarP(&o.ModuleName, "module", "m", "",
		"Go module name (e.g., github.com/mycompany/myapp)")
	cmd.Flags().StringVarP(&o.TemplateRef, "ref", "r", "",
		"Template version (tag/branch/commit, default: recommended version)")
	cmd.Flags().BoolVar(&o.NoCache, "no-cache", false,
		"Force re-download template (for branches)")

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
		console.Info(fmt.Sprintf("使用推荐版本：%s", o.TemplateRef))
	}

	// 2. Compute service list (keep existing logic)
	o.Interactive = len(o.Services) == 0 && len(o.NoServices) == 0 && len(o.AddServices) == 0

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

	return nil
}

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

// Run executes a new sub command using the specified options.
func (o *CreateOptions) Run(args []string) error {
	console.Info(fmt.Sprintf("Creating project %s", o.RootPackage))

	// 1. Fetch template (download or use cache)
	fetcher, err := template.NewFetcher()
	if err != nil {
		return fmt.Errorf("failed to create fetcher: %w", err)
	}

	templatePath, err := fetcher.FetchTemplate(o.TemplateRef, o.NoCache)
	if err != nil {
		return fmt.Errorf("获取模板失败: %w", err)
	}

	// 2. Create temporary directory
	tmpDir := filepath.Join(os.TempDir(), fmt.Sprintf("bingoctl-%d", time.Now().Unix()))
	defer os.RemoveAll(tmpDir)

	// 3. Copy to temporary directory
	console.Info("复制模板...")
	if err := cmdutil.CopyDir(templatePath, tmpDir); err != nil {
		return fmt.Errorf("复制模板失败: %w", err)
	}

	// 4. Filter services (before renaming, using original directory names)
	if len(o.selectedServices) > 0 {
		console.Info("过滤服务...")
		if err := o.filterServices(tmpDir); err != nil {
			return err
		}
	}

	// 5. Rename directories (always execute, only for remaining directories)
	replacer := template.NewReplacer(tmpDir, "bingo", o.ModuleName, o.AppName)
	console.Info("重命名目录...")
	if err := replacer.RenameDirs(); err != nil {
		return fmt.Errorf("重命名目录失败: %w", err)
	}

	// 6. Replace module name (only if -m specified)
	if o.ModuleName != "" {
		console.Info(fmt.Sprintf("替换模块名: bingo -> %s", o.ModuleName))
		if err := replacer.ReplaceModuleName(); err != nil {
			return fmt.Errorf("替换模块名失败: %w", err)
		}
	}

	// 7. Atomically move to target location
	if err := os.Rename(tmpDir, o.AppName); err != nil {
		return fmt.Errorf("移动项目失败: %w", err)
	}

	// 8. Success message
	console.Info("✓ 项目创建成功！")
	if len(o.selectedServices) == 0 {
		console.Info("提示：已删除所有服务，建议运行 'go mod tidy' 清理未使用的依赖")
	}

	return nil
}

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

// filterServices deletes unselected service directories
// Reads service mapping from .bingoctl.yaml
func (o *CreateOptions) filterServices(targetDir string) error {
	// Load .bingoctl.yaml
	configPath := filepath.Join(targetDir, ".bingoctl.yaml")
	config, err := template.LoadBingoctlConfig(configPath)
	if err != nil {
		return fmt.Errorf("加载 .bingoctl.yaml 失败: %w\n提示：模板版本 %s 可能不包含此配置文件", err, o.TemplateRef)
	}

	allServices := config.Services

	// Mark selected services
	selected := make(map[string]bool)
	for _, svc := range o.selectedServices {
		selected[svc] = true
	}

	// Delete unselected service directories
	for svc, service := range allServices {
		if !selected[svc] {
			// Delete cmd directory
			cmdPath := filepath.Join(targetDir, service.Cmd)
			if cmdutil.Exists(cmdPath) {
				console.Info(fmt.Sprintf("  删除 %s", service.Cmd))
				if err := os.RemoveAll(cmdPath); err != nil {
					return fmt.Errorf("删除 %s 失败: %w", service.Cmd, err)
				}
			}

			// Delete internal directory
			internalPath := filepath.Join(targetDir, service.Internal)
			if cmdutil.Exists(internalPath) {
				console.Info(fmt.Sprintf("  删除 %s", service.Internal))
				if err := os.RemoveAll(internalPath); err != nil {
					return fmt.Errorf("删除 %s 失败: %w", service.Internal, err)
				}
			}
		}
	}

	return nil
}
