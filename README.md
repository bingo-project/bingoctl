# Bingo CLI

[English](./README.md) | [中文](./README.zh.md)

Bingo CLI is a scaffolding and code generation tool written in Go, designed for quickly creating and developing applications based on the Bingo framework.

## Features

- 🚀 Rapid project scaffolding
- 📝 Code generators for all layers
- 🔄 Database migration management
- 🗄️ Auto-generate model code from database tables
- 🛠️ Flexible configuration system
- 🎯 Support for HTTP, gRPC, and WebSocket services

## Installation

```bash
go install github.com/bingo-project/bingoctl/cmd/bingo@latest
```

> See [CHANGELOG](./CHANGELOG.md) for version history

## Shell Completion

bingo supports command-line auto-completion for multiple shells.

### Zsh

```bash
# Temporary (current session)
source <(bingo completion zsh)

# Permanent
## Linux
bingo completion zsh > "${fpath[1]}/_bingo"

## macOS (Homebrew)
bingo completion zsh > $(brew --prefix)/share/zsh/site-functions/_bingo
```

> If completion doesn't work, ensure `.zshrc` has: `autoload -U compinit; compinit`

### Bash

```bash
# Temporary (current session)
source <(bingo completion bash)

# Permanent
## Linux
bingo completion bash > /etc/bash_completion.d/bingo

## macOS (Homebrew)
bingo completion bash > $(brew --prefix)/etc/bash_completion.d/bingo
```

> Requires the `bash-completion` package

### Fish

```bash
bingo completion fish > ~/.config/fish/completions/bingo.fish
```

### PowerShell

```powershell
bingo completion powershell > bingo.ps1
# Add the generated script to your PowerShell profile
```

## Configuration File

Create a configuration file `.bingo.yaml` in your project root:

```yaml
version: v1

rootPackage: github.com/your-org/your-project

directory:
  cmd: internal/bingoctl/cmd
  model: internal/pkg/model
  store: internal/apiserver/store
  biz: internal/apiserver/biz
  handler: internal/apiserver/handler/http
  middleware: internal/pkg/middleware/http
  request: pkg/api/apiserver/v1
  migration: internal/pkg/database/migration
  seeder: internal/pkg/database/seeder

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

## Commands

### Global Options

```bash
-c, --config string   Config file path (defaults to .bingo.yaml)
```

### create - Create Project

Create a new project scaffold from scratch. Downloads and caches Bingo project templates from GitHub.

```bash
bingo create <package-name> [options]

# Example
bingo create github.com/myorg/myapp
```

#### Create Command Options

**Template Version**

```bash
# Use recommended version (default)
bingo create myapp

# Use specific version
bingo create myapp -r v1.2.3

# Use branch (development version)
bingo create myapp -r main

# Force re-download branch template
bingo create myapp -r main --no-cache
```

**Custom Module Name**

```bash
# Replace package name
bingo create myapp -m github.com/mycompany/myapp
```

**Git Initialization**

```bash
# Create project and initialize git repo (default)
bingo create myapp

# Create project without git initialization
bingo create myapp --init-git=false
```

**Build Options**

```bash
# Create project without building (default)
bingo create myapp

# Create project and run make build
bingo create myapp --build
```

**Service Selection**

```bash
# Include only apiserver (default)
bingo create myapp

# Create all available services
bingo create myapp --all
# or
bingo create myapp -a

# Explicitly specify services
bingo create myapp --services apiserver,ctl,scheduler

# Add service to default apiserver
bingo create myapp --add-service admserver

# Exclude service
bingo create myapp --no-service bot

# Skeleton only, no services
bingo create myapp --services none
```

**Cache Management**

```bash
# Use cache (default) - speeds up creation
bingo create myapp

# Force refresh template (for branches)
bingo create myapp -r main --no-cache

# Cache location: ~/.bingo/templates/
```

**Mirror Configuration**

For regions with difficult GitHub access, configure a mirror:

```bash
# Using environment variable
export BINGO_TEMPLATE_MIRROR=https://ghproxy.com/
bingo create myapp

