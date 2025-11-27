# bingoctl create 命令重构设计文档

> **注意**：此设计取代了之前的 `2025-11-27-part2-template-reorganization.md` 计划。
> 旧计划基于重组 `embed.FS` 模板，新设计直接从 GitHub 拉取 bingo 项目，更简单、更灵活。

## 背景

bingoctl 是 bingo 脚手架的 CLI 工具，`create` 命令用于创建新项目。目前模板是通过 `embed.FS` 嵌入到二进制文件中的，导致以下问题：

- 模板更新需要重新编译 bingoctl
- 无法灵活选择 bingo 的不同版本
- 维护两份代码（bingo 主仓库 + bingoctl 嵌入模板）

## 设计目标

改为从 GitHub 拉取指定版本的 bingo 项目作为模板，实现：

1. 支持指定版本（tag/branch/commit）创建项目
2. 默认使用推荐稳定版本（硬编码）
3. 本地缓存机制，避免重复下载
4. 支持可选的包名替换（从 `bingo` 改为用户指定的模块名）
5. 简化模板逻辑，不使用 `.tpl` 文件
6. 不依赖 GitHub API，避免限流问题
7. 支持镜像加速
8. 使用临时目录，确保原子操作

## 设计变更（相对于初始版本）

基于 code review 反馈，以下是关键设计变更：

### 1. 版本管理策略

- **变更前**：通过 GitHub API 获取最新 stable tag
- **变更后**：硬编码推荐版本（如 `v1.0.0`），每次 bingoctl 发布时手动更新
- **原因**：避免 GitHub API 限流，提高可靠性和速度

### 2. 分支缓存策略

- **变更前**：分支缓存使用 `main-{commit-hash}` 格式
- **变更后**：分支缓存使用分支名（如 `main`），通过 `--no-cache` 标志强制更新
- **原因**：简化缓存逻辑，用户体验更直观

### 3. 包名替换范围

- **变更前**：只替换 `*.go` 和 `go.mod`
- **变更后**：使用文件类型白名单，包括 Makefile, Dockerfile, 配置文件等
- **原因**：确保替换完整，避免遗漏

### 4. 服务过滤实现

- **变更前**：硬编码服务映射
- **变更后**：临时使用硬编码 + TODO，未来迁移到 `.bingoctl.yaml` 元数据文件
- **原因**：降低维护成本，适应 bingo 项目结构变化

### 5. 操作原子性

- **变更前**：直接在目标目录操作
- **变更后**：在临时目录完成所有操作，最后原子移动到目标位置
- **原因**：避免中间状态污染，失败时不影响用户目录

### 6. 镜像支持

- **新增功能**：支持 `BINGOCTL_TEMPLATE_MIRROR` 环境变量
- **原因**：解决某些网络环境下 GitHub 访问慢或不可达的问题

### 7. 错误处理

- **变更前**：简单错误提示
- **变更后**：友好且具体的错误提示，包含原因和建议
- **原因**：提升用户体验，快速定位问题

### 8. 并发安全

- **新增功能**：使用文件锁保证并发安全
- **原因**：支持多个 bingoctl 实例同时运行

### 9. 用户体验优化

- **新增功能**：下载进度条、缓存目录权限检查、友好提示
- **原因**：提升用户体验

### 10. 测试策略

- **变更**：测试部分先备注，优先实现功能，测试后补
- **原因**：快速交付核心功能

## 设计决策

### 1. 命令行参数

```bash
# 使用默认推荐版本，不修改包名
bingoctl create demo

# 指定版本
bingoctl create demo -r v1.2.3
bingoctl create demo -r main

# 指定包名（自动替换）
bingoctl create demo -m github.com/mycompany/demo

# 组合使用
bingoctl create demo -m github.com/mycompany/demo -r v1.2.3

# 强制更新缓存（用于分支）
bingoctl create demo -r main --no-cache

# 服务选择（保留现有功能）
bingoctl create demo --services apiserver,ctl
bingoctl create demo --no-service bot,scheduler
```

**参数说明：**
- `NAME`：项目目录名（必需）
- `-m, --module`：Go 模块名（可选），指定后会替换所有 `bingo` 包名引用
- `-r, --ref`：模板版本（可选），支持 tag/branch/commit，默认为硬编码的推荐版本
- `--no-cache`：强制重新下载，不使用缓存（主要用于分支）

