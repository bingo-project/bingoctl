# Create Command Refinement Design

## 背景

当前 `bingoctl create` 命令会创建包含 5 个服务的完整项目模板：apiserver、bot、ctl、admserver、scheduler。由于已经有了 `make service` 命令可以按需生成服务，create 命令的模板显得过于臃肿。

另外，当前 `make` 子命令（如 `make model`、`make biz` 等）使用 `.bingoctl.yaml` 中配置的固定路径，当项目中有多个服务时，需要为不同服务生成代码比较麻烦，需要每次手动指定 `-d` 参数。

## 目标

1. 精简 create 命令的项目模板，支持用户灵活选择需要创建的服务，同时保持良好的用户体验
2. 增强 make 命令，支持通过 `--service` 参数指定目标服务，自动推断代码生成路径

## 设计方案

## Part 1: Create 命令精简

### 1. 命令行接口

#### 交互式模式（默认）

```bash
bingoctl create github.com/myorg/myapp
```

运行后弹出多选界面：
```
? 选择要创建的服务: (空格选择，回车确认)
✓ apiserver (默认选中)
✓ ctl (默认选中)
✗ admserver
✗ bot
✗ scheduler
```

#### 命令行参数模式（跳过交互）

**方式1：明确指定服务（覆盖默认）**
```bash
bingoctl create myapp --services apiserver,ctl,bot
bingoctl create myapp --services apiserver  # 只创建 apiserver
bingoctl create myapp --services none       # 不创建任何服务
```

**方式2：基于默认调整**
```bash
bingoctl create myapp --no-service ctl                    # 只创建 apiserver
bingoctl create myapp --no-service apiserver,ctl          # 不创建任何服务
bingoctl create myapp --add-service bot,scheduler         # apiserver + ctl + bot + scheduler
```

**方式3：混合使用**
```bash
bingoctl create myapp --services apiserver --add-service bot  # apiserver + bot
```

**参数优先级：**
- `--services` 优先级最高，会覆盖 `--no-service` 和 `--add-service`
- 如果只使用 `--no-service` 或 `--add-service`，则基于默认选项（apiserver + ctl）进行调整

### 2. 项目结构

#### 最小骨架（不选择任何服务时）

```
myapp/
├── .bingoctl.yaml          # bingoctl 配置文件
├── .gitignore
├── .air.example.toml       # 热加载配置示例
├── .golangci.yaml          # 代码检查配置
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── LICENSE
├── cmd/                    # 空目录，预留给服务入口
├── configs/                # 空目录，预留给配置文件
├── internal/               # 空目录，预留给内部代码
├── pkg/                    # 空目录，预留给公共包
├── scripts/                # 构建和部署脚本
├── deployments/            # 部署配置（k8s/docker等）
└── storage/                # 本地存储目录
```

当不选择任何服务时：
1. 显示警告信息：`Warning: 未选择任何服务，将创建最小项目骨架`
2. 要求用户确认：`继续? (Y/n)`
3. 确认后创建最小骨架

#### 选择服务后的变化

- **apiserver**：添加 `cmd/{app}-apiserver/`、`internal/apiserver/`、配置文件、Swagger 文档等
- **ctl**：添加 `cmd/{app}ctl/`、`internal/{app}ctl/cmd/`（如果有）
- **admserver**：添加 `cmd/{app}-admserver/`、相关配置等
- **bot**：添加 `cmd/{app}-bot/`、相关配置等
- **scheduler**：添加 `cmd/{app}-scheduler/`、相关配置等

### 3. 服务模板内容

#### apiserver 服务包含

```
cmd/{app}-apiserver/main.go
internal/apiserver/
  ├── biz/
  │   ├── auth/          # 认证相关业务逻辑（示例）
  │   ├── user/          # 用户业务逻辑（示例）
  │   ├── file/          # 文件业务逻辑（示例）
  │   ├── syscfg/        # 系统配置业务逻辑（示例）
  │   ├── app/           # 应用业务逻辑（示例）
  │   ├── common/        # 公共业务逻辑
  │   └── biz.go         # biz 注册
  ├── controller/v1/     # 控制器层（示例）
  ├── store/             # 存储层（示例）
  ├── model/             # 数据模型（示例）
  ├── router/            # 路由配置
  └── database/
      ├── migration/     # 数据库迁移
      └── seeder/        # 数据填充
configs/{app}-apiserver.yaml
api/swagger/apiserver/   # Swagger 文档
```

#### ctl 服务包含

