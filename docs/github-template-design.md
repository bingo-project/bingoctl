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
2. 默认使用最新稳定版本
3. 本地缓存机制，避免重复下载
4. 支持可选的包名替换（从 `bingo` 改为用户指定的模块名）
5. 简化模板逻辑，不使用 `.tpl` 文件

## 设计决策

### 1. 命令行参数

```bash
# 使用默认最新稳定版，不修改包名
bingoctl create demo

# 指定版本
bingoctl create demo -r v1.2.3
bingoctl create demo -r main

# 指定包名（自动替换）
bingoctl create demo -m github.com/mycompany/demo

# 组合使用
bingoctl create demo -m github.com/mycompany/demo -r v1.2.3

# 服务选择（保留现有功能）
bingoctl create demo --services apiserver,ctl
bingoctl create demo --no-service bot,scheduler
```

**参数说明：**
- `NAME`：项目目录名（必需）
- `-m, --module`：Go 模块名（可选），指定后会替换所有 `bingo` 包名引用
- `-r, --ref`：模板版本（可选），支持 tag/branch/commit，默认为最新 stable tag

### 2. 版本管理

- **最新稳定版**：通过 GitHub API 获取所有 tags，过滤符合 `v*.*.*` 模式的 semver tag，排序后取最新
- **分支支持**：支持 `main` 等分支，缓存时使用 `main-{commit-hash}` 避免缓存过期
- **固定版本**：支持具体的 tag 或 commit hash

### 3. 缓存策略

```
~/.bingoctl/templates/
  ├── v1.2.3/              # tag 版本缓存
  ├── v1.2.4/
  └── main-abc123def/      # 分支缓存（分支名-commit hash）
```

- Tag 版本永久缓存
- 分支版本带 commit hash，避免内容变化导致的问题
- 如果缓存存在，直接使用，不重新下载

### 4. 下载方式

使用 GitHub Tarball API：
```
GET https://github.com/bingo-project/bingo/archive/refs/tags/{ref}.tar.gz
```

**优点：**
- 快速，体积小（不包含 .git）
- 不依赖系统 git 命令
- Go 标准库可处理 tar.gz

**超时设置：** 30 秒

### 5. 模板处理流程

**不使用 `.tpl` 模板文件**，而是直接处理真实项目：

```
1. 下载并缓存 bingo 项目
   ↓
2. 复制到用户指定目录
   ↓
3. 过滤服务（根据服务选择参数）
   - 删除未选中的服务的 cmd/ 和 internal/ 目录
   - 保留 internal/pkg/（共享代码）
   - 例如：只选择 apiserver 时删除 cmd/bingo-admserver, internal/admserver 等
   ↓
4. 重命名目录（总是执行）
   - cmd/bingo-apiserver → cmd/{app}-apiserver
   - cmd/bingoctl → cmd/{app}ctl
   - 其他保留的 bingo-* 目录
   ↓
5. 替换包名（仅当指定 -m 时）
   - go.mod: module bingo → module {newModule}
   - *.go: import "bingo/xxx" → import "{newModule}/xxx"
   - *.go: 字符串中的 "bingo/" → "{newModule}/"
```

### 6. 错误处理

- **网络不可用**：直接报错，提示检查网络连接
- **版本不存在**：提示用户访问 GitHub releases 页面查看可用版本
- **下载超时**：明确提示超时（30秒），建议检查网络

## 技术实现

### 核心模块

#### 1. `pkg/template/fetcher.go`

模板获取器，负责下载和缓存。

```go
type Fetcher struct {
    cacheDir string        // ~/.bingoctl/templates
    timeout  time.Duration // 30s
    repoURL  string        // https://github.com/bingo-project/bingo
}

// FetchTemplate 下载模板到缓存（如不存在），返回缓存路径
func (f *Fetcher) FetchTemplate(ref string) (string, error)

// downloadWithTimeout 下载 tarball，30秒超时
func (f *Fetcher) downloadWithTimeout(url string) (string, error)

// extractTarball 解压 tarball 到缓存目录
func (f *Fetcher) extractTarball(tarPath, destDir string) error

// getCacheKey 生成缓存目录名
// - tag: v1.2.3 → v1.2.3
// - branch: main → main-{commit-hash}
func (f *Fetcher) getCacheKey(ref string) (string, error)
```

#### 2. `pkg/template/version.go`

版本解析器，从 GitHub 获取版本信息。

```go
const githubAPIURL = "https://api.github.com/repos/bingo-project/bingo"

// GetLatestStableTag 获取最新稳定 tag
func GetLatestStableTag() (string, error)

// GetCommitHash 获取 branch/tag 的 commit hash
func GetCommitHash(ref string) (string, error)

// filterSemverTags 过滤符合 v*.*.* 的 tags
func filterSemverTags(tags []Tag) []Tag

// sortAndGetLatest 排序并返回最新版本
func sortAndGetLatest(tags []Tag) string
```

**实现细节：**
- 调用 `GET /repos/bingo-project/bingo/tags` 获取所有 tags
- 使用正则 `^v\d+\.\d+\.\d+$` 过滤 semver tags
- 排序后返回最新版本

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
func (r *Replacer) ReplaceModuleName() error

// replaceGoMod 替换 go.mod 中的 module 声明
func (r *Replacer) replaceGoMod() error