### 2. 版本管理

**设计原则：避免依赖 GitHub API，防止限流问题。**

- **默认推荐版本**：在 bingoctl 中硬编码一个推荐版本（如 `v1.0.0`），每次 bingoctl 发布时更新
  - 优点：不需要网络请求，快速可靠
  - 用户不指定 `-r` 时使用此版本

- **固定版本**：用户通过 `-r` 指定具体的 tag 或 commit hash
  - 示例：`-r v1.2.3` 或 `-r abc123def`

- **分支支持**：支持 `main` 等分支，使用 `--no-cache` 强制更新
  - 分支直接使用分支名作为缓存 key（如 `main`）
  - 如果需要最新代码，使用 `--no-cache` 标志

- **版本过滤**：只推荐使用正式发布版本（如 `v1.2.3`），不包含 pre-release 版本（如 `v1.2.3-beta.1`）
  - 文档中明确说明推荐版本的选择标准

### 3. 缓存策略

```
~/.bingoctl/templates/
  ├── v1.2.3/              # tag 版本缓存（永久）
  ├── v1.2.4/
  └── main/                # 分支缓存（可通过 --no-cache 更新）
```

- **Tag 版本**：永久缓存，认为内容不变
- **分支版本**：也缓存，但可通过 `--no-cache` 标志强制更新
- **缓存命中**：如果缓存存在且未指定 `--no-cache`，直接使用，不重新下载
- **并发安全**：使用文件锁防止多个 bingoctl 实例同时写缓存
- **权限**：缓存目录权限为 0755

### 4. 下载方式

使用 GitHub Archive URL 下载 tarball：
```
# Tag
https://github.com/bingo-project/bingo/archive/refs/tags/{ref}.tar.gz

# Branch
https://github.com/bingo-project/bingo/archive/refs/heads/{ref}.tar.gz
```

**优点：**
- 快速，体积小（不包含 .git）
- 不依赖系统 git 命令
- Go 标准库可处理 tar.gz
- 不使用 GitHub API，无限流问题

**镜像支持：**
- 支持通过环境变量 `BINGOCTL_TEMPLATE_MIRROR` 配置镜像地址
- 示例：`export BINGOCTL_TEMPLATE_MIRROR=https://ghproxy.com/`
- 实际下载 URL：`${MIRROR}https://github.com/...`

**超时设置：** 30 秒

**下载进度：** 显示进度条，提升用户体验

### 5. 模板处理流程

**不使用 `.tpl` 模板文件**，而是直接处理真实项目：

```
1. 下载并缓存 bingo 项目
   - 下载 tarball 到临时文件
   - 解压到缓存目录（处理 tarball 根目录：bingo-{ref}/）
   - 使用文件锁保证并发安全
   ↓
2. 复制到临时目录（非用户目标目录）
   - 使用 os.TempDir() 创建临时目录
   - 复制缓存内容到临时目录
   ↓
3. 过滤服务（根据服务选择参数）
   - 删除未选中的服务的 cmd/ 和 internal/ 目录
   - 保留 internal/pkg/（共享代码）
   - 例如：只选择 apiserver 时删除 cmd/bingo-admserver, internal/admserver 等
   ↓
4. 重命名目录（总是执行）
   - cmd/bingo-apiserver → cmd/{app}-apiserver
   - cmd/bingo-admserver → cmd/{app}-admserver
   - cmd/bingo-bot → cmd/{app}-bot
   - cmd/bingo-scheduler → cmd/{app}-scheduler
   - cmd/bingoctl → cmd/{app}ctl
   - 只重命名这些明确列出的目录
   ↓
5. 替换包名（仅当指定 -m 时）
   - go.mod: module bingo → module {newModule}
   - *.go: import "bingo/xxx" → import "{newModule}/xxx"
   - *.go: 字符串中的 "bingo/" → "{newModule}/"
   - Makefile, Dockerfile, *.yaml, *.toml, *.json, *.sh: 文本替换
   - 完整文件类型列表见下文
   ↓
6. 原子移动到目标位置
   - os.Rename(tmpDir, targetDir)
   - 如果失败，清理临时目录
   - 成功后提示用户运行 go mod tidy（如有需要）
```

