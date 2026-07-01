# 基于微服务的后台管理系统

这是一个课程设计用的 B/S 后台管理系统。它不是双击某个 `.exe` 打开的程序，而是：

1. 先启动后端依赖：MySQL、Redis、Consul。
2. 再启动 4 个 Go 后端服务。
3. 再启动 Vue 前端。
4. 最后用浏览器打开 `http://127.0.0.1:5173`。

默认账号：

```text
用户名：admin
密码：admin123
```

## 一、最快打开方式（Windows）

如果你已经安装好了 Docker Desktop、Go、Node.js，可以直接双击项目根目录里的：

```text
start-windows.bat
```

脚本会自动做这些事情：

- 启动 MySQL、Redis、Consul。
- 打开 4 个后端服务窗口。
- 安装前端依赖。
- 启动前端页面。
- 自动打开浏览器访问系统。

如果双击后提示缺少 `docker`、`go` 或 `npm`，请先看下面的“运行前必须安装的软件”。

## 二、运行前必须安装的软件

请先安装这 3 个软件：

| 软件 | 作用 | 检查命令 |
| --- | --- | --- |
| Docker Desktop | 启动 MySQL、Redis、Consul | `docker --version` |
| Go 1.22 或更高版本 | 启动后端服务 | `go version` |
| Node.js 20 或更高版本 | 启动前端页面 | `node -v` |

Windows 安装后，如果命令行提示找不到 `go`、`node`、`npm` 或 `docker`，请重启电脑或重新打开 PowerShell。

Docker Desktop 安装后要先手动打开一次，等左下角显示 Docker 已运行，再启动本项目。

## 三、手动启动方式（推荐答辩前熟悉一遍）

下面步骤适合 Windows、Linux、macOS。Windows 推荐用 PowerShell。

### 第 1 步：进入项目目录

```powershell
cd E:\操作系统课设\micro-admin-system
```

如果你的项目放在别的位置，就把路径换成你自己的。

### 第 2 步：启动数据库、Redis、Consul

```powershell
cd backend
docker compose up -d
```

启动成功后，可以打开 Consul 页面检查：

```text
http://127.0.0.1:8500
```

### 第 3 步：启动 4 个后端服务

这一步要打开 4 个 PowerShell 窗口，每个窗口运行一个命令。

窗口 1：启动用户权限服务

```powershell
cd E:\操作系统课设\micro-admin-system\backend
go run ./user-service
```

窗口 2：启动设备管理服务

```powershell
cd E:\操作系统课设\micro-admin-system\backend
go run ./device-service
```

窗口 3：启动文件管理服务

```powershell
cd E:\操作系统课设\micro-admin-system\backend
go run ./file-service
```

窗口 4：启动 API 网关

```powershell
cd E:\操作系统课设\micro-admin-system\backend
go run ./api-gateway
```

看到类似下面的输出就说明后端启动成功：

```text
user-service started at :9001
device-service started at :9002
file-service started at :9003
api-gateway started at :8080
```

### 第 4 步：启动前端

再打开一个新的 PowerShell 窗口：

```powershell
cd E:\操作系统课设\micro-admin-system\frontend
npm install
npm run dev
```

看到类似下面的输出就说明前端启动成功：

```text
Local: http://localhost:5173/
```

### 第 5 步：浏览器打开系统

浏览器访问：

```text
http://127.0.0.1:5173
```

登录：

```text
admin / admin123
```

## 四、启动顺序不要乱

正确顺序是：

```text
Docker Desktop
  ↓
docker compose up -d
  ↓
user-service
device-service
file-service
api-gateway
  ↓
npm run dev
  ↓
浏览器打开 http://127.0.0.1:5173
```

如果顺序乱了，常见现象是前端页面能打开，但是登录失败或接口请求失败。

## 五、常见问题

### 1. 浏览器打不开 `http://127.0.0.1:5173`

说明前端没有启动成功。进入 `frontend` 目录重新执行：

```powershell
npm install
npm run dev
```

### 2. 登录失败，提示服务异常

通常是后端服务没有全部启动。请确认这 4 个窗口都在运行：

- user-service
- device-service
- file-service
- api-gateway

还要确认 Docker 里的 MySQL、Redis、Consul 已启动：

```powershell
cd backend
docker compose ps
```

### 3. `docker` 不是内部或外部命令

说明 Docker Desktop 没安装，或者安装后没有重启终端。请安装 Docker Desktop，然后重新打开 PowerShell。

### 4. `go` 不是内部或外部命令

说明 Go 没安装，或者环境变量没有生效。请安装 Go 1.22 以上版本，然后重新打开 PowerShell。

### 5. `npm` 无法运行或脚本被禁止

如果 PowerShell 提示 `npm.ps1 cannot be loaded`，可以用：

```powershell
npm.cmd install
npm.cmd run dev
```

### 6. MySQL 里没有初始化数据

如果之前启动过 Docker，旧数据卷可能还在。可以重置数据库：

```powershell
cd backend
docker compose down -v
docker compose up -d
```

注意：这个命令会删除 Docker 里的旧 MySQL 数据。

## 六、项目功能

前端页面：

- 登录页
- Dashboard 首页
- 用户管理
- 角色管理
- 菜单管理
- 部门管理
- 设备管理
- 文件管理

后端微服务：

- `api-gateway`：HTTP REST API 网关，端口 `8080`
- `user-service`：用户、角色、菜单、部门、登录认证，端口 `9001`
- `device-service`：设备管理、设备类型、设备统计，端口 `9002`
- `file-service`：文件元数据管理，端口 `9003`

基础组件：

- MySQL：业务数据库，端口 `3306`
- Redis：Token 黑名单，端口 `6379`
- Consul：服务注册发现，端口 `8500`

## 七、数据库初始化脚本

数据库脚本在：

```text
backend/scripts/init.sql
```

Docker 第一次启动 MySQL 时会自动执行这个脚本，创建表并插入管理员账号。

包含核心表：

- `sys_user`
- `sys_role`
- `sys_menu`
- `sys_dept`
- `sys_user_role`
- `sys_role_menu`
- `device`
- `device_type`
- `file_info`

## 八、接口文档和答辩稿

接口文档：

```text
backend/API.md
```

答辩稿：

```text
docs/defense-script.md
```

截图可以放在：

```text
docs/screenshots/
```

## 九、测试和造数

生成 10000 条设备数据：

```powershell
cd backend
go run ./scripts/seed.go
```

1000 用户请求压测：

```powershell
cd backend
go test ./scripts -run TestPressure1000Users -count=1 -v
```

## 十、GitHub 地址

```text
https://github.com/Xinyang-lol/micro-admin-system
```
