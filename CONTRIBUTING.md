# Contributing to git-proxy-go

感谢您对 git-proxy-go 项目的关注！本文档将帮助您了解如何为项目做出贡献。

## 开发环境

### 前置要求

- Go 1.21 或更高版本
- Git

### 设置开发环境

1. Fork 仓库
2. 克隆您的 Fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/git-proxy-go.git
   cd git-proxy-go
   ```
3. 添加上游仓库:
   ```bash
   git remote add upstream https://github.com/git-proxy/go.git
   ```
4. 安装依赖:
   ```bash
   go mod tidy
   ```

## 开发工作流

### 1. 创建功能分支

```bash
git checkout main
git pull upstream main
git checkout -b feature/your-feature-name
```

### 2. 编写代码

- 遵循 Go 代码规范
- 添加必要的注释
- 确保代码清晰易懂

### 3. 编写测试

```bash
# 运行所有测试
go test ./...

# 运行测试并查看覆盖率
go test -cover ./...

# 运行特定包的测试
go test -v ./internal/config/...
```

### 4. 提交代码

提交信息格式:
```
<type>: <subject>

<body>

<footer>
```

类型 (type):
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式（不影响功能）
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

示例:
```
feat: 添加 Prometheus 指标支持

添加 /metrics 端点用于暴露 Prometheus 格式的指标数据。
包括请求计数、延迟直方图和活跃连接数。

Closes #123
```

### 5. 推送并创建 Pull Request

```bash
git push origin feature/your-feature-name
```

然后在 GitHub 上创建 Pull Request。

## 代码规范

### Go 代码规范

- 使用 `gofmt` 格式化代码
- 遵循 Effective Go 指南
- 变量命名清晰、描述性强
- 函数应该简短，职责单一

### 命名规范

- 包名：简短、简洁、全小写
- 函数名：驼峰命名
- 常量名：全大写或驼峰
- 接口名：单数形式，如 `Logger`, `Proxy`

### 错误处理

- 错误应该被处理或显式返回
- 使用上下文丰富的错误信息
- 避免忽略 `_` 变量

## 项目结构

```
internal/
├── config/          # 配置加载和验证
├── proxy/           # HTTP/HTTPS 代理核心
├── ssh/             # SSH 代理
├── logger/          # 日志模块
└── middleware/       # HTTP 中间件
```

## 测试指南

- 每个包应该有对应的 `_test.go` 文件
- 测试函数以 `Test` 开头
- 使用 `t.Run` 组织相关测试用例
- 确保测试覆盖主要功能路径

示例:
```go
func TestConfigValidation(t *testing.T) {
    tests := []struct {
        name    string
        config  Config
        wantErr bool
    }{
        {
            name:    "valid config",
            config:  Config{...},
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.config.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Commit 规范

### 格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

### 示例

```
feat(proxy): 添加请求重试机制

当目标服务器返回临时错误时，自动重试请求。
最大重试次数可通过配置设置。

Closes #45
Fixes #46
```

## 问题反馈

### 创建 Issue

在创建 Issue 时，请包含:

- 清晰的标题和描述
- 复现步骤
- 期望行为 vs 实际行为
- Go 版本、操作系统信息
- 相关日志或错误信息

## 许可证

通过提交代码，您同意您的贡献将在 MIT 许可证下发布。