**Tarball 解压注意事项：**
- GitHub tarball 解压后有根目录（如 `bingo-v1.2.3/` 或 `bingo-main/`）
- 需要检测并进入根目录，提取其中的内容到缓存

### 6. 错误处理

提供友好且具体的错误提示：

- **网络不可用**：
  ```
  错误：无法下载模板

  可能的原因：
  1. 网络连接失败
  2. GitHub 不可访问

  建议：
  1. 检查网络连接
  2. 配置镜像：export BINGOCTL_TEMPLATE_MIRROR=https://ghproxy.com/
  3. 或稍后重试
  ```

- **版本不存在**：
  ```
  错误：模板版本 'v999.0.0' 不存在

  建议：
  1. 访问 https://github.com/bingo-project/bingo/tags 查看可用版本
  2. 或不指定 -r 参数使用默认推荐版本
  ```

- **下载超时**：
  ```
  错误：下载超时（30秒）

  建议：
  1. 检查网络连接速度
  2. 配置镜像加速下载
  3. 或稍后重试
  ```

- **磁盘空间不足**：
  ```
  错误：磁盘空间不足

  缓存目录：~/.bingoctl/templates/
  建议：清理缓存或释放磁盘空间
  ```

- **目标目录已存在**：
  ```
  错误：目录 'myapp' 已存在

  建议：
  1. 使用其他目录名
  2. 删除现有目录
  3. 或使用 --force 覆盖（如果实现该选项）
  ```

## 技术实现

### 核心模块

#### 1. `pkg/template/fetcher.go`

模板获取器，负责下载和缓存。

```go
type Fetcher struct {
    cacheDir string        // ~/.bingoctl/templates
    timeout  time.Duration // 30s
    mirror   string        // 镜像地址（从环境变量读取）
}

// FetchTemplate 下载模板到缓存（如不存在），返回缓存路径
// noCache: 是否强制重新下载
func (f *Fetcher) FetchTemplate(ref string, noCache bool) (string, error)

// downloadWithTimeout 下载 tarball，30秒超时，显示进度条
func (f *Fetcher) downloadWithTimeout(url string) (string, error)

// extractTarball 解压 tarball 到缓存目录
// 处理 GitHub tarball 根目录（如 bingo-v1.2.3/）
func (f *Fetcher) extractTarball(tarPath, destDir string) error

// buildDownloadURL 构建下载 URL（支持镜像）
// 示例：tag: https://github.com/.../archive/refs/tags/v1.2.3.tar.gz
//      branch: https://github.com/.../archive/refs/heads/main.tar.gz
func (f *Fetcher) buildDownloadURL(ref string) string

// acquireLock 获取文件锁，保证并发安全
func (f *Fetcher) acquireLock() (*flock.Flock, error)
```

#### 2. `pkg/template/version.go`

版本管理，定义默认推荐版本。

```go
// DefaultTemplateVersion 默认推荐的模板版本
// 每次 bingoctl 发布时更新此版本号
const DefaultTemplateVersion = "v1.0.0"

// isValidRef 检查 ref 格式是否有效
// 支持：v1.2.3, main, abc123def 等
func isValidRef(ref string) bool

// refType 返回 ref 类型：tag, branch, commit
func refType(ref string) string
```

**设计说明：**
- 不使用 GitHub API，避免限流和网络依赖
- 硬编码推荐版本，每次 bingoctl 发布时手动更新
- 用户可通过 `-r` 参数指定任意版本

#### 3. `pkg/template/replacer.go`

包名和目录名替换器。

```go
type Replacer struct {
    targetDir string  // 目标目录
    oldModule string  // "bingo"
    newModule string  // "github.com/mycompany/demo"
    appName   string  // "demo"
}

// ReplaceModuleName 替换所有文件中的模块名
// 遍历目标目录，根据文件扩展名进行替换
func (r *Replacer) ReplaceModuleName() error

// replaceInFile 替换单个文件中的模块名
// 使用字符串替换，避免误伤二进制文件
func (r *Replacer) replaceInFile(path string) error

// RenameDirs 重命名目录（明确列表）
func (r *Replacer) RenameDirs() error

// shouldReplaceFile 判断文件是否需要替换
// 根据文件扩展名白名单
func (r *Replacer) shouldReplaceFile(path string) bool
```

