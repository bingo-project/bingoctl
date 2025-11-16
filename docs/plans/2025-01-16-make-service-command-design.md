# Make Service 命令设计文档

**日期：** 2025-01-16
**目标：** 为 bingoctl 添加 `make service` 子命令，支持生成与 cmd/{app}-apiserver 平级的服务模块

## 概述

`make service` 命令用于在现有项目中生成新的微服务模块，包括：
- `cmd/<name>/main.go` - 服务入口
- `internal/<name>/` - 服务实现
- `configs/<name>.yaml` - 服务配置

服务支持可选的 HTTP/gRPC 服务器，以及可配置的业务层目录结构。

## 命令接口

### 基本用法

```bash
bingoctl make service <name> [flags]
```

### 标志列表

- `--http` - 生成 HTTP 服务器相关代码
- `--grpc` - 生成 gRPC 服务器相关代码
- `--with-biz` - 生成 biz/ 业务逻辑层目录 (默认启用)
- `--no-biz` - 不生成 biz/ 业务逻辑层目录 (覆盖 --with-biz)
- `--with-store` - 生成 store/ 数据访问层目录
- `--with-controller` - 生成 controller/ 控制器目录
- `--with-middleware` - 生成 middleware/ 中间件目录
- `--with-router` - 生成 router/ 路由目录

### 使用示例

```bash
# 最小化服务（只有基础框架，默认包含 biz 目录）
bingoctl make service payment

# 不包含 biz 目录的最小化服务
bingoctl make service payment --no-biz

# HTTP API 服务
bingoctl make service payment --http --with-controller --with-router

# 完整的 HTTP+gRPC 服务
bingoctl make service payment --http --grpc --with-store --with-controller --with-router --with-middleware
```

## 生成的目录结构

### 最小化结构

执行 `bingoctl make service payment` 生成：

```
cmd/payment/
  └── main.go              # 服务入口

internal/payment/
  ├── app.go              # Cobra 命令定义，应用初始化
  ├── biz/                # 业务逻辑层目录（默认创建）
  │   └── biz.go          # 基本的业务逻辑接口和实现
  └── run.go              # 服务启动逻辑（空框架）

configs/
  └── payment.yaml        # 服务配置文件
```

### 添加 HTTP 服务器

使用 `--http` 标志：

```
internal/payment/
  ├── app.go
  ├── run.go              # 包含 HTTP 服务器启动代码
  └── server.go           # HTTP 服务器配置和初始化
```

### 添加 gRPC 服务器

使用 `--grpc` 标志：

```
internal/payment/
  ├── app.go
  ├── run.go              # 包含 gRPC 服务器启动代码
  ├── grpc.go             # gRPC 服务器配置和初始化
  └── grpc/               # gRPC 服务实现目录
      └── .gitkeep
```

### 完整结构示例

使用所有标志：

```
cmd/payment/
  └── main.go

internal/payment/
  ├── app.go
  ├── run.go
  ├── server.go
  ├── grpc.go
  ├── biz/
  │   └── .gitkeep
  ├── store/
  │   └── .gitkeep
  ├── controller/
  │   └── .gitkeep
  ├── router/
  │   ├── http.go        # HTTP 路由注册
  │   └── grpc.go        # gRPC 服务注册
  ├── middleware/
  │   └── .gitkeep
  └── grpc/
      └── .gitkeep

configs/
  └── payment.yaml
```

## 实现设计

### 代码位置

- 新建文件：`pkg/cmd/make/make_service.go`
- 修改文件：`pkg/cmd/make/make.go` - 注册新子命令

### 模板文件组织

```
pkg/cmd/make/tpl/service/
  ├── cmd_main.go.tpl           # cmd/<name>/main.go 模板
  ├── app.go.tpl                # internal/<name>/app.go 模板
  ├── run_minimal.go.tpl        # 最小化的 run.go（无服务器）
  ├── run_http.go.tpl           # 带 HTTP 的 run.go
  ├── run_grpc.go.tpl           # 带 gRPC 的 run.go
  ├── run_both.go.tpl           # 带 HTTP+gRPC 的 run.go
  ├── server.go.tpl             # HTTP 服务器模板
  ├── grpc.go.tpl               # gRPC 服务器模板
  ├── router_http.go.tpl        # router/http.go 模板
  ├── router_grpc.go.tpl        # router/grpc.go 模板
  └── config.yaml.tpl           # configs/<name>.yaml 模板
```

### 数据结构

