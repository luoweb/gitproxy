# 部署指南

本文档提供 git-proxy-go 的多种部署方案。

## 目录

- [Docker 部署](#docker-部署)
- [Systemd 部署](#systemd-部署)
- [Kubernetes 部署](#kubernetes-部署)
- [生产环境配置](#生产环境配置)
- [反向代理配置](#反向代理配置)
- [监控配置](#监控配置)

## Docker 部署

### 构建镜像

```bash
docker build -t git-proxy:latest .
```

### 运行容器

```bash
docker run -d \
  --name git-proxy \
  --restart unless-stopped \
  -p 8080:8080 \
  -p 8443:8443 \
  -v /path/to/config.yaml:/app/config.yaml:ro \
  git-proxy:latest
```

### Docker Compose

```yaml
version: '3.8'

services:
  git-proxy:
    image: git-proxy:latest
    container_name: git-proxy
    restart: unless-stopped
    ports:
      - "8080:8080"
      - "8443:8443"
    volumes:
      - ./config.yaml:/app/config.yaml:ro
    environment:
      - TZ=Asia/Shanghai
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

## Systemd 部署

### 创建用户

```bash
sudo useradd -r -s /sbin/nologin git-proxy
```

### 创建服务文件

```bash
sudo vim /etc/systemd/system/git-proxy.service
```

```ini
[Unit]
Description=Git Proxy Server
After=network.target

[Service]
Type=simple
User=git-proxy
Group=git-proxy
WorkingDirectory=/opt/git-proxy
ExecStart=/opt/git-proxy/git-proxy -config /opt/git-proxy/config.yaml
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

### 部署步骤

```bash
# 编译二进制文件
go build -o git-proxy cmd/proxy/main.go

# 复制文件
sudo cp git-proxy /opt/git-proxy/
sudo cp config.yaml /opt/git-proxy/
sudo chown -R git-proxy:git-proxy /opt/git-proxy

# 重新加载 systemd
sudo systemctl daemon-reload

# 启动服务
sudo systemctl enable git-proxy
sudo systemctl start git-proxy

# 检查状态
sudo systemctl status git-proxy
```

## Kubernetes 部署

### ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: git-proxy-config
data:
  config.yaml: |
    server:
      host: "0.0.0.0"
      http_port: 8080
      https_port: 8443
    target:
      host: "github.com"
      port: 443
      ssh:
        enabled: false
        port: 22
    logging:
      level: "info"
    proxy:
      timeout: 30
      max_connections: 100
```

### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: git-proxy
spec:
  replicas: 2
  selector:
    matchLabels:
      app: git-proxy
  template:
    metadata:
      labels:
        app: git-proxy
    spec:
      containers:
        - name: git-proxy
          image: git-proxy:latest
          ports:
            - containerPort: 8080
              name: http
            - containerPort: 8443
              name: https
          volumeMounts:
            - name: config
              mountPath: /app/config.yaml
              subPath: config.yaml
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 10
            periodSeconds: 30
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 10
          resources:
            requests:
              memory: "64Mi"
              cpu: "100m"
            limits:
              memory: "128Mi"
              cpu: "500m"
      volumes:
        - name: config
          configMap:
            name: git-proxy-config
```

### Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: git-proxy
spec:
  type: LoadBalancer
  selector:
    app: git-proxy
  ports:
    - port: 80
      targetPort: 8080
      name: http
    - port: 443
      targetPort: 8443
      name: https
```

## 生产环境配置

### 完整配置文件示例

```yaml
server:
  host: "0.0.0.0"
  http_port: 8080
  https_port: 8443
  cert_file: "/etc/ssl/git-proxy.crt"
  key_file: "/etc/ssl/git-proxy.key"

target:
  host: "your-gitlab.example.com"
  port: 443
  username: "deploy"
  password: "${GIT_PASSWORD}"
  ssh:
    enabled: true
    port: 22

logging:
  level: "info"
  format: "json"

proxy:
  timeout: 60
  max_connections: 200
  path_rewrite: true
  allowed_paths:
    - "/group/project1"
    - "/group/project2"
  blocked_paths:
    - "/internal"
    - "/archive"
```

### 环境变量

密码等敏感信息建议使用环境变量：

```bash
export GIT_PASSWORD="your-secure-password"
./git-proxy -config config.yaml
```

## 反向代理配置

### Nginx

```nginx
upstream git_proxy {
    server 127.0.0.1:8080;
    keepalive 32;
}

server {
    listen 80;
    server_name git-proxy.example.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name git-proxy.example.com;

    ssl_certificate /etc/ssl/git-proxy.crt;
    ssl_certificate_key /etc/ssl/git-proxy.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    client_max_body_size 2G;

    location / {
        proxy_pass http://git_proxy;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    location /health {
        proxy_pass http://git_proxy/health;
        access_log off;
    }
}
```

## 监控配置

### Prometheus 抓取配置

```yaml
scrape_configs:
  - job_name: 'git-proxy'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

### Grafana Dashboard

建议监控指标：
- HTTP 请求率和延迟
- SSH 连接数
- 错误率
- 活跃连接数

## 性能调优

### 文件描述符限制

```bash
# /etc/security/limits.conf
git-proxy soft nofile 65536
git-proxy hard nofile 65536
```

### 内核参数

```bash
# /etc/sysctl.conf
net.core.somaxconn = 1024
net.ipv4.tcp_max_syn_backlog = 2048
```

## 备份策略

配置文件应该纳入版本控制或配置管理：

```bash
# 使用 git 管理配置
git add config.yaml
git commit -m "Update production config"
```

## 安全建议

1. 使用 TLS 终止 HTTPS
2. 配置防火墙规则
3. 定期更新证书
4. 使用最小权限原则
5. 启用日志监控
6. 定期审计访问日志
