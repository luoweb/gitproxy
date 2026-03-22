# Git Proxy Go

[![Build Status](https://github.com/git-proxy/go/workflows/Build/badge.svg)](https://github.com/git-proxy/go/actions)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Go 实现的 Git 代理程序，支持 HTTPS、HTTP、SSH 三种协议的代码仓库请求代理。

## 功能特性

- **HTTP/HTTPS 代理**: 作为反向代理转发 Git 客户端请求到目标 Git 服务器
- **SSH 代理**: 支持 SSH 协议的 git clone/fetch/push 操作
- **认证支持**: 可配置 HTTP Basic 认证
- **路径控制**: 支持白名单/黑名单路径控制
- **健康检查**: 提供 `/health` 端点用于监控
- **CORS 支持**: 支持跨域请求
- **TLS 安全**: 强制 TLS 1.2+，配置安全的加密套件
- **优雅关闭**: 支持信号处理，平滑停止服务

## 快速开始

### 前置要求

- Go 1.21 或更高版本
- Git 客户端

### 1. 配置文件

复制并编辑 `config.yaml`:

```yaml
server:
  host: "0.0.0.0"
  http_port: 8080
  https_port: 8443
  cert_file: "server.crt"
  key_file: "server.key"

target:
  host: "github.com"
  port: 443
  username: ""
  password: ""
  ssh:
    enabled: true
    port: 22

logging:
  level: "info"
  format: "text"

proxy:
  timeout: 30
  max_connections: 100
  path_rewrite: true
  allowed_paths: []
  blocked_paths: []
```

### 2. 构建

```bash
# 下载依赖
go mod tidy

# 构建
go build -o git-proxy cmd/proxy/main.go
```

### 3. 运行

```bash
# 使用默认配置
./git-proxy

# 指定配置文件
./git-proxy -config /path/to/config.yaml

# 查看版本
./git-proxy -version
```

### 4. 使用代理

```bash
# HTTP/HTTPS 协议
git clone http://localhost:8080/owner/repo.git
git clone https://localhost:8443/owner/repo.git

# SSH 协议 (需要配置 SSH 代理)
git clone ssh://git@localhost:22/owner/repo.git
```

### 5. 健康检查

```bash
curl http://localhost:8080/health
# {"status":"healthy"}
```

## 配置说明

### Server 配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `server.host` | 监听地址 | `0.0.0.0` |
| `server.http_port` | HTTP 端口 | `8080` |
| `server.https_port` | HTTPS 端口 | `8443` |
| `server.cert_file` | TLS 证书文件路径 | - |
| `server.key_file` | TLS 私钥文件路径 | - |

### Target 配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `target.host` | 目标 Git 服务器地址 | `github.com` |
| `target.port` | 目标端口 | `443` |
| `target.username` | 认证用户名 | - |
| `target.password` | 认证密码 | - |
| `target.ssh.enabled` | 启用 SSH 代理 | `true` |
| `target.ssh.port` | SSH 端口 | `22` |

### Logging 配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `logging.level` | 日志级别 (debug/info/warn/error) | `info` |
| `logging.format` | 日志格式 (text/json) | `text` |

### Proxy 配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `proxy.timeout` | 请求超时时间 (秒) | `30` |
| `proxy.max_connections` | 最大并发连接数 | `100` |
| `proxy.path_rewrite` | 启用路径重写 | `true` |
| `proxy.allowed_paths` | 允许的路径前缀列表 | `[]` |
| `proxy.blocked_paths` | 禁止的路径前缀列表 | `[]` |

## 项目结构

```
git-proxy-go/
├── cmd/proxy/
│   └── main.go              # 主程序入口
├── internal/
│   ├── config/              # 配置加载和验证
│   ├── proxy/               # HTTP/HTTPS 代理核心
│   ├── ssh/                 # SSH 代理
│   ├── logger/              # 日志模块
│   └── middleware/          # HTTP 中间件
├── config.yaml              # 配置文件示例
├── Dockerfile               # Docker 镜像构建
├── docker-compose.yaml      # Docker Compose 示例
├── CONTRIBUTING.md          # 贡献指南
├── LICENSE                  # MIT License
├── README.md                # 本文件
└── docs/
    └── DEPLOYMENT.md        # 部署文档
```

## 测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test -cover ./...

# 运行测试并输出详细日志
go test -v ./...
```

## Docker 部署

### 使用 Docker

```bash
# 构建镜像
docker build -t git-proxy .

# 运行容器
docker run -d \
  --name git-proxy \
  -p 8080:8080 \
  -p 8443:8443 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  git-proxy
```

### 使用 Docker Compose

```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

## 安全性

- TLS 1.2+ 强制启用
- 安全 TLS 加密套件配置
- Basic 认证支持
- 路径白名单/黑名单控制
- 请求超时保护
- Panic 恢复中间件

## 常见问题

### Q: 如何配置 HTTPS？

设置 `server.cert_file` 和 `server.key_file` 指向您的 TLS 证书和私钥文件。

### Q: 如何限制可访问的仓库？

使用 `proxy.allowed_paths` 或 `proxy.blocked_paths` 配置路径控制。

### Q: SSH 代理不工作？

确保 `target.ssh.enabled` 为 `true`，并正确配置 `target.ssh.port`。

## License

MIT License - 详见 [LICENSE](LICENSE) 文件

## 贡献

欢迎提交 Issue 和 Pull Request！详见 [CONTRIBUTING.md](CONTRIBUTING.md)