**替换文件类型（白名单）：**

```go
var replaceableExtensions = []string{
    // Go 相关
    ".go", ".mod", ".sum",

    // 文档
    ".md", ".txt",

    // 构建和脚本
    "Makefile", ".mk", ".sh", ".bash",

    // 配置文件
    ".yaml", ".yml", ".toml", ".json",

    // Docker
    "Dockerfile", ".dockerignore",
}
```

**替换规则：**

| 原始内容 | 替换后 | 文件类型 |
|---------|--------|---------|
| `module bingo` | `module {newModule}` | go.mod |
| `import "bingo/xxx"` | `import "{newModule}/xxx"` | *.go |
| `"bingo/internal/xxx"` | `"{newModule}/internal/xxx"` | *.go (字符串) |
| `bingo/cmd/...` | `{newModule}/cmd/...` | Makefile |
| `WORKDIR /app/bingo` | `WORKDIR /app/{appName}` | Dockerfile |
| 任何文本文件中的 `bingo` | `{newModule}` 或 `{appName}` | 配置文件、脚本、文档 |

**目录重命名规则（明确列表）：**

```go
var renameRules = map[string]string{
    "cmd/bingo-apiserver":   "cmd/{app}-apiserver",
    "cmd/bingo-admserver":   "cmd/{app}-admserver",
    "cmd/bingo-bot":         "cmd/{app}-bot",
    "cmd/bingo-scheduler":   "cmd/{app}-scheduler",
    "cmd/bingoctl":          "cmd/{app}ctl",
}
```

#### 4. 修改 `pkg/cmd/create/create.go`

移除 `embed.FS`，使用新的 template 模块。