```
cmd/{app}ctl/main.go
internal/{app}ctl/cmd/   # CLI 命令定义
```

#### 其他服务

admserver、bot、scheduler 类似，包含各自的 main.go、配置文件、业务逻辑示例等。

#### .bingoctl.yaml 配置文件

根据选择的服务，动态生成对应的 directory 配置。例如只选择 ctl 时，不包含 apiserver 相关的目录配置。

### 4. create vs make service 的区别

- **`create` 创建的服务**：包含完整示例代码（auth、user、file 等业务逻辑）
- **`make service` 创建的服务**：只有基础结构骨架，不包含示例代码

这样设计的原因：
- `create` 用于快速启动新项目，提供参考示例
- `make service` 用于在已有项目中添加新服务，避免污染现有代码

## Part 2: Make 命令支持服务选择

### 5. 问题场景

当项目中有多个服务（如 apiserver、admserver）时，配置文件中只能配置一个默认路径：

```yaml
# .bingoctl.yaml
directory:
  model: internal/apiserver/model
  store: internal/apiserver/store
  biz: internal/apiserver/biz
```

如果要为 admserver 生成代码，需要每次手动指定路径：
```bash
bingoctl make model user -d internal/admserver/model
```

这很不方便。

### 6. 解决方案：渐进式路径推断

在 `make` 子命令中添加 `--service` 参数，自动推断目标路径。

#### 命令行接口

```bash
# 使用配置默认路径（现有行为，保持不变）
bingoctl make model user

# 为指定服务生成代码（新增）
bingoctl make model user --service admserver

# 明确指定路径（优先级最高）
bingoctl make model user -d custom/path
```

#### 路径推断逻辑

采用三层推断策略：

**1. 智能替换（优先）**
- 扫描 `cmd/` 目录，识别已存在的服务名称
  - 例如：`cmd/myapp-apiserver` → 识别服务名 `apiserver`
  - 例如：`cmd/myapp-admserver` → 识别服务名 `admserver`
- 如果配置路径中包含已识别的服务名，则替换
  - `internal/apiserver/model` + `--service admserver` → `internal/admserver/model`
  - `internal/apiserver/biz/user` + `--service admserver` → `internal/admserver/biz/user`

**2. 固定模式回退**
- 如果智能替换失败（配置路径中没有可识别的服务名）
- 使用固定模式：提取配置路径的后缀部分，拼接到 `internal/{service}/` 后面
  - `internal/pkg/model` + `--service admserver` → `internal/admserver/model`
  - `pkg/model` + `--service admserver` → `internal/admserver/model`

**3. 参数优先级**
- `-d` 明确指定目录 > `--service` 推断 > 配置默认路径

#### 适用命令

所有基于目录配置的 make 子命令都支持 `--service` 参数：
- `make model --service <name>`
- `make store --service <name>`
- `make biz --service <name>`
- `make controller --service <name>`
- `make request --service <name>`
- `make middleware --service <name>`
- `make job --service <name>`
- `make migration --service <name>`
- `make seeder --service <name>`
- `make crud --service <name>`

注：`make cmd` 和 `make service` 不适用，因为它们是创建服务级别的代码。

### 7. 使用示例

假设项目结构：
```
myapp/
├── cmd/
│   ├── myapp-apiserver/
│   └── myapp-admserver/
├── internal/
│   ├── apiserver/
│   │   ├── model/
│   │   ├── store/
│   │   └── biz/
│   └── admserver/
│       ├── model/
│       ├── store/
│       └── biz/
```

配置文件：
```yaml
directory:
  model: internal/apiserver/model
  store: internal/apiserver/store
  biz: internal/apiserver/biz
```

使用示例：
```bash
# 为 apiserver 生成（使用默认配置）
bingoctl make model user
# → 生成到 internal/apiserver/model/user.go

# 为 admserver 生成（自动推断）
bingoctl make model user --service admserver
# → 扫描 cmd/ 发现 myapp-admserver，识别服务名 admserver
# → 替换路径：internal/apiserver/model → internal/admserver/model
# → 生成到 internal/admserver/model/user.go

# 完整 CRUD（自动推断所有层级）
bingoctl make crud order --service admserver
# → model: internal/admserver/model/order.go
# → store: internal/admserver/store/order.go
# → biz: internal/admserver/biz/order.go
# → controller: internal/admserver/controller/v1/order.go
# → request: pkg/api/admserver/v1/order.go
```

## 实现方案

