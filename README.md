# BingoCTL

Go CLI 工具

功能规划

| 进度   | 命令         | 说明                           |
|------|------------|------------------------------|
| -[x] | make 命令    | 生成 cmd, crud, middleware 等代码 |
| -[ ] | create 命令  | 从零创建项目脚手架代                   |
| -[ ] | migrate 命令 | 通过 struct 生成表                |
| -[ ] | create 命令  | 通过 sql 生成 struct             |

# Usage

Install

```bash
go install github.com/bingo-project/bingoctl@latest
```

在项目根目录下创建配置文件，`touch .bingoctl.yaml`, 并写入以下内容

```yaml
version: v1

root-package: your-pkg-name

directory:
  cmd: internal/bingoctl/cmd
  model: internal/pkg/model
  store: internal/apiserver/store
  request: pkg/api/v1
  biz: internal/apiserver/biz
  controller: internal/apiserver/controller
  middleware: internal/pkg/middleware
```

查看命令

```bash
bingoctl make -h
bingoctl make [COMMAND] -h
```

生成 cmd 代码

```bash
bingoctl make cmd NAME [-d dir] [-p package]
```

生成 crud 代码
自动生成 model store biz controller request 代码

```bash
bingoctl make crud NAME
```

生成单文件代码

```bash
bingoctl make model NAME [-d dir] [-p package]
bingoctl make store NAME [-d dir] [-p package]
bingoctl make request NAME [-d dir] [-p package]
bingoctl make biz NAME [-d dir] [-p package]
bingoctl make controller NAME [-d dir] [-p package]
bingoctl make middleware NAME [-d dir] [-p package]
```