```go
type ServiceOptions struct {
    *generator.Options
    EnableHTTP     bool
    EnableGRPC     bool
    WithBiz        bool
    WithStore      bool
    WithController bool
    WithMiddleware bool
    WithRouter     bool
}

type ServiceTemplateData struct {
    ServiceName   string  // 服务名称，如 "payment"
    RootPackage   string  // 根包路径，从当前项目 go.mod 读取
    EnableHTTP    bool    // 是否启用 HTTP
    EnableGRPC    bool    // 是否启用 gRPC
    WithBiz       bool    // 是否生成 biz 层
    WithStore     bool    // 是否生成 store 层
    WithController bool   // 是否生成 controller
    WithMiddleware bool   // 是否生成 middleware
    WithRouter    bool    // 是否生成 router
}
```

### 核心逻辑

1. **验证阶段**：
   - 检查服务名称是否提供
   - 检查 `cmd/` 和 `internal/` 目录是否存在
   - 检查服务名是否已存在（避免覆盖）

2. **准备阶段**：
   - 读取 `go.mod` 获取根包路径
   - 根据标志组合确定要生成的文件和目录

3. **生成阶段**：
   - 生成 `cmd/<name>/main.go`
   - 生成 `internal/<name>/app.go`
   - 根据 HTTP/gRPC 组合选择合适的 `run.go` 模板
   - 如果 `--http`，生成 `server.go`
   - 如果 `--grpc`，生成 `grpc.go` 和 `grpc/` 目录
   - 根据 `--with-*` 标志创建对应目录（使用 `.gitkeep` 占位）
   - 如果 `--with-router`，根据启用的服务生成路由文件
   - 生成 `configs/<name>.yaml`

### run.go 模板选择逻辑

| EnableHTTP | EnableGRPC | 使用的模板 |
|-----------|-----------|---------|
| false | false | run_minimal.go.tpl |
| true | false | run_http.go.tpl |
| false | true | run_grpc.go.tpl |
| true | true | run_both.go.tpl |

## 核心模板示例

### cmd/<name>/main.go.tpl

```go
package main

import (
	"github.com/spf13/cobra"

	"{{.RootPackage}}/internal/{{.ServiceName}}"
)

func main() {
	command := {{.ServiceName}}.NewAppCommand()
	cobra.CheckErr(command.Execute())
}
```

### internal/<name>/app.go.tpl

```go
package {{.ServiceName}}

import (
	"github.com/spf13/cobra"
	"github.com/bingo-project/component-base/cli"
)

func NewAppCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "{{.ServiceName}}",
		Short: "{{.ServiceName}} service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}

	cli.AddConfigFlag(cmd, "{{.ServiceName}}")

	return cmd
}
```

### configs/<name>.yaml.tpl

```yaml
server:
{{- if .EnableHTTP}}
  http:
    addr: :8080
    mode: release
{{- end}}
{{- if .EnableGRPC}}
  grpc:
    addr: :9090
{{- end}}

log:
  level: info
  format: console
  output-paths:
    - stdout
```

## 设计决策

1. **默认最小化** - 遵循 YAGNI 原则，默认只生成基础框架，其他功能通过标志按需添加
2. **独立配置** - 每个服务有独立的配置文件，符合微服务架构原则
3. **灵活组合** - 通过独立标志支持任意功能组合，而不是预设模板
4. **自定义命名** - 服务名完全由用户指定，不强制加前缀或后缀
5. **空目录占位** - 使用 `.gitkeep` 占位空目录，确保目录结构被 git 追踪

## 后续扩展

- 支持从现有服务复制结构（`--from` 标志）
- 支持生成 Dockerfile 和 kubernetes 部署文件
- 支持自定义模板路径

---

## 实现状态

**实现完成日期:** 2025-01-16

**已实现功能:**
- ✅ `bingoctl make service` 基础命令
- ✅ `--http` 标志 - 生成 HTTP 服务器
- ✅ `--grpc` 标志 - 生成 gRPC 服务器
- ✅ `--with-biz` 标志 - 生成 biz 目录 (默认启用)
- ✅ `--no-biz` 标志 - 不生成 biz 目录 (覆盖 --with-biz)
- ✅ `--with-store` 标志 - 生成 store 目录
- ✅ `--with-controller` 标志 - 生成 controller 目录
- ✅ `--with-middleware` 标志 - 生成 middleware 目录
- ✅ `--with-router` 标志 - 生成 router 目录
- ✅ 配置文件生成
- ✅ 模板系统
- ✅ 手动测试验证

**测试结果:** 所有功能测试通过

**使用示例:**
```bash
# 最小化服务（默认包含 biz 目录）
bingoctl make service payment

# 不包含 biz 目录的最小化服务
bingoctl make service payment --no-biz

# HTTP API 服务
bingoctl make service order --http --with-router

# 完整服务
bingoctl make service inventory --http --grpc --with-store --with-controller --with-router
```
