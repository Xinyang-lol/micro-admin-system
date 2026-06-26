<template>
  <div class="page">
    <div class="panel">
      <div class="toolbar">
        <div class="toolbar-left">
          <el-input v-model="query.keyword" clearable placeholder="搜索文件名" style="width: 240px" @keyup.enter="load" />
          <el-button @click="load"><el-icon><Search /></el-icon>查询</el-button>
        </div>
        <el-upload :http-request="doUpload" :show-file-list="false">
          <el-button type="primary" :loading="uploading"><el-icon><Upload /></el-icon>上传文件</el-button>
        </el-upload>
      </div>
    </div>
    <div class="panel">
      <el-table :data="rows" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="original_name" label="原始文件名" min-width="220" />
        <el-table-column prop="content_type" label="类型" min-width="160" />
        <el-table-column label="大小" width="120">
          <template #default="{ row }">{{ formatSize(row.size) }}</template>
        </el-table-column>
        <el-table-column prop="created_at" label="上传时间" min-width="170" />
        <el-table-column label="操作" width="170" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" @click="download(row)">下载</el-button>
            <el-button text type="danger" @click="remove(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="query.page"
        v-model:page-size="query.pageSize"
        layout="total, sizes, prev, pager, next"
        :total="total"
        class="pagination"
        @current-change="load"
        @size-change="load"
      />
    </div>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { deleteFile, downloadFile, listFiles, uploadFile } from '../api/file'

const loading = ref(false)
const uploading = ref(false)
const rows = ref([])
const total = ref(0)
const query = reactive({ page: 1, pageSize: 10, keyword: '' })

const load = async () => {
  loading.value = true
  try {
    const res = await listFiles(query)
    rows.value = res.items || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

const doUpload = async ({ file }) => {
  uploading.value = true
  try {
    const data = new FormData()
    data.append('file', file)
    await uploadFile(data)
    ElMessage.success('上传成功')
    load()
  } finally {
    uploading.value = false
  }
}

const download = async (row) => {
  const resp = await downloadFile(row.id)
  const url = URL.createObjectURL(resp.data)
  const a = document.createElement('a')
  a.href = url
  a.download = row.original_name
  a.click()
  URL.revokeObjectURL(url)
}

const remove = async (row) => {
  await ElMessageBox.confirm(`确定删除文件 ${row.original_name}？`, '删除确认', { type: 'warning' })
  await deleteFile(row.id)
  ElMessage.success('删除成功')
  load()
}

const formatSize = (size) => {
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`
  return `${(size / 1024 / 1024).toFixed(1)} MB`
}

onMounted(load)
</script>

<style scoped>
.pagination {
  justify-content: flex-end;
  margin-top: 14px;
}
</style>
