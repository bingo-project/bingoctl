# BingoCTL

Go CLI 工具

功能规划

| 进度   | 命令        | 说明                           |
|------|-----------|------------------------------|
| -[x] | make 命令   | 生成 cmd, crud, middleware 等代码 |
| -[x] | create 命令 | 从零创建项目脚手架代                   |
| -[x] | gen 命令    | 通过 sql 生成 struct             |

# Usage

Install

```bash
go install github.com/bingo-project/bingoctl@latest
```

在项目根目录下创建配置文件，`touch .bingoctl.yaml`, 并写入以下内容

```yaml
version: v1

rootPackage: bingo

directory:
  cmd: internal/bingoctl/cmd
  model: internal/pkg/model
  store: internal/apiserver/store
  request: pkg/api/v1
  biz: internal/apiserver/biz
  controller: internal/apiserver/controller/v1
  middleware: internal/pkg/middleware
  job: internal/watcher/watcher

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

## 查看命令

```bash
bingoctl create -h
bingoctl make -h
bingoctl make [COMMAND] -h
```

## Create

```bash
bingoctl create NAME
```

## Make

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
bingoctl make job NAME [-d dir] [-p package]
```