### Part 1 实现：Create 命令

#### 代码结构调整

```go
// CreateOptions 新增字段
type CreateOptions struct {
    GoVersion    string
    TemplatePath string
    RootPackage  string
    AppName      string
    AppNameCamel string

    // 新增：服务选择相关
    Services      []string  // 明确指定的服务列表
    NoServices    []string  // 要排除的服务
    AddServices   []string  // 要添加的服务
    Interactive   bool      // 是否交互模式（默认 true）
}
```

### Flags 定义

```go
cmd.Flags().StringSliceVar(&o.Services, "services", nil, "明确指定要创建的服务（覆盖默认），使用 'none' 表示不创建任何服务")
cmd.Flags().StringSliceVar(&o.NoServices, "no-service", nil, "从默认中排除的服务")
cmd.Flags().StringSliceVar(&o.AddServices, "add-service", nil, "在默认基础上添加的服务")
```

### 执行流程

1. **解析参数**
   - 检查是否提供了 `--services`、`--no-service` 或 `--add-service`
   - 如果提供了任一参数，则 `Interactive = false`

2. **确定服务列表**
   - 如果 `Interactive = true`：使用 promptui 显示多选列表（apiserver 和 ctl 默认选中）
   - 如果 `Interactive = false`：根据参数计算最终服务列表
     - 优先使用 `--services`（如果提供）
     - 否则从默认列表（apiserver, ctl）出发，应用 `--no-service` 和 `--add-service`

3. **验证和确认**
   - 如果服务列表为空：显示警告并二次确认
   - 如果用户取消：退出程序

4. **生成项目**
   - 复制 common 目录下的通用文件
   - 根据服务列表，复制对应服务的模板文件
   - 生成 `.bingoctl.yaml`，只包含选中服务相关的配置

### 模板文件重组

将现有的 `pkg/cmd/create/tpl/` 目录重新组织：

```
tpl/
├── common/              # 所有项目共用
│   ├── .gitignore
│   ├── .air.example.toml
│   ├── .golangci.yaml
│   ├── go.mod.tpl
│   ├── go.sum.tpl
│   ├── Makefile
│   ├── README.md
│   ├── LICENSE
│   ├── scripts/
│   ├── deployments/
│   └── storage/
├── services/
│   ├── apiserver/       # apiserver 相关文件
│   │   ├── cmd/{app}-apiserver/
│   │   ├── internal/apiserver/
│   │   ├── configs/{app}-apiserver.yaml
│   │   └── api/swagger/apiserver/
│   ├── ctl/             # ctl 相关文件
│   │   ├── cmd/{app}ctl/
│   │   └── internal/{app}ctl/
│   ├── admserver/       # admserver 相关文件
│   ├── bot/             # bot 相关文件
│   └── scheduler/       # scheduler 相关文件
└── configs/
    └── .bingoctl.yaml.tpl  # 配置模板，根据服务动态生成
```

### 主要代码变更

1. **create.go**
   - 添加新的 flags
   - 实现交互式选择逻辑
   - 实现参数解析和服务列表计算
   - 修改文件复制逻辑，根据服务列表选择性复制

2. **模板文件**
   - 重新组织目录结构
   - 调整 `.bingoctl.yaml.tpl` 使其可以根据服务动态生成

### Part 2 实现：Make 命令服务选择

#### 代码结构调整

```go
// generator.Options 新增字段
type Options struct {
    // ... 现有字段

    // 新增：服务选择
    Service string  // 目标服务名称
}
```

#### Flags 定义

在 `pkg/cmd/make/make.go` 中添加全局 flag：
```go
cmd.PersistentFlags().StringVarP(&opt.Service, "service", "s", "", "Target service name")
```

#### 核心实现：路径推断函数

在 `pkg/generator/generate.go` 中添加新函数：