# Or temporary setting
BINGO_TEMPLATE_MIRROR=https://ghproxy.com/ bingo create myapp
```

### make - Code Generation

Generate various types of code files.

#### Global Options

```bash
-d, --directory string   Specify the directory for generated files
-p, --package string     Specify package name
-t, --table string       Read fields from database table
-s, --service string     Target service name for automatic path inference
```

#### Service Selection

When a project contains multiple services, use the `--service` parameter for automatic path inference. Path inference priority:

1. **Explicit directory** (`-d`) - Highest priority
2. **Service parameter** (`--service`) - Auto-infer path
3. **Config default path** - Usually apiserver path

```bash
# Generate code for default service (usually apiserver)
bingo make model user

# Auto-infer path for specific service
bingo make model user --service admserver

# Generate complete CRUD (for specified service)
bingo make crud order --service admserver

# Explicitly specify directory (highest priority)
bingo make model user -d custom/path
```

**Path Inference Rules:**
1. Scan `cmd/` directory to identify existing services
2. If config path contains service name, intelligently replace (e.g., `internal/apiserver/model` → `internal/admserver/model`)
3. Otherwise use default pattern: `internal/{service}/{suffix}`

#### crud - Generate Complete CRUD Code

Generate complete code for model, store, biz, handler, and request at once.

```bash
bingo make crud <name>

# Example
bingo make crud user
```

#### model - Generate Model Code

```bash
bingo make model <name> [-d dir] [-p package] [-t table]

# Examples
bingo make model user
bingo make model user -t users  # Generate from users table
```

#### store - Generate Store Layer Code

```bash
bingo make store <name> [-d dir] [-p package]

# Example
bingo make store user
```

#### biz - Generate Business Logic Layer Code

```bash
bingo make biz <name> [-d dir] [-p package]

# Example
bingo make biz user
```

#### handler - Generate Handler Code

```bash
bingo make handler <name> [-d dir] [-p package]

# Example
bingo make handler user
```

#### request - Generate Request Validation Code

```bash
bingo make request <name> [-d dir] [-p package]

# Example
bingo make request user
```

#### middleware - Generate Middleware Code

```bash
bingo make middleware <name> [-d dir] [-p package]

# Example
bingo make middleware auth
```

#### cmd - Generate Command Line Code

```bash
bingo make cmd <name> [-d dir] [-p package]

# Example
bingo make cmd serve
```

#### job - Generate Scheduled Job Code

```bash
bingo make job <name> [-d dir] [-p package]

# Example
bingo make job cleanup
```

#### migration - Database Migration

**Generate Migration File**

```bash
bingo make migration <name> [-d dir] [-p package] [-t table]

# Examples
bingo make migration create_users_table
bingo make migration create_posts_table -t posts
```

**Run Migrations**

```bash
bingo migrate <command> [options]

# Options
-v, --verbose   Show detailed compilation output
    --rebuild   Force recompile migration program
-f, --force     Force execution in production environment

# Subcommands
bingo migrate up          # Run all pending migrations
bingo migrate rollback    # Rollback the last batch of migrations
bingo migrate reset       # Rollback all migrations
bingo migrate refresh     # Rollback all and re-run migrations
bingo migrate fresh       # Drop all tables and re-run migrations
```

**Configure Migration Table Name** (optional, in `.bingo.yaml`):

```yaml
migrate:
  table: bingo_migration  # Default value
```

#### seeder - Generate Seeder File

```bash
bingo make seeder <name> [-d dir] [-p package]

# Example
bingo make seeder users
```

### db - Database Management

#### seed - Run Database Seeders

Run user-defined seeders to populate the database.

```bash
bingo db seed [options]

# Options
-v, --verbose      Show detailed compilation output
    --rebuild      Force recompile seeder program
    --seeder       Specify seeder class name to run

# Examples
bingo db seed                    # Run all seeders
bingo db seed --seeder=User      # Run only UserSeeder
bingo db seed -v                 # Show detailed output
```

#### service - Generate Service Module

Generate a complete service module with HTTP/gRPC/WebSocket server configuration.

```bash
bingo make service <name> [options]

# Server Options
--http                  Enable HTTP server
--grpc                  Enable gRPC server
--ws                    Enable WebSocket server

