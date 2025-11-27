# Part 2 Template Reorganization Tasks

## Context

Part 2 (Create Command Service Selection) 的基础框架已完成,包括:
- ✅ 命令行 flags (--services, --no-service, --add-service)
- ✅ 服务列表计算逻辑
- ✅ 交互式选择框架

但是模板文件重组和选择性复制功能被有意推迟,需要在后续迭代中完成。

## Remaining Work

### 1. Reorganize template directory structure

**目标**: 重组模板目录,将服务特定文件移到子目录

**当前结构** (假设):
```
templates/
  ├── main.go
  ├── config.yaml
  ├── apiserver/
  ├── ctl/
  └── ...
```

**目标结构**:
```
templates/
  ├── common/           # 公共文件
  │   ├── .gitignore
  │   ├── go.mod
  │   ├── Makefile
  │   └── ...
  ├── services/         # 服务特定文件
  │   ├── apiserver/
  │   │   ├── cmd/
  │   │   ├── internal/
  │   │   └── ...
  │   ├── ctl/
  │   │   ├── cmd/
  │   │   └── ...
  │   ├── admserver/
  │   ├── bot/
  │   └── scheduler/
  └── config/           # 配置模板
      └── .bingoctl.yaml.tmpl
```

**子任务**:
- [ ] 识别公共文件 vs 服务特定文件
- [ ] 创建新的目录结构
- [ ] 移动文件到相应目录
- [ ] 更新 embed.FS 路径 (`//go:embed` 指令)
- [ ] 验证嵌入文件系统仍然工作

### 2. Implement selective template copying

**目标**: 修改 create 命令的 Run() 方法,只复制选中的服务

**当前行为**: 复制所有模板文件

**目标行为**: 根据 `selectedServices` 选择性复制

**实现步骤**:

#### 2.1 修改 Run() 方法

```go
func (o *CreateOptions) Run(args []string) error {
    // 1. 复制公共文件
    if err := copyCommonFiles(o.AppName); err != nil {
        return err
    }

    // 2. 为每个选中的服务复制文件
    for _, service := range o.selectedServices {
        if err := copyServiceFiles(o.AppName, service); err != nil {
            return err
        }
    }

    // 3. 生成动态配置文件
    if err := generateConfig(o.AppName, o.selectedServices); err != nil {
        return err
    }

    return nil
}
```

#### 2.2 实现 copyCommonFiles

```go
func copyCommonFiles(appName string) error {
    // 从 templates/common/ 复制文件
    // 包括: .gitignore, go.mod, Makefile, README.md, etc.
}
```

#### 2.3 实现 copyServiceFiles

```go
func copyServiceFiles(appName, service string) error {
    // 从 templates/services/{service}/ 复制文件
    // 应用模板变量替换 (AppName, ServiceName, etc.)
}
```

#### 2.4 实现 generateConfig

```go
func generateConfig(appName string, services []string) error {
    // 根据选中的服务生成 .bingoctl.yaml
    // 只包含选中服务的目录配置

    // 例如,如果只选择了 apiserver:
    // directory:
    //   model: internal/apiserver/model
    //   store: internal/apiserver/store
    //   ...

    // 如果选择了多个服务,使用第一个作为默认
}
```

**子任务**:
- [ ] 实现 copyCommonFiles 函数
- [ ] 实现 copyServiceFiles 函数
- [ ] 实现 generateConfig 函数
- [ ] 更新 Run() 方法调用这些函数
- [ ] 处理模板变量替换 (AppName, AppNameCamel, ServiceName, etc.)

### 3. Testing

**目标**: 验证所有服务选择场景都能正确工作

#### 3.1 最小骨架测试

```bash
bingoctl create minimal-app --services none
```

**预期结果**:
- 只创建公共文件 (go.mod, Makefile, .gitignore, etc.)
- 没有 cmd/ 目录
- 没有 internal/ 目录
- .bingoctl.yaml 为空或最小配置

#### 3.2 单服务测试

```bash
bingoctl create api-app --services apiserver
bingoctl create ctl-app --services ctl
bingoctl create admin-app --services admserver
```

**预期结果**:
- 只创建选中服务的文件
- .bingoctl.yaml 包含该服务的配置
- 项目可以正常构建和运行

#### 3.3 多服务组合测试

```bash
bingoctl create combo-app --services apiserver,admserver
bingoctl create full-app --services apiserver,ctl,bot
```

**预期结果**:
- 创建所有选中服务的文件
- .bingoctl.yaml 包含所有服务的配置
- 没有重复文件或冲突

#### 3.4 服务修改测试

```bash
bingoctl create modified-app --no-service ctl --add-service bot
```

**预期结果**:
- 创建 apiserver + bot (排除默认的 ctl)
- 配置正确反映选择的服务

#### 3.5 交互模式测试

```bash
bingoctl create interactive-app
# 使用交互式提示选择服务
```

**预期结果**:
- 显示服务选择界面
- 根据选择创建相应文件

**子任务**:
- [ ] 编写测试脚本自动化上述场景
- [ ] 验证生成的项目可以构建 (`go build`)
- [ ] 验证生成的项目可以运行
- [ ] 验证 bingoctl 命令在生成的项目中工作

### 4. Update Documentation

完成模板重组后需要更新文档:

- [ ] 更新 README.md create 命令部分
- [ ] 添加服务选择示例
- [ ] 说明 --services none 用法
- [ ] 说明如何组合服务
- [ ] 添加最佳实践建议

## Implementation Strategy

建议的实现顺序:

1. **Phase 1: 重组模板** (1-2 hours)
   - 分析当前模板结构
   - 创建新的目录结构
   - 移动文件
   - 更新 embed.FS

2. **Phase 2: 实现选择性复制** (2-3 hours)
   - 实现辅助函数
   - 修改 Run() 方法
   - 处理边缘情况

3. **Phase 3: 测试和修复** (1-2 hours)
   - 运行所有测试场景
   - 修复发现的问题
   - 验证构建和运行

4. **Phase 4: 文档更新** (0.5-1 hour)
   - 更新 README
   - 添加示例
   - 更新计划文档

## Estimated Effort

**总计**: 4-8 小时的集中工作

- 模板重组: 1-2 小时
- 选择性复制实现: 2-3 小时
- 测试和修复: 1-2 小时
- 文档更新: 0.5-1 小时

## Prerequisites

在开始之前:
- ✅ Part 1 已完成并测试
- ✅ Part 2 基础框架已完成
- ✅ 当前实现已通过代码审查
- ⚠️ 需要充分理解当前模板结构
- ⚠️ 需要制定模板文件分类策略

## Success Criteria

模板重组完成的标准:

1. ✅ 所有测试场景通过
2. ✅ 生成的项目可以正常构建和运行
3. ✅ 代码通过审查
4. ✅ 文档完整且准确
5. ✅ 没有破坏现有功能
6. ✅ 向后兼容(或有清晰的迁移路径)

## Notes

- 模板重组是一个较大的变更,建议在独立分支上进行
- 考虑使用 TDD 方法,先写集成测试
- 可能需要与项目维护者讨论模板分类策略
- 考虑未来添加更多服务的可扩展性
