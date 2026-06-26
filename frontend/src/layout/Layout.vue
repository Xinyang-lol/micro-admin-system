<template>
  <el-container class="app-shell">
    <el-aside class="app-aside" :width="collapsed ? '64px' : '224px'">
      <div class="brand">
        <el-icon><Grid /></el-icon>
        <span v-if="!collapsed">Micro Admin</span>
      </div>
      <el-menu :default-active="route.path" router :collapse="collapsed" class="side-menu">
        <el-menu-item v-for="item in visibleRoutes" :key="item.path" :index="item.path">
          <el-icon><component :is="item.meta.icon" /></el-icon>
          <template #title>{{ item.meta.title }}</template>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="app-header">
        <div class="header-left">
          <el-button circle text @click="collapsed = !collapsed">
            <el-icon><Fold v-if="!collapsed" /><Expand v-else /></el-icon>
          </el-button>
          <el-breadcrumb separator="/">
            <el-breadcrumb-item>后台管理</el-breadcrumb-item>
            <el-breadcrumb-item>{{ route.meta.title || '首页' }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <span>{{ store.profile?.nickname || store.profile?.username }}</span>
          <el-button text @click="doLogout">
            <el-icon><SwitchButton /></el-icon>
            退出
          </el-button>
        </div>
      </el-header>
      <el-main class="app-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'
import { menuRoutes } from '../router'

const collapsed = ref(false)
const route = useRoute()
const router = useRouter()
const store = useAuthStore()

const visibleRoutes = computed(() => menuRoutes.filter((item) => store.hasPermission(item.meta.permission)))

const doLogout = async () => {
  await store.logout()
  router.push('/login')
}
</script>

<style scoped>
.app-shell {
  height: 100vh;
  background: #f4f7fb;
}

.app-aside {
  color: #dbeafe;
  background: #182335;
  transition: width 0.2s ease;
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 56px;
  padding: 0 20px;
  color: #fff;
  font-weight: 700;
  white-space: nowrap;
}

.side-menu {
  border-right: 0;
  background: transparent;
}

.side-menu :deep(.el-menu-item) {
  color: #cbd5e1;
}

.side-menu :deep(.el-menu-item.is-active) {
  color: #fff;
  background: #2563eb;
}

.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 56px;
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
}

.header-left,
.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.app-main {
  height: calc(100vh - 56px);
  overflow: auto;
  padding: 16px;
}

@media (max-width: 720px) {
  .app-aside {
    position: fixed;
    z-index: 10;
    height: 100vh;
  }

  .app-main {
    padding-left: 76px;
  }
}
</style>