```go
type CreateOptions struct {
    AppName      string   // 目录名
    ModuleName   string   // Go 模块名（可选）
    TemplateRef  string   // 模板版本
    NoCache      bool     // 是否强制重新下载
    GoVersion    string

    // Service selection（保留现有字段）
    Services     []string
    NoServices   []string
    AddServices  []string
    Interactive  bool
    selectedServices []string
}

func NewCmdCreate() *cobra.Command {
    cmd.Flags().StringVarP(&o.ModuleName, "module", "m", "",
        "Go module name (e.g., github.com/mycompany/myapp)")
    cmd.Flags().StringVarP(&o.TemplateRef, "ref", "r", "",
        "Template version (tag/branch/commit, default: recommended version)")
    cmd.Flags().BoolVar(&o.NoCache, "no-cache", false,
        "Force re-download template (for branches)")

    // 保留服务选择 flags
    cmd.Flags().StringSliceVar(&o.Services, "services", nil, "...")
    cmd.Flags().StringSliceVar(&o.NoServices, "no-service", nil, "...")
    cmd.Flags().StringSliceVar(&o.AddServices, "add-service", nil, "...")

    return cmd
}

func (o *CreateOptions) Complete(cmd *cobra.Command, args []string) error {
    // 1. 解析模板版本
    if o.TemplateRef == "" {
        o.TemplateRef = template.DefaultTemplateVersion
        console.Info(fmt.Sprintf("使用推荐版本: %s", o.TemplateRef))
    }

    // 2. 计算服务列表（保留现有逻辑）
    if o.Interactive {
        o.selectedServices = o.selectServicesInteractively()
    } else {
        o.selectedServices = o.computeServiceList()
    }

    return nil
}

func (o *CreateOptions) Run(args []string) error {
    // 1. 获取模板（下载或使用缓存）
    fetcher := template.NewFetcher()
    templatePath, err := fetcher.FetchTemplate(o.TemplateRef, o.NoCache)
    if err != nil {
        return fmt.Errorf("获取模板失败: %w", err)
    }

    // 2. 创建临时目录
    tmpDir := filepath.Join(os.TempDir(), fmt.Sprintf("bingoctl-%d", time.Now().Unix()))
    defer os.RemoveAll(tmpDir)

    // 3. 复制到临时目录
    if err := copyDir(templatePath, tmpDir); err != nil {
        return fmt.Errorf("复制模板失败: %w", err)
    }

    // 4. 过滤服务（在重命名之前删除未选中的服务目录）
    if len(o.selectedServices) < len(availableServices) {
        console.Info("过滤服务...")
        if err := o.filterServices(tmpDir); err != nil {
            return err
        }
    }

    // 5. 重命名目录（总是执行）
    replacer := template.NewReplacer(tmpDir, "bingo", o.ModuleName, o.AppName)
    if err := replacer.RenameDirs(); err != nil {
        return fmt.Errorf("重命名目录失败: %w", err)
    }

    // 6. 替换模块名（仅当指定 -m 时）
    if o.ModuleName != "" {
        console.Info(fmt.Sprintf("替换模块名: bingo -> %s", o.ModuleName))
        if err := replacer.ReplaceModuleName(); err != nil {
            return fmt.Errorf("替换模块名失败: %w", err)
        }
    }

    // 7. 原子移动到目标位置
    if err := os.Rename(tmpDir, o.AppName); err != nil {
        return fmt.Errorf("移动项目失败: %w", err)
    }

    // 8. 提示用户后续操作
    console.Success("项目创建成功！")
    if len(o.selectedServices) == 0 {
        console.Info("提示：已删除所有服务，建议运行 'go mod tidy' 清理未使用的依赖")
    }

    return nil
}

// filterServices 删除未选中的服务目录
func (o *CreateOptions) filterServices(targetDir string) error {
    // 方案 A：从 bingo 项目的 .bingoctl.yaml 读取服务映射（推荐）
    // 方案 B：约定优于配置，自动扫描 cmd/ 和 internal/ 目录

    // 这里采用方案 B（简单实现）
    // TODO: 未来迁移到方案 A，在 bingo 项目中添加 .bingoctl.yaml

    // 硬编码服务映射（临时方案）
    allServices := map[string]struct {
        cmdDir      string
        internalDir string
    }{
        "apiserver": {"cmd/bingo-apiserver", "internal/apiserver"},
        "admserver": {"cmd/bingo-admserver", "internal/admserver"},
        "bot":       {"cmd/bingo-bot", "internal/bot"},
        "scheduler": {"cmd/bingo-scheduler", "internal/scheduler"},
        "ctl":       {"cmd/bingoctl", "internal/bingoctl"},
    }

    // 标记选中的服务
    selected := make(map[string]bool)
    for _, svc := range o.selectedServices {
        selected[svc] = true
    }

    // 删除未选中的服务目录
    for svc, dirs := range allServices {
        if !selected[svc] {
            // 删除 cmd 目录
            cmdPath := filepath.Join(targetDir, dirs.cmdDir)
            if exists(cmdPath) {
                console.Info(fmt.Sprintf("  删除 %s", dirs.cmdDir))
                if err := os.RemoveAll(cmdPath); err != nil {
                    return fmt.Errorf("删除 %s 失败: %w", dirs.cmdDir, err)
                }
            }

            // 删除 internal 目录
            internalPath := filepath.Join(targetDir, dirs.internalDir)
            if exists(internalPath) {
                console.Info(fmt.Sprintf("  删除 %s", dirs.internalDir))
                if err := os.RemoveAll(internalPath); err != nil {
                    return fmt.Errorf("删除 %s 失败: %w", dirs.internalDir, err)
                }
            }
        }
    }

    return nil
}
```

**服务过滤改进方案（未来）：**

在 bingo 项目中添加 `.bingoctl.yaml`：

```yaml
# bingo/.bingoctl.yaml
version: 1
services:
  apiserver:
    cmd: cmd/bingo-apiserver
    internal: internal/apiserver
    description: API 服务器
  admserver:
    cmd: cmd/bingo-admserver
    internal: internal/admserver
    description: 管理后台服务器
  bot:
    cmd: cmd/bingo-bot
    internal: internal/bot
    description: Bot 服务
  scheduler:
    cmd: cmd/bingo-scheduler
    internal: internal/scheduler
    description: 定时任务调度器
  ctl:
    cmd: cmd/bingoctl
    internal: internal/bingoctl
    description: 命令行工具
```

优点：
- 消除硬编码，bingo 项目结构变化时无需修改 bingoctl
- 可以添加更多元数据（描述、依赖等）
- 支持更灵活的服务管理

