# 前端说明

## 技术栈

- Vue3
- Vite
- Element Plus
- Axios
- Pinia
- Vue Router

## 启动

```bash
cd frontend
npm install
npm run dev
```

访问：

```text
http://127.0.0.1:5173
```

默认账号：

```text
admin / admin123
```

## 页面

- 登录页
- Dashboard
- 用户管理
- 角色管理
- 菜单管理
- 部门管理
- 设备管理
- 文件管理

Axios 拦截器位于 `src/utils/request.js`，会自动携带 JWT Token。路由守卫位于 `src/router/index.js`，会根据用户权限控制访问。
