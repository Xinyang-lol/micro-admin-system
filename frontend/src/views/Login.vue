<template>
  <div class="login-page">
    <div class="login-panel">
      <div class="login-title">
        <h1>微服务后台管理系统</h1>
        <p>Vue3 + Go + gRPC + MySQL + Redis + Consul</p>
      </div>
      <el-form :model="form" :rules="rules" ref="formRef" @keyup.enter="submit">
        <el-form-item prop="username">
          <el-input v-model="form.username" size="large" placeholder="用户名">
            <template #prefix><el-icon><User /></el-icon></template>
          </el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="form.password" size="large" type="password" show-password placeholder="密码">
            <template #prefix><el-icon><Lock /></el-icon></template>
          </el-input>
        </el-form-item>
        <el-button type="primary" size="large" :loading="loading" class="login-button" @click="submit">登录</el-button>
      </el-form>
      <div class="login-tip">默认账号：admin / admin123</div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'

const router = useRouter()
const store = useAuthStore()
const formRef = ref()
const loading = ref(false)
const form = reactive({ username: 'admin', password: 'admin123' })
const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const submit = async () => {
  await formRef.value.validate()
  loading.value = true
  try {
    await store.login(form)
    router.push('/dashboard')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  display: grid;
  place-items: center;
  min-height: 100vh;
  padding: 24px;
  background:
    linear-gradient(135deg, rgba(15, 23, 42, 0.86), rgba(37, 99, 235, 0.62)),
    url("https://images.unsplash.com/photo-1558494949-ef010cbdcc31?auto=format&fit=crop&w=1800&q=80") center/cover;
}

.login-panel {
  width: min(420px, 100%);
  padding: 28px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.96);
  box-shadow: 0 18px 50px rgba(15, 23, 42, 0.25);
}

.login-title h1 {
  margin: 0 0 8px;
  font-size: 26px;
  color: #111827;
}

.login-title p {
  margin: 0 0 24px;
  color: #6b7280;
}

.login-button {
  width: 100%;
}

.login-tip {
  margin-top: 16px;
  color: #6b7280;
  text-align: center;
}
</style>
