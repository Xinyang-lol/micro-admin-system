# 接口文档

统一返回格式：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

错误示例：

```json
{
  "code": 401,
  "message": "Token 无效或已过期",
  "data": {}
}
```

除登录接口外，请在请求头携带：

```text
Authorization: Bearer <token>
```

## 认证

| 方法 | 路径 | 说明 | Body |
| --- | --- | --- | --- |
| POST | `/api/auth/login` | 登录 | `{"username":"admin","password":"admin123"}` |
| POST | `/api/auth/logout` | 登出并加入 Redis Token 黑名单 | 无 |
| GET | `/api/auth/profile` | 当前用户信息 | 无 |

## 用户

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/users?page=1&pageSize=10&keyword=` | 用户分页 |
| POST | `/api/users` | 新增用户 |
| PUT | `/api/users/:id` | 修改用户 |
| DELETE | `/api/users/:id` | 删除用户 |
| PUT | `/api/users/:id/status` | 启用/禁用，Body: `{"status":1}` |
| PUT | `/api/users/:id/password` | 重置密码，Body: `{"password":"123456"}` |
| PUT | `/api/users/:id/roles` | 分配角色，Body: `{"role_ids":[1,2]}` |

用户保存 Body：

```json
{
  "username": "test",
  "password": "123456",
  "nickname": "测试用户",
  "email": "test@example.com",
  "phone": "13900000000",
  "status": 1,
  "dept_id": 2,
  "role_ids": [2]
}
```

## 角色、菜单、部门

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/roles?page=1&pageSize=10&keyword=` | 角色列表 |
| POST | `/api/roles` | 新增角色 |
| PUT | `/api/roles/:id` | 修改角色 |
| DELETE | `/api/roles/:id` | 删除角色 |
| PUT | `/api/roles/:id/menus` | 分配菜单权限，Body: `{"menu_ids":[1,2,3]}` |
| GET | `/api/menus/tree` | 菜单树 |
| POST | `/api/menus` | 新增菜单 |
| PUT | `/api/menus/:id` | 修改菜单 |
| DELETE | `/api/menus/:id` | 删除菜单 |
| GET | `/api/depts/tree` | 部门树 |
| POST | `/api/depts` | 新增部门 |
| PUT | `/api/depts/:id` | 修改部门 |
| DELETE | `/api/depts/:id` | 删除部门 |

## 设备

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/devices?page=1&pageSize=10&keyword=&typeId=&status=` | 设备分页 |
| POST | `/api/devices` | 新增设备 |
| PUT | `/api/devices/:id` | 修改设备 |
| DELETE | `/api/devices/:id` | 删除设备 |
| GET | `/api/devices/statistics` | 设备数量统计 |
| GET | `/api/device-types` | 设备类型 |
| POST | `/api/device-types` | 新增设备类型 |
| PUT | `/api/device-types/:id` | 修改设备类型 |
| DELETE | `/api/device-types/:id` | 删除设备类型 |

设备保存 Body：

```json
{
  "name": "实验室交换机",
  "code": "SW-1001",
  "type_id": 2,
  "status": "online",
  "location": "实验楼 301",
  "remark": "课程设计演示"
}
```

## 文件

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/api/files/upload` | multipart 上传，字段名 `file` |
| GET | `/api/files?page=1&pageSize=10&keyword=` | 文件列表 |
| GET | `/api/files/:id/download` | 文件下载 |
| DELETE | `/api/files/:id` | 删除文件 |

## gRPC 服务

gRPC 源文件位于 `backend/proto`，Go 绑定位于 `backend/proto/gen`。

| 服务 | 地址默认值 | 职责 |
| --- | --- | --- |
| user-service | `127.0.0.1:9001` | 登录认证、用户、角色、菜单、部门、权限检查 |
| device-service | `127.0.0.1:9002` | 设备、设备类型、设备统计 |
| file-service | `127.0.0.1:9003` | 文件元数据 |

网关使用 Consul 健康实例列表做轮询发现，并通过 `context.WithTimeout` 控制每次 gRPC 调用在 3 秒内返回。
