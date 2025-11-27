# BingoCTL

BingoCTL 是一个 Go 语言的脚手架和代码生成工具，用于快速创建和开发基于 Bingo 框架的应用程序。

## 功能特性

- 🚀 快速创建项目脚手架
- 📝 代码生成器，支持生成各层代码
- 🔄 数据库迁移管理
- 🗄️ 从数据库表自动生成模型代码
- 🛠️ 灵活的配置系统
- 🎯 支持 HTTP 和 gRPC 服务

## 安装

```bash
go install github.com/bingo-project/bingoctl@latest
```

## 配置文件

在项目根目录下创建配置文件 `.bingoctl.yaml`：

```yaml
version: v1

rootPackage: github.com/your-org/your-project

directory:
  cmd: internal/bingoctl/cmd
  model: internal/pkg/model
  store: internal/apiserver/store
  request: pkg/api/v1
  biz: internal/apiserver/biz
  controller: internal/apiserver/controller/v1
  middleware: internal/pkg/middleware
  job: internal/watcher/watcher
  migration: internal/apiserver/database/migration
  seeder: internal/apiserver/database/seeder

registries:
  router: internal/apiserver/router/api.go
  store:
    filePath: internal/apiserver/store/store.go
    interface: "IStore"
  biz:
    filePath: internal/apiserver/biz/biz.go
    interface: "IBiz"

mysql:
  host: 127.0.0.1:3306
  username: root
  password:
  database: bingo
```

## 命令使用

### 全局选项

```bash
-c, --config string   配置文件路径（默认使用 .bingoctl.yaml）
```

### create - 创建项目

从零创建一个新的项目脚手架。

```bash
bingoctl create <package-name>

# 示例
bingoctl create github.com/myorg/myapp
```

### make - 代码生成

生成各种类型的代码文件。

#### 全局选项

```bash
-d, --directory string   指定生成文件的目录
-p, --package string     指定包名
-t, --table string       从数据库表读取字段
-s, --service string     目标服务名称，用于自动推断路径
```

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

#### crud - 生成完整 CRUD 代码

一次性生成 model、store、biz、controller、request 的完整代码。

```bash
bingoctl make crud <name>

# 示例
bingoctl make crud user
```

#### model - 生成模型代码

```bash
bingoctl make model <name> [-d dir] [-p package] [-t table]

# 示例
bingoctl make model user
bingoctl make model user -t users  # 从 users 表生成
```

#### store - 生成存储层代码

```bash
bingoctl make store <name> [-d dir] [-p package]

# 示例
bingoctl make store user
```

#### biz - 生成业务逻辑层代码

```bash
bingoctl make biz <name> [-d dir] [-p package]

# 示例
bingoctl make biz user
```

#### controller - 生成控制器代码

```bash
bingoctl make controller <name> [-d dir] [-p package]

# 示例
bingoctl make controller user
```

#### request - 生成请求验证代码

```bash
bingoctl make request <name> [-d dir] [-p package]

# 示例
bingoctl make request user
```

#### middleware - 生成中间件代码

```bash
bingoctl make middleware <name> [-d dir] [-p package]

# 示例
bingoctl make middleware auth
```

#### cmd - 生成命令行代码

```bash
bingoctl make cmd <name> [-d dir] [-p package]

# 示例
bingoctl make cmd serve
```

#### job - 生成定时任务代码

```bash
bingoctl make job <name> [-d dir] [-p package]

# 示例
bingoctl make job cleanup
```

#### migration - 生成数据库迁移文件

```bash
bingoctl make migration <name> [-d dir] [-p package] [-t table]

# 示例
bingoctl make migration create_users_table
bingoctl make migration create_posts_table -t posts
```

#### seeder - 生成数据填充文件

```bash
bingoctl make seeder <name> [-d dir] [-p package]

# 示例
bingoctl make seeder users
```

#### service - 生成服务模块

生成一个完整的服务模块，包括 HTTP/gRPC 服务器配置。

```bash
bingoctl make service <name> [选项]

# 选项
--http                  启用 HTTP 服务器
--grpc                  启用 gRPC 服务器
--with-biz              生成业务层（默认 true）
--no-biz                不生成业务层（覆盖 --with-biz）
--with-store            生成存储层
--with-controller       生成控制器层
--with-middleware       生成中间件目录
--with-router           生成路由目录

# 示例
bingoctl make service api --http
bingoctl make service gateway --http --grpc --with-store --with-controller
bingoctl make service worker --no-biz
```

### gen - 从数据库生成代码

从数据库表自动生成 model 代码。

```bash
bingoctl gen -t <table1,table2,...>

# 示例
bingoctl gen -t users
bingoctl gen -t users,posts,comments
```

### version - 查看版本

```bash
bingoctl version
```

## 使用示例

### 1. 创建新项目

```bash
# 创建项目
bingoctl create github.com/myorg/blog

# 进入项目目录
cd blog

# 生成用户模块的完整 CRUD 代码
bingoctl make crud user

# 生成文章模块的完整 CRUD 代码
bingoctl make crud post
```

### 2. 从数据库生成模型

```bash
# 从现有数据库表生成模型
bingoctl gen -t users,posts,comments
```

### 3. 生成新服务

```bash
# 生成一个带 HTTP 服务器的 API 服务
bingoctl make service api --http --with-store --with-controller

# 生成一个纯业务处理的 worker 服务
bingoctl make service worker --no-biz
```

### 4. 生成迁移和数据填充

```bash
# 生成数据库迁移文件
bingoctl make migration create_users_table

# 生成数据填充文件
bingoctl make seeder users
```

## 目录结构

使用 bingoctl 创建的项目典型目录结构：

```
myapp/
├── cmd/                          # 命令行入口
│   └── myapp/
├── internal/
│   ├── apiserver/
│   │   ├── biz/                 # 业务逻辑层
│   │   ├── controller/          # 控制器
│   │   ├── database/
│   │   │   ├── migration/       # 数据库迁移
│   │   │   └── seeder/          # 数据填充
│   │   ├── model/               # 数据模型
│   │   ├── router/              # 路由
│   │   └── store/               # 存储层
│   ├── pkg/
│   │   └── middleware/          # 中间件
│   └── watcher/
│       └── watcher/             # 定时任务
├── pkg/
│   └── api/
│       └── v1/                  # API 请求/响应定义
├── .bingoctl.yaml               # bingoctl 配置文件
└── go.mod
```

## 开发工作流

1. **初始化项目**：使用 `bingoctl create` 创建新项目
2. **配置数据库**：在 `.bingoctl.yaml` 中配置数据库连接
3. **生成代码**：
   - 使用 `bingoctl make crud` 快速生成 CRUD 代码
   - 使用 `bingoctl gen` 从数据库生成模型
4. **数据库管理**：
   - 使用 `bingoctl make migration` 创建迁移文件
   - 使用 `bingoctl make seeder` 创建数据填充文件
5. **扩展功能**：根据需要使用 `make` 命令生成其他组件

## 许可证

[许可证信息]
