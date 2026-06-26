# 基于微服务的后台管理系统答辩稿

## 一、开场介绍

各位老师好，我本次课程设计的题目是《基于微服务的后台管理系统的设计与实现》。

本项目实现了一个 B/S 架构的后台管理系统，前端使用 Vue3、Vite、Element Plus、Axios、Pinia 和 Vue Router，后端使用 Go 语言实现 API Gateway 和多个微服务。系统包含用户权限服务、设备管理服务和文件管理服务，微服务之间通过 gRPC 通信，网关对外提供 REST API。

项目的主要目标是把操作系统课程中学习到的进程、线程、并发、进程间通信等思想，结合后端微服务架构进行实践。系统支持登录认证、RBAC 权限控制、设备管理、文件管理、服务注册发现、健康检查和简单负载均衡。

## 二、系统整体架构

系统整体分为前端、网关、微服务和基础设施四层。

第一层是前端页面，使用 Vue3 实现后台管理界面，包括登录页、首页 Dashboard、用户管理、角色管理、菜单管理、部门管理、设备管理和文件管理页面。

第二层是 API Gateway，也就是网关服务。网关使用 Go + Gin 实现，对外提供 HTTP REST API，例如 `/api/auth/login`、`/api/users`、`/api/devices`、`/api/files/upload` 等接口。前端所有请求都先进入网关。

第三层是微服务层，包括三个服务：

1. user-service：负责登录认证、用户、角色、菜单、部门和权限校验。
2. device-service：负责设备信息、设备类型和设备统计。
3. file-service：负责文件元数据保存、查询和删除。

第四层是基础组件，包括 MySQL、Redis 和 Consul。MySQL 保存业务数据，Redis 用来实现 Token 黑名单，Consul 用来实现服务注册、服务发现和健康检查。

系统调用流程是：浏览器访问前端页面，前端通过 Axios 请求 API Gateway，网关校验 JWT 和权限，然后从 Consul 发现对应微服务实例，再通过 gRPC 调用具体服务。业务服务访问 MySQL 完成数据操作，最后结果返回给前端。

## 三、主要功能实现

本系统主要实现了以下功能。

第一是登录认证模块。用户输入用户名和密码后，网关调用 user-service 的 Login 方法。user-service 从数据库查询用户信息，并使用 bcrypt 校验密码。校验成功后生成 JWT Token，前端保存 Token，后续请求通过 Axios 拦截器自动携带 `Authorization` 请求头。

第二是 RBAC 权限控制。数据库中设计了用户表、角色表、菜单表、部门表、用户角色关联表和角色菜单关联表。用户登录后会获得角色和权限标识，例如 `sys:user:list`、`device:create`、`file:upload`。网关中间件会根据当前接口所需权限进行拦截，如果权限不足，会返回明确的 403 响应。

第三是用户、角色、菜单和部门管理。用户管理支持分页查询、新增、修改、删除、启用禁用、重置密码和分配角色。角色管理支持新增、修改、删除和分配菜单权限。菜单和部门使用树形结构展示，符合后台管理系统常见设计。

第四是设备管理模块。设备管理支持设备分页查询、新增、修改、删除、状态管理、设备类型管理和设备数量统计。统计接口会返回设备总数、在线数量、离线数量、维修数量以及按设备类型分组的数量。

第五是文件管理模块。文件上传接口由网关接收 multipart 文件并保存到本地上传目录，然后调用 file-service 保存文件元数据。文件列表、下载和删除也都通过 file-service 管理元数据。

第六是服务治理模块。三个微服务启动后都会注册到 Consul，并使用 TTL 健康检查定时上报状态。网关调用服务时，会从 Consul 查询健康实例，并采用简单轮询方式选择服务实例，实现基础负载均衡。

## 四、数据库设计

数据库使用 MySQL，共设计了九张核心表：

1. `sys_user`：用户表，保存用户名、bcrypt 密码、昵称、手机号、邮箱、状态和部门 ID。
2. `sys_role`：角色表，保存角色名称、角色编码和状态。
3. `sys_menu`：菜单表，保存菜单树、路由、组件名和权限标识。
4. `sys_dept`：部门表，保存部门树。
5. `sys_user_role`：用户角色关联表。
6. `sys_role_menu`：角色菜单关联表。
7. `device`：设备表。
8. `device_type`：设备类型表。
9. `file_info`：文件元数据表。

初始化脚本位于 `backend/scripts/init.sql`，里面已经插入了管理员账号 `admin`，默认密码是 `admin123`，密码使用 bcrypt 加密保存，不是明文。

## 五、gRPC 与微服务通信

项目中定义了三个 proto 文件：

1. `user.proto`
2. `device.proto`
3. `file.proto`

user-service 提供 Login、GetUserInfo、ListUsers、CreateUser、UpdateUser、DeleteUser、ListRoles、ListMenus、CheckPermission 等方法。

device-service 提供 ListDevices、CreateDevice、UpdateDevice、DeleteDevice、GetDeviceStatistics 等方法。

file-service 提供 UploadFileMeta、ListFiles、GetFile、DeleteFile 等方法。

API Gateway 本身不直接处理这些业务数据库逻辑，而是通过 gRPC 调用微服务。这样做可以体现微服务拆分和进程间通信思想，也方便后续对单个服务独立扩展。

## 六、并发与操作系统课程相关点

本项目重点体现了操作系统课程中的并发和进程通信思想。

第一，多个服务本质上是多个独立进程。api-gateway、user-service、device-service 和 file-service 可以分别启动、分别监听不同端口，并通过 gRPC 进行进程间通信。

第二，Go 的 goroutine 用来实现轻量级并发。比如每个服务向 Consul 注册后，会启动 goroutine 定时进行健康上报。

