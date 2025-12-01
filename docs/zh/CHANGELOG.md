# 更新日志

本项目的所有重要变更都将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
项目遵循 [语义化版本](https://semver.org/lang/zh-CN/spec/v2.0.0.html)。

## [1.6.0] - 2025-12-01

### 新增

- 新增 `bingo migrate` 数据库迁移管理命令
  - `migrate up` - 运行所有未执行的迁移
  - `migrate rollback` - 回滚最后一批迁移
  - `migrate reset` - 回滚所有迁移
  - `migrate refresh` - 回滚所有迁移并重新运行
  - `migrate fresh` - 删除所有表并重新运行迁移
  - 支持 `--force` 参数在生产环境强制执行
  - 支持 `--verbose` 和 `--rebuild` 参数
- 新增 `bingo db seed` 数据库填充命令
  - 支持 `--seeder` 参数指定要运行的 seeder
  - 支持 `--verbose` 和 `--rebuild` 参数
- 新增 `bingo make migration` 生成迁移文件命令
- 新增 `bingo make seeder` 生成 seeder 文件命令

### 变更

- 项目创建时先运行 `make protoc` 再运行 `go mod tidy`

## [1.5.0] - 2025-11-28

### 变更 - 重大变更

**模板系统重构：从内置模板改为在线拉取**

这是一个重大架构变更，将模板系统从内置改为在线拉取模式：

- **在线模板拉取**：从 GitHub 在线下载项目模板，而不是内置在二进制文件中
- **模板缓存机制**：支持本地缓存模板，加快创建速度（缓存位置：`~/.bingo/templates/`）
- **版本/分支选择**：支持指定模板版本或分支创建项目（`-r` 参数）
- **强制刷新**：支持 `--no-cache` 参数强制重新下载模板
- **镜像配置**：支持通过 `BINGO_TEMPLATE_MIRROR` 环境变量配置 GitHub 镜像

### 新增

- 支持通过 `-r` 参数指定模板版本或分支
- 支持通过 `--no-cache` 参数强制刷新模板
- 支持通过环境变量配置 GitHub 镜像

### 迁移指南

从 v1.4.x 升级到 v1.5.0：

- 首次使用 `bingo create` 时会自动从 GitHub 下载模板
- 如果网络访问 GitHub 有困难，可配置镜像：`export BINGO_TEMPLATE_MIRROR=https://ghproxy.com/`
- 如需继续使用内置模板系统，请继续使用 v1.4.7 版本

## [1.4.7] - 2024-XX-XX

最后一个使用内置模板系统的版本。

### 功能

- 内置项目模板，无需网络连接即可创建项目
- 支持生成各层代码（model, store, biz, controller 等）
- 支持从数据库表生成模型代码
- 支持多服务架构

---

## 版本说明

### v1.5+ - 在线模板系统

- ✅ 从 GitHub 在线拉取项目模板
- ✅ 支持模板缓存，加快创建速度
- ✅ 支持指定版本/分支创建项目
- ✅ 支持镜像配置，解决网络访问问题
- ✅ 灵活选择服务组件

### v1.4.x 及更早版本 - 内置模板系统

- ✅ 模板直接内置在 bingo 中
- ✅ 无需网络连接
- ✅ 模板版本与 bingo 版本绑定

[1.6.0]: https://github.com/bingo-project/bingoctl/compare/v1.5.0...v1.6.0
[1.5.0]: https://github.com/bingo-project/bingoctl/compare/v1.4.7...v1.5.0
[1.4.7]: https://github.com/bingo-project/bingoctl/releases/tag/v1.4.7
