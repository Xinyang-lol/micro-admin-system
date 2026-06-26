import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../store/auth'
import Layout from '../layout/Layout.vue'
import Login from '../views/Login.vue'
import Dashboard from '../views/Dashboard.vue'
import UserManage from '../views/UserManage.vue'
import RoleManage from '../views/RoleManage.vue'
import MenuManage from '../views/MenuManage.vue'
import DeptManage from '../views/DeptManage.vue'
import DeviceManage from '../views/DeviceManage.vue'
import FileManage from '../views/FileManage.vue'

export const menuRoutes = [
  { path: '/dashboard', name: 'Dashboard', component: Dashboard, meta: { title: '首页', icon: 'House', permission: 'dashboard:view' } },
  { path: '/system/users', name: 'UserManage', component: UserManage, meta: { title: '用户管理', icon: 'User', permission: 'sys:user:list' } },
  { path: '/system/roles', name: 'RoleManage', component: RoleManage, meta: { title: '角色管理', icon: 'Avatar', permission: 'sys:role:list' } },
  { path: '/system/menus', name: 'MenuManage', component: MenuManage, meta: { title: '菜单管理', icon: 'Menu', permission: 'sys:menu:list' } },
  { path: '/system/depts', name: 'DeptManage', component: DeptManage, meta: { title: '部门管理', icon: 'OfficeBuilding', permission: 'sys:dept:list' } },
  { path: '/devices', name: 'DeviceManage', component: DeviceManage, meta: { title: '设备管理', icon: 'Monitor', permission: 'device:list' } },
  { path: '/files', name: 'FileManage', component: FileManage, meta: { title: '文件管理', icon: 'FolderOpened', permission: 'file:list' } }
]

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/dashboard' },
    { path: '/login', component: Login },
    {
      path: '/',
      component: Layout,
      children: menuRoutes
    }
  ]
})

router.beforeEach(async (to) => {
  const store = useAuthStore()
  if (to.path === '/login') {
    return true
  }
  if (!store.isLoggedIn) {
    return '/login'
  }
  if (!store.profile) {
    await store.loadProfile()
  }
  const permission = to.meta?.permission
  if (permission && !store.hasPermission(permission)) {
    return '/dashboard'
  }
  return true
})

export default router