```go
// InferDirectoryForService 根据服务名推断目录路径
func (o *Options) InferDirectoryForService(baseDir, serviceName string) (string, error) {
    if serviceName == "" {
        return baseDir, nil
    }

    // 1. 扫描 cmd/ 目录，识别已存在的服务
    services, err := discoverServices()
    if err != nil {
        return "", err
    }

    // 2. 智能替换：如果路径中包含已知服务名，则替换
    for _, svc := range services {
        if strings.Contains(baseDir, svc) {
            return strings.ReplaceAll(baseDir, svc, serviceName), nil
        }
    }

    // 3. 固定模式回退：提取后缀，拼接到 internal/{service}/
    suffix := extractSuffix(baseDir)
    return filepath.Join("internal", serviceName, suffix), nil
}

// discoverServices 扫描 cmd/ 目录发现服务名称
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

        // 提取服务名：myapp-apiserver → apiserver, myappctl → ctl
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

// extractSuffix 提取路径后缀
func extractSuffix(path string) string {
    parts := strings.Split(filepath.Clean(path), string(filepath.Separator))

    // 找到 internal/ 或其他已知前缀之后的部分
    for i, part := range parts {
        if part == "internal" || part == "pkg" {
            if i+2 < len(parts) {
                return strings.Join(parts[i+2:], string(filepath.Separator))
            }
        }
    }

    // 如果找不到，返回最后一段
    if len(parts) > 0 {
        return parts[len(parts)-1]
    }

    return path
}
```

#### 修改 GetMapDirectory 函数

```go
func GetMapDirectory(tmpl string) (dir string) {
    // ... 现有的映射逻辑

    return
}

// 新增：根据服务名调整目录
func (o *Options) GetDirectoryForService(tmpl string) (string, error) {
    baseDir := GetMapDirectory(tmpl)
    return o.InferDirectoryForService(baseDir, o.Service)
}
```

#### 修改 GenerateCode 函数

```go
func (o *Options) GenerateCode(tmpl, path string) error {
    // 获取基础目录
    baseDir := GetMapDirectory(tmpl)

    // 如果指定了服务，推断目录
    var dir string
    var err error
    if o.Service != "" {
        dir, err = o.InferDirectoryForService(baseDir, o.Service)
        if err != nil {
            return fmt.Errorf("failed to infer directory for service %s: %w", o.Service, err)
        }
    } else {
        dir = baseDir
    }

    o.SetName(tmpl)
    o.ReadCodeTemplates()
    o.GenerateAttributes(dir, path)

    // ... 其余逻辑保持不变
}
```

#### 主要代码变更

1. **pkg/generator/option.go**
   - 添加 `Service` 字段

2. **pkg/generator/generate.go**
   - 实现 `InferDirectoryForService()` 函数
   - 实现 `discoverServices()` 辅助函数
   - 实现 `extractSuffix()` 辅助函数
   - 修改 `GenerateCode()` 使用新的推断逻辑

3. **pkg/cmd/make/make.go**
   - 添加 `--service` 全局 flag

4. **文档**
   - 更新 README.md 说明 `--service` 参数的使用

## 兼容性

### Part 1：Create 命令
- **向后兼容性影响**：不提供任何参数时，行为从"创建所有 5 个服务"变为"交互式选择（默认 apiserver + ctl）"
- **迁移建议**：
  - 如果脚本中使用了 `bingoctl create`，需要添加 `--services apiserver,ctl,admserver,bot,scheduler` 来保持原有行为
  - 或者使用交互式时默认全选所有服务
- **文档更新**：需要更新 README.md 中的 create 命令说明

### Part 2：Make 命令
- **完全向后兼容**：新增的 `--service` 参数是可选的，不影响现有用法
- **优先级保证**：`-d` 参数优先级最高，确保现有使用 `-d` 的脚本不受影响
- **文档更新**：需要更新 README.md 添加 `--service` 参数说明

## 测试计划

### Part 1：Create 命令测试
1. 测试交互式模式的多选功能
2. 测试各种参数组合（`--services`、`--no-service`、`--add-service`）
3. 测试不选择任何服务的警告和确认流程
4. 测试生成的项目结构是否正确
5. 测试生成的 `.bingoctl.yaml` 配置是否正确
6. 测试与 `make service` 的配合使用

### Part 2：Make 命令测试
1. 测试 `--service` 参数的路径推断逻辑
   - 智能替换场景
   - 固定模式回退场景
2. 测试参数优先级（`-d` > `--service` > 默认配置）
3. 测试不同服务名的推断（apiserver、admserver、ctl 等）
4. 测试在没有 `cmd/` 目录时的错误处理
5. 测试 `make crud --service` 的完整流程
6. 测试与现有 `-d` 参数的兼容性

## 后续工作

1. 实现 Part 1：Create 命令精简
   - 重组模板文件结构
   - 实现交互式选择
   - 实现参数解析逻辑
2. 实现 Part 2：Make 命令服务选择
   - 实现路径推断函数
   - 添加 `--service` flag
   - 修改代码生成逻辑
3. 更新文档和示例
4. 添加单元测试和集成测试
5. 更新 README.md