第三，网关调用微服务时使用 `context.WithTimeout` 设置超时时间，避免服务调用长时间阻塞，提高系统稳定性。

第四，device-service 的设备统计接口使用多个 goroutine 并发查询设备总数、在线数量、离线数量、维修数量和类型统计，然后通过 channel 汇总结果。

第五，file-service 在保存文件元数据后，使用 goroutine 和 channel 异步记录上传日志，避免日志处理影响主请求响应。

第六，`scripts/seed.go` 使用 worker pool 模式批量生成 10000 条设备数据，体现了任务队列、工作线程和结果汇总的思想。

## 七、前端实现

前端采用 Vue3 单页应用结构。

登录页面调用 `/api/auth/login` 获取 Token。登录成功后，Pinia 保存用户信息和 Token，Axios 请求拦截器会自动携带 Token。Vue Router 路由守卫会判断用户是否登录，以及是否拥有当前页面权限。

后台布局包括侧边栏、顶部栏、面包屑和主内容区域。侧边栏根据用户权限显示可访问页面。业务页面使用 Element Plus 的表格、表单、弹窗、分页、树形选择等组件实现，整体界面清晰、可操作。

## 八、运行与测试

项目支持 Windows、Linux 和 macOS 运行。

启动基础组件：

```bash
cd backend
docker compose up -d
```

启动后端服务：

```bash
go run ./user-service
go run ./device-service
go run ./file-service
go run ./api-gateway
```

启动前端：

```bash
cd frontend
npm install
npm run dev
```

浏览器访问：

```text
http://127.0.0.1:5173
```

默认账号密码：

```text
admin / admin123
```

生成 10000 条设备测试数据：

```bash
cd backend
go run ./scripts/seed.go
```

进行 1000 用户请求压测：

```bash
go test ./scripts -run TestPressure1000Users -count=1 -v
```

压测脚本会模拟并发登录和访问用户信息接口，用于验证系统在并发请求下的响应能力。

## 九、项目亮点

我认为本项目的主要亮点有四点。

第一，项目不是单体后端，而是拆成了 API Gateway 和三个微服务，服务职责比较清晰。

第二，网关不直接访问所有数据库，而是通过 gRPC 调用微服务，体现了微服务之间的进程通信。

第三，实现了 Consul 服务注册发现和 TTL 健康检查，网关可以从注册中心发现服务实例，并进行简单轮询负载均衡。

第四，项目中多处使用 goroutine、channel 和 context，例如健康检查、文件日志、设备统计和批量造数，能够体现操作系统课程中的并发编程思想。

## 十、总结

通过本次课程设计，我完成了一个前后端分离、基于微服务架构的后台管理系统。项目实现了登录认证、RBAC 权限控制、用户管理、角色管理、菜单管理、部门管理、设备管理、文件管理、服务注册发现和并发测试等功能。

在实现过程中，我加深了对进程、线程、协程、进程间通信、服务治理和并发控制的理解。同时也熟悉了 Go 语言在微服务开发中的使用方式，以及 Vue3 在后台管理系统中的实际应用。

我的汇报到此结束，请各位老师批评指正。

## 十一、可能被问到的问题与回答

### 1. 为什么使用 API Gateway？

API Gateway 可以统一对外提供 REST API，前端不需要直接调用多个微服务。同时网关可以集中处理 JWT 校验、权限控制、跨域、服务发现和超时控制，让各个微服务专注自己的业务逻辑。

### 2. 为什么微服务之间使用 gRPC？

gRPC 基于 HTTP/2 和 proto 定义接口，适合服务之间的高效通信。它比普通 HTTP JSON 调用更适合内部服务通信，也能清晰约束服务接口。

### 3. RBAC 是怎么实现的？

RBAC 分为用户、角色、菜单/权限三层。用户通过 `sys_user_role` 关联角色，角色通过 `sys_role_menu` 关联菜单和权限。接口访问时，网关根据 JWT 中的权限或调用 user-service 的 CheckPermission 方法判断是否允许访问。

### 4. Redis 在项目中有什么作用？

Redis 用于保存退出登录后的 JWT Token 黑名单。用户登出后，当前 Token 会写入 Redis，在 Token 过期前如果再次使用这个 Token 访问接口，网关会拒绝请求。

### 5. Consul 在项目中有什么作用？

Consul 用来实现服务注册、服务发现和健康检查。微服务启动后向 Consul 注册，网关调用微服务前会从 Consul 查询健康实例，再选择一个实例进行 gRPC 调用。

### 6. 项目如何体现负载均衡？

网关内部维护每个服务的轮询计数器。每次调用服务时，先从 Consul 获取健康实例列表，然后按照轮询方式选择一个实例，实现简单的客户端负载均衡。

### 7. 项目如何体现并发？

项目中使用 goroutine 和 channel 处理并发任务。比如设备统计接口并发查询多个统计指标，最后通过 channel 汇总；文件服务异步记录上传日志；造数脚本使用 worker pool 并发插入设备数据；健康检查使用 goroutine 定时上报。

### 8. 如果某个微服务挂了怎么办？

微服务挂掉后，Consul 的 TTL 健康检查会变成异常状态。网关服务发现时只查询健康实例，因此不会继续调用异常实例。如果所有实例都不可用，网关会返回明确的服务异常响应。

### 9. 密码是如何保证安全的？

密码使用 bcrypt 加密后保存到数据库，不保存明文。用户登录时使用 bcrypt 校验输入密码和数据库中的哈希值。

### 10. 这个系统还有哪些可以改进的地方？

后续可以增加更多测试用例、接入对象存储保存文件、增加操作日志、使用 Dockerfile 容器化后端服务、增加链路追踪和 Prometheus 监控等功能。
