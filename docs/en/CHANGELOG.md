# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.6.0] - 2025-12-01

### Added

- Add `bingo migrate` command for database migration management
  - `migrate up` - Run all pending migrations
  - `migrate rollback` - Rollback the last batch of migrations
  - `migrate reset` - Rollback all migrations
  - `migrate refresh` - Rollback all and re-run migrations
  - `migrate fresh` - Drop all tables and re-run migrations
  - Supports `--force` flag for production environment
  - Supports `--verbose` and `--rebuild` flags
- Add `bingo db seed` command for running database seeders
  - Supports `--seeder` flag to run specific seeder
  - Supports `--verbose` and `--rebuild` flags
- Add `bingo make migration` command to generate migration files
- Add `bingo make seeder` command to generate seeder files

### Changed

- Run `make protoc` before `go mod tidy` during project creation

## [1.5.0] - 2025-11-28

### Changed - Breaking Changes

**Template System Refactoring: From Built-in to Online Fetching**

This is a major architectural change, moving the template system from built-in to online fetching:

- **Online Template Fetching**: Download project templates from GitHub instead of embedding them in the binary
- **Template Caching**: Support local template caching for faster creation (cache location: `~/.bingo/templates/`)
- **Version/Branch Selection**: Support specifying template version or branch for project creation (`-r` flag)
- **Force Refresh**: Support `--no-cache` flag to force re-download templates
- **Mirror Configuration**: Support GitHub mirror configuration via `BINGO_TEMPLATE_MIRROR` environment variable

### Added

- Support specifying template version or branch via `-r` flag
- Support force refresh templates via `--no-cache` flag
- Support GitHub mirror configuration via environment variable

### Migration Guide

Upgrading from v1.4.x to v1.5.0:

- First use of `bingo create` will automatically download templates from GitHub
- If GitHub access is difficult, configure a mirror: `export BINGO_TEMPLATE_MIRROR=https://ghproxy.com/`
- To continue using the built-in template system, stay on v1.4.7

## [1.4.7] - 2024-XX-XX

Last version using the built-in template system.

### Features

- Built-in project templates, no network connection required
- Support generating code for all layers (model, store, biz, handler, etc.)
- Support generating model code from database tables
- Support multi-service architecture

---

## Version Notes

### v1.5+ - Online Template System

- ✅ Fetch project templates from GitHub online
- ✅ Support template caching for faster creation
- ✅ Support specifying version/branch for project creation
- ✅ Support mirror configuration for network access issues
- ✅ Flexible service component selection

### v1.4.x and Earlier - Built-in Template System

- ✅ Templates embedded directly in bingo
- ✅ No network connection required
- ✅ Template version bound to bingo version

[1.6.0]: https://github.com/bingo-project/bingoctl/compare/v1.5.0...v1.6.0
[1.5.0]: https://github.com/bingo-project/bingoctl/compare/v1.4.7...v1.5.0
[1.4.7]: https://github.com/bingo-project/bingoctl/releases/tag/v1.4.7
