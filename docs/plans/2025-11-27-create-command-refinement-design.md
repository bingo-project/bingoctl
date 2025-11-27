# Create Command Refinement Design

## 背景

当前 `bingoctl create` 命令会创建包含 5 个服务的完整项目模板：apiserver、bot、ctl、admserver、scheduler。由于已经有了 `make service` 命令可以按需生成服务，create 命令的模板显得过于臃肿。

## 目标

精简 create 命令的项目模板，支持用户灵活选择需要创建的服务，同时保持良好的用户体验。

## 设计方案

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

## 实现方案

### 代码结构调整

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

## 兼容性

- **向后兼容**：不提供任何参数时，行为变为交互式选择（默认 apiserver + ctl）
- **文档更新**：需要更新 README.md 中的 create 命令说明

## 测试计划

1. 测试交互式模式的多选功能
2. 测试各种参数组合
3. 测试不选择任何服务的警告和确认流程
4. 测试生成的项目结构是否正确
5. 测试生成的 `.bingoctl.yaml` 配置是否正确
6. 测试与 `make service` 的配合使用

## 后续工作

1. 实现本设计方案
2. 更新文档和示例
3. 添加单元测试和集成测试
