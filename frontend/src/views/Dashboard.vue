<template>
  <div class="page">
    <div class="stat-grid">
      <div class="stat-card">
        <span class="muted">设备总数</span>
        <strong>{{ stats.total || 0 }}</strong>
      </div>
      <div class="stat-card">
        <span class="muted">在线设备</span>
        <strong>{{ stats.online || 0 }}</strong>
      </div>
      <div class="stat-card">
        <span class="muted">离线设备</span>
        <strong>{{ stats.offline || 0 }}</strong>
      </div>
      <div class="stat-card">
        <span class="muted">维修设备</span>
        <strong>{{ stats.repair || 0 }}</strong>
      </div>
    </div>
    <div class="panel">
      <div class="toolbar">
        <h3>设备类型统计</h3>
        <el-button :loading="loading" @click="load">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
      <el-table :data="typeRows" stripe>
        <el-table-column prop="name" label="类型" />
        <el-table-column prop="count" label="数量" width="160" />
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { deviceStatistics } from '../api/device'

const loading = ref(false)
const stats = ref({})
const typeRows = computed(() => Object.entries(stats.value.type_stats || {}).map(([name, count]) => ({ name, count })))

const load = async () => {
  loading.value = true
  try {
    stats.value = await deviceStatistics()
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