### 工具函数

`pkg/util/file.go` 中需要添加：

```go
// copyDir 递归复制目录
func copyDir(src, dst string) error

// exists 检查路径是否存在
func exists(path string) bool
```

## 项目结构变化

```diff
bingoctl/
├── pkg/
│   ├── cmd/
│   │   └── create/
-│   │       ├── tpl/                    # 删除：嵌入的模板目录
│   │       ├── create.go               # 修改：移除 embed.FS
│   │       └── create_test.go
+│   ├── template/                       # 新增：模板处理模块
+│   │   ├── fetcher.go                 # 下载和缓存
+│   │   ├── version.go                 # 版本解析
+│   │   └── replacer.go                # 替换逻辑
│   └── util/
│       └── file.go                      # 添加文件操作工具函数
```

## 使用示例

### 场景 1: 快速创建项目（默认配置）

```bash
bingoctl create myapp
```

- 使用最新稳定版本
- 不修改包名（保持 `module bingo`）
- 目录名改为 `myapp-apiserver`、`myappctl` 等
- 包含默认服务（apiserver, ctl）

### 场景 2: 指定版本创建

```bash
bingoctl create myapp -r v1.2.0
```

- 使用指定版本 v1.2.0
- 其他同场景 1

### 场景 3: 自定义包名

```bash
bingoctl create myapp -m github.com/mycompany/myapp
```

生成的项目：
- 目录：`myapp/`
- `go.mod`: `module github.com/mycompany/myapp`
- 所有 import: `import "github.com/mycompany/myapp/internal/xxx"`
- 目录名：`cmd/myapp-apiserver/`, `cmd/myappctl/`

### 场景 4: 完整自定义

```bash
bingoctl create myapp \
  -m github.com/mycompany/myapp \
  -r v1.2.0 \
  --services apiserver,ctl
```

- 版本：v1.2.0
- 包名：github.com/mycompany/myapp
- 只包含 apiserver 和 ctl 服务

### 场景 5: 单一服务项目

```bash
# 只创建 API 服务
bingoctl create api-only --services apiserver

# 只创建命令行工具
bingoctl create cli-only --services ctl
```

生成的项目：
- 只包含选中服务的 `cmd/` 和 `internal/` 目录
- 保留 `internal/pkg/`（共享代码）
- 其他服务目录被删除

### 场景 6: 排除某些服务

```bash
# 创建项目，但不包括 bot 和 scheduler
bingoctl create myapp --no-service bot,scheduler
```

- 从默认服务（apiserver, ctl）开始
- 排除指定的服务
- 结果：只包含 apiserver 和 ctl

### 场景 7: 添加额外服务

```bash
# 在默认基础上添加 admserver
bingoctl create myapp --add-service admserver
```

- 从默认服务（apiserver, ctl）开始
- 添加 admserver
- 结果：包含 apiserver, ctl, admserver

### 场景 8: 最小项目骨架

```bash
bingoctl create minimal --services none
```

- 不包含任何服务
- 只有基础项目结构（go.mod, Makefile, configs/ 等）
- 适合从零开始构建自定义服务

### 场景 9: 使用主分支

```bash
bingoctl create myapp -r main
```

- 使用 main 分支最新代码
- 缓存为 `main-{commit-hash}`

## 向后兼容性

**破坏性变更：**
- 移除嵌入的模板，首次运行需要网络连接

**保留功能：**
- 服务选择参数（`--services`, `--no-service`, `--add-service`）
- 交互式服务选择
- 项目覆盖确认

## 测试计划

> **注意**：测试部分先备注，优先实现功能。测试将在功能稳定后补充。

### 单元测试（待实现）

1. `template/version.go`
   - 测试 ref 格式验证
   - 测试 ref 类型判断

2. `template/fetcher.go`
   - 测试缓存命中逻辑
   - 测试下载和解压
   - 测试超时处理
   - 测试 Tarball 根目录处理
   - 测试并发安全（文件锁）
   - 测试镜像支持

3. `template/replacer.go`
   - 测试模块名替换
   - 测试目录重命名
   - 测试文件类型白名单
   - 测试边缘情况（特殊字符等）

