# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.5.0] - 2025-01-29

### Changed - 重大变更

**模板系统重构：从内置模板改为在线拉取**

这是一个重大架构变更，将模板系统从内置改为在线拉取模式：

- **在线模板拉取**：从 GitHub 在线下载项目模板，而不是内置在二进制文件中
- **模板缓存机制**：支持本地缓存模板，加快创建速度（缓存位置：`~/.bingo/templates/`）
- **版本/分支选择**：支持指定模板版本或分支创建项目（`-r` 参数）
- **强制刷新**：支持 `--no-cache` 参数强制重新下载模板
- **镜像配置**：支持通过 `BINGO_TEMPLATE_MIRROR` 环境变量配置 GitHub 镜像

### Added

- 支持通过 `-r` 参数指定模板版本或分支
- 支持通过 `--no-cache` 参数强制刷新模板
- 支持通过环境变量配置 GitHub 镜像

### Migration Guide

从 v1.4.x 升级到 v1.5.0：

- 首次使用 `bingo create` 时会自动从 GitHub 下载模板
- 如果网络访问 GitHub 有困难，可配置镜像：`export BINGO_TEMPLATE_MIRROR=https://ghproxy.com/`
- 如需继续使用内置模板系统，请继续使用 v1.4.7 版本

## [1.4.7] - 2024-XX-XX

最后一个使用内置模板系统的版本。

### Features

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

[1.5.0]: https://github.com/bingo-project/bingoctl/compare/v1.4.7...v1.5.0
[1.4.7]: https://github.com/bingo-project/bingoctl/releases/tag/v1.4.7
