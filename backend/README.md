# 后端说明

## 技术栈

- Go + Gin：HTTP API Gateway
- Go + gRPC：user-service、device-service、file-service
- MySQL：业务数据
- Redis：JWT Token 黑名单
- Consul：服务注册、发现、健康检查
- goroutine/channel/context：并发保存日志、设备统计、worker pool 造数、超时控制

## 启动基础组件

```bash
cd backend
docker compose up -d
```

MySQL 首次启动会自动执行 `scripts/init.sql`。如果已经有旧数据卷，需要重置：

```bash
docker compose down -v
docker compose up -d
```

## 启动微服务

打开 4 个终端：

```bash
cd backend
go run ./user-service
```

```bash
cd backend
go run ./device-service
```

```bash
cd backend
go run ./file-service
```

```bash
cd backend
go run ./api-gateway
```

默认端口：

| 服务 | 端口 |
| --- | --- |
| api-gateway | 8080 |
| user-service | 9001 |
| device-service | 9002 |
| file-service | 9003 |
| Consul UI | 8500 |
| MySQL | 3306 |
| Redis | 6379 |

## 环境变量

代码优先读取环境变量，`config/config.yaml` 是课堂展示用配置示例。

常用变量：

```bash
MYSQL_DSN='root:123456@tcp(127.0.0.1:3306)/micro_admin?charset=utf8mb4&parseTime=true&loc=Local'
CONSUL_ADDR='127.0.0.1:8500'
REDIS_ADDR='127.0.0.1:6379'
JWT_SECRET='micro-admin-secret'
SERVICE_HOST='127.0.0.1'
```

## 测试脚本

生成 10000 条设备数据：

```bash
cd backend
go run ./scripts/seed.go
```

1000 用户请求压测：

```bash
cd backend
go test ./scripts -run TestPressure1000Users -count=1 -v
```

可调整规模：

```bash
PRESSURE_TOTAL=1000 PRESSURE_CONCURRENCY=100 go test ./scripts -run TestPressure1000Users -count=1 -v
```

Windows PowerShell：

```powershell
$env:PRESSURE_TOTAL=1000
$env:PRESSURE_CONCURRENCY=100
go test ./scripts -run TestPressure1000Users -count=1 -v
```

## 验收点

- 登录账号：`admin`
- 登录密码：`admin123`
- 初始化 SQL：`scripts/init.sql`
- gRPC proto：`proto/user.proto`、`proto/device.proto`、`proto/file.proto`
- 接口文档：`API.md`
- Consul UI：`http://127.0.0.1:8500`

测试结果填写位置：

| 项目 | 结果 |
| --- | --- |
| 1000 用户并发耗时 | 待填写 |
| 10000 设备数据生成耗时 | 待填写 |
| 平均响应时间 | 待填写 |
| 最大响应时间 | 待填写 |