4. 网络失败场景
   - 测试 GitHub 不可达
   - 测试下载超时
   - 测试磁盘空间不足

5. 并发场景
   - 同时运行多个 bingoctl create
   - 测试缓存目录冲突

6. 边界情况
   - 目标目录已存在且非空
   - 缓存目录权限不足
   - 无效的模块名（特殊字符）
   - 不存在的版本号

### 集成测试（待实现）

```bash
# 测试默认行为
bingoctl create test-default
cd test-default && go mod tidy && go build ./...

# 测试包名替换
bingoctl create test-custom -m github.com/test/custom
cd test-custom && go mod tidy && go build ./...

# 测试版本指定
bingoctl create test-v120 -r v1.2.0

# 测试分支和缓存更新
bingoctl create test-main -r main
bingoctl create test-main-nocache -r main --no-cache

# 测试服务过滤
bingoctl create test-minimal --services apiserver
bingoctl create test-none --services none

# 测试镜像
export BINGOCTL_TEMPLATE_MIRROR=https://ghproxy.com/
bingoctl create test-mirror
```

## 实施步骤

1. ✅ 完成设计文档
2. [ ] 实现 `pkg/template/version.go`（版本解析）
3. [ ] 实现 `pkg/template/fetcher.go`（下载和缓存）
4. [ ] 实现 `pkg/template/replacer.go`（替换逻辑）
5. [ ] 修改 `pkg/cmd/create/create.go`（集成新模块）
6. [ ] 删除 `pkg/cmd/create/tpl/` 目录
7. [ ] 删除 `pkg/generator/template.go` 中的 embed（如果不再需要）
8. [ ] 编写单元测试
9. [ ] 编写集成测试
10. [ ] 文档更新（README, 用户指南）

## 风险与缓解

| 风险 | 影响 | 缓解措施 |
|-----|------|---------|
| 网络不可用 | 首次使用失败 | 支持镜像配置，提供友好错误提示 |
| 下载超时 | 用户体验差 | 30秒超时 + 镜像支持 + 显示进度条 |
| bingo 项目结构变化 | 替换逻辑失效 | 使用 .bingoctl.yaml 元数据文件（未来） |
| 包名替换不完整 | 生成的项目无法编译 | 白名单文件类型 + 充分测试 |
| 缓存目录冲突 | 多实例同时运行失败 | 使用文件锁保证并发安全 |
| Tarball 格式变化 | 解压失败 | 检测并处理根目录名 |
| 默认版本过时 | 用户使用旧版本 | 每次 bingoctl 发布时更新推荐版本 |

## 附录

### Tarball URL 格式

```
# Tag
https://github.com/bingo-project/bingo/archive/refs/tags/v1.2.3.tar.gz

# Branch
https://github.com/bingo-project/bingo/archive/refs/heads/main.tar.gz

# Commit（不推荐，优先使用 tag 或 branch）
https://github.com/bingo-project/bingo/archive/{commit-hash}.tar.gz
```

### 镜像配置示例

```bash
# 使用 ghproxy 镜像
export BINGOCTL_TEMPLATE_MIRROR=https://ghproxy.com/

# 使用其他镜像
export BINGOCTL_TEMPLATE_MIRROR=https://mirror.ghproxy.com/

# 清除镜像配置
unset BINGOCTL_TEMPLATE_MIRROR
```

### 缓存管理（未来功能）

未来可考虑添加缓存管理命令：

```bash
bingoctl cache list              # 列出缓存的模板
bingoctl cache clean             # 清理所有缓存
bingoctl cache clean v1.2.0      # 清理指定版本
bingoctl cache info              # 显示缓存目录大小和位置
```

### 依赖库

实现需要的 Go 依赖：

```go
// 文件锁（并发安全）
github.com/gofrs/flock

// 进度条（可选，提升用户体验）
github.com/schollz/progressbar/v3

// YAML 解析（用于 .bingoctl.yaml，未来）
gopkg.in/yaml.v3
```

### 相关文档

- [bingo 项目](https://github.com/bingo-project/bingo)
- [GitHub Archive 下载说明](https://docs.github.com/en/repositories/working-with-files/using-files/downloading-source-code-archives)
- [Go embed 文档](https://pkg.go.dev/embed)（旧方案参考）