# Layer Options (when server is enabled, biz/router/handler are generated by default)
--no-biz                Don't generate business layer
--no-router             Don't generate router
--no-handler            Don't generate handler
--with-store            Generate store layer
--with-middleware       Generate middleware directory

# Examples
bingo make service api --http
bingo make service gateway --http --grpc
bingo make service realtime --ws
bingo make service chat --http --ws --with-store
bingo make service worker --no-biz
```

The generated service follows the `cmd/{app}-{service}/` naming convention. For example, if your root package is `github.com/myorg/demo` and you run `bingo make service admin`, it creates `cmd/demo-admin/main.go`.

### gen - Generate Code from Database

Auto-generate model code from database tables.

```bash
bingo gen -t <table1,table2,...>

# Examples
bingo gen -t users
bingo gen -t users,posts,comments
```

### version - Show Version

```bash
bingo version
```

## Usage Examples

### 1. Create New Project

```bash
# Create project (includes apiserver service by default)
bingo create github.com/myorg/blog

# Create project with all services
bingo create github.com/myorg/blog --all

# Create with specific services
bingo create github.com/myorg/blog --services apiserver,admserver

# Enter project directory
cd blog

# Generate complete CRUD code for user module
bingo make crud user

# Generate CRUD code for admserver service
bingo make crud user --service admserver
```

### 2. Generate Models from Database

```bash
# Generate models from existing database tables
bingo gen -t users,posts,comments
```

### 3. Generate New Service

```bash
# Generate an API service with HTTP server
bingo make service api --http --with-store

# Generate a WebSocket service
bingo make service realtime --ws

# Generate a service with HTTP and WebSocket
bingo make service chat --http --ws

# Generate a pure business processing worker service
bingo make service worker --no-biz
```

### 4. Generate Migrations and Seeders

```bash
# Generate database migration file
bingo make migration create_users_table

# Run migrations
bingo migrate up

# Generate seeder file
bingo make seeder users

# Run seeders
bingo db seed
```

## Directory Structure

Typical directory structure for a project created with bingo:

```
myapp/
├── cmd/                          # Command entry points
│   └── myapp/
├── internal/
│   ├── apiserver/
│   │   ├── biz/                 # Business logic layer
│   │   ├── handler/             # Handlers
│   │   ├── database/
│   │   │   ├── migration/       # Database migrations
│   │   │   └── seeder/          # Database seeders
│   │   ├── model/               # Data models
│   │   ├── router/              # Routes
│   │   └── store/               # Store layer
│   ├── pkg/
│   │   └── middleware/          # Middleware
│   └── watcher/
│       └── watcher/             # Scheduled jobs
├── pkg/
│   └── api/
│       └── v1/                  # API request/response definitions
├── .bingo.yaml                  # bingo configuration file
└── go.mod
```

## Development Workflow

1. **Initialize Project**: Create a new project with `bingo create`
2. **Configure Database**: Set up database connection in `.bingo.yaml`
3. **Generate Code**:
   - Use `bingo make crud` to quickly generate CRUD code
   - Use `bingo gen` to generate models from database
4. **Database Management**:
   - Use `bingo make migration` to create migration files
   - Use `bingo migrate up` to run migrations
   - Use `bingo make seeder` to create seeder files
   - Use `bingo db seed` to run seeders
5. **Extend Functionality**: Use `make` commands to generate other components as needed

## Development Checklist

### Core Features ✅
- [x] `bingo create` - Create project from GitHub template
- [x] `bingo make` - Code generation (model, store, biz, handler, etc.)
- [x] `bingo make service` - Generate complete service module (HTTP/gRPC/WebSocket)
- [x] `bingo gen` - Generate model code from database tables
- [x] `bingo migrate` - Database migration management (up, rollback, reset, refresh, fresh)
- [x] `bingo db seed` - Run database seeders
- [x] Service selection (`--services`, `--no-service`, `--add-service`, `--all`)
- [x] Make commands support multi-service (`--service` parameter for auto path inference)

### Pending Tasks 📋
- [ ] Cache management commands: `bingo cache list/clean` (future version)

### Documentation 📚
- [x] README updated with latest features
- [x] All new parameters documented
- [x] Usage examples cover main scenarios

## License

[License information]