// replaceGoFile 替换 .go 文件中的 import 和字符串
func (r *Replacer) replaceGoFile(path string) error

// RenameDirs 重命名目录
func (r *Replacer) RenameDirs() error
```

**替换规则：**

| 原始内容 | 替换后 | 文件类型 |
|---------|--------|---------|
| `module bingo` | `module {newModule}` | go.mod |
| `import "bingo/xxx"` | `import "{newModule}/xxx"` | *.go |
| `"bingo/internal/xxx"` | `"{newModule}/internal/xxx"` | *.go (字符串) |
| `cmd/bingo-apiserver/` | `cmd/{app}-apiserver/` | 目录名 |
| `cmd/bingoctl/` | `cmd/{app}ctl/` | 目录名 |

#### 4. 修改 `pkg/cmd/create/create.go`

移除 `embed.FS`，使用新的 template 模块。

```go
type CreateOptions struct {
    AppName      string   // 目录名
    ModuleName   string   // Go 模块名（可选）
    TemplateRef  string   // 模板版本
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
        "Template version (tag/branch/commit, default: latest stable)")

    // 保留服务选择 flags
    cmd.Flags().StringSliceVar(&o.Services, "services", nil, "...")
    cmd.Flags().StringSliceVar(&o.NoServices, "no-service", nil, "...")
    cmd.Flags().StringSliceVar(&o.AddServices, "add-service", nil, "...")

    return cmd
}

func (o *CreateOptions) Complete(cmd *cobra.Command, args []string) error {
    // 1. 解析模板版本
    if o.TemplateRef == "" {
        ver, err := template.GetLatestStableTag()
        o.TemplateRef = ver
        console.Info(fmt.Sprintf("使用最新稳定版本: %s", ver))
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
    templatePath, err := fetcher.FetchTemplate(o.TemplateRef)

    // 2. 复制到目标目录
    copyDir(templatePath, o.AppName)

    // 3. 过滤服务（在重命名之前删除未选中的服务目录）
    if len(o.selectedServices) < len(availableServices) {
        console.Info("过滤服务...")
        o.filterServices()
    }

    // 4. 重命名目录（总是执行）
    replacer := template.NewReplacer(o.AppName, "bingo", o.ModuleName, o.AppName)
    replacer.RenameDirs()

    // 5. 替换模块名（仅当指定 -m 时）
    if o.ModuleName != "" {
        replacer.ReplaceModuleName()
    }

    return nil
}

// filterServices 删除未选中的服务目录
func (o *CreateOptions) filterServices() error {
    // 构建服务映射
    allServices := map[string]struct{
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
            cmdPath := filepath.Join(o.AppName, dirs.cmdDir)
            if exists(cmdPath) {
                console.Info(fmt.Sprintf("  删除 %s", dirs.cmdDir))
                os.RemoveAll(cmdPath)
            }

            // 删除 internal 目录
            internalPath := filepath.Join(o.AppName, dirs.internalDir)
            if exists(internalPath) {
                console.Info(fmt.Sprintf("  删除 %s", dirs.internalDir))
                os.RemoveAll(internalPath)
            }
        }
    }

    return nil
}
```

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

### 单元测试

1. `template/version.go`
   - 测试 GitHub API 调用
   - 测试 semver tag 过滤
   - 测试版本排序

2. `template/fetcher.go`
   - 测试缓存命中逻辑
   - 测试下载和解压
   - 测试超时处理

3. `template/replacer.go`
   - 测试模块名替换
   - 测试目录重命名
   - 测试边缘情况（特殊字符等）

### 集成测试

```bash
# 测试默认行为
bingoctl create test-default
cd test-default && go mod tidy && go build ./...

# 测试包名替换
bingoctl create test-custom -m github.com/test/custom
cd test-custom && go mod tidy && go build ./...

# 测试版本指定
bingoctl create test-v120 -r v1.2.0

# 测试服务过滤
bingoctl create test-minimal --services apiserver
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
| GitHub API 限流 | 无法获取最新版本 | 缓存最新版本号，降级到缓存 |
| 网络不可用 | 首次使用失败 | 提供离线安装包选项 |
| bingo 项目结构变化 | 替换逻辑失效 | 版本化替换规则，针对不同版本使用不同策略 |
| 包名替换不完整 | 生成的项目无法编译 | 充分测试，收集边缘 case |

## 附录

### GitHub API 参考

- 获取 tags: `GET /repos/bingo-project/bingo/tags`
- 获取 commit: `GET /repos/bingo-project/bingo/commits/{ref}`
- 下载 tarball: `GET /repos/bingo-project/bingo/tarball/{ref}`（或使用 archive URL）

### Tarball URL 格式

```
# Tag
https://github.com/bingo-project/bingo/archive/refs/tags/v1.2.3.tar.gz

# Branch
https://github.com/bingo-project/bingo/archive/refs/heads/main.tar.gz

# Commit
https://github.com/bingo-project/bingo/archive/{commit-hash}.tar.gz
```

### 缓存目录清理

未来可考虑添加缓存管理命令：

```bash
bingoctl cache list              # 列出缓存的模板
bingoctl cache clean             # 清理所有缓存
bingoctl cache clean v1.2.0      # 清理指定版本
```
