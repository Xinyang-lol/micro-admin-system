<template>
  <div class="page">
    <div class="stat-grid">
      <div class="stat-card"><span class="muted">总数</span><strong>{{ stats.total || 0 }}</strong></div>
      <div class="stat-card"><span class="muted">在线</span><strong>{{ stats.online || 0 }}</strong></div>
      <div class="stat-card"><span class="muted">离线</span><strong>{{ stats.offline || 0 }}</strong></div>
      <div class="stat-card"><span class="muted">维修</span><strong>{{ stats.repair || 0 }}</strong></div>
    </div>

    <div class="panel">
      <div class="toolbar">
        <div class="toolbar-left">
          <el-input v-model="query.keyword" clearable placeholder="搜索设备名称/编码/位置" style="width: 240px" @keyup.enter="load" />
          <el-select v-model="query.typeId" clearable placeholder="设备类型" style="width: 150px">
            <el-option v-for="type in types" :key="type.id" :label="type.name" :value="type.id" />
          </el-select>
          <el-select v-model="query.status" clearable placeholder="状态" style="width: 130px">
            <el-option label="在线" value="online" />
            <el-option label="离线" value="offline" />
            <el-option label="维修" value="repair" />
          </el-select>
          <el-button @click="load"><el-icon><Search /></el-icon>查询</el-button>
        </div>
        <div class="toolbar-right">
          <el-button @click="typeDialogVisible = true"><el-icon><Collection /></el-icon>类型管理</el-button>
          <el-button type="primary" @click="openCreate"><el-icon><Plus /></el-icon>新增设备</el-button>
        </div>
      </div>
    </div>

    <div class="panel">
      <el-table :data="rows" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="设备名称" min-width="150" />
        <el-table-column prop="code" label="编码" min-width="130" />
        <el-table-column prop="type_name" label="类型" width="120" />
        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)">{{ statusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="location" label="位置" min-width="160" />
        <el-table-column prop="remark" label="备注" min-width="180" />
        <el-table-column label="操作" width="170" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" @click="openEdit(row)">编辑</el-button>
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

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑设备' : '新增设备'" width="560px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="设备名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="设备编码"><el-input v-model="form.code" /></el-form-item>
        <el-form-item label="设备类型">
          <el-select v-model="form.type_id" style="width: 100%">
            <el-option v-for="type in types" :key="type.id" :label="type.name" :value="type.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio-button label="online">在线</el-radio-button>
            <el-radio-button label="offline">离线</el-radio-button>
            <el-radio-button label="repair">维修</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="位置"><el-input v-model="form.location" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="form.remark" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="save">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="typeDialogVisible" title="设备类型管理" width="680px">
      <div class="toolbar">
        <span class="muted">设备类型用于设备分类和统计</span>
        <el-button type="primary" @click="openTypeCreate"><el-icon><Plus /></el-icon>新增类型</el-button>
      </div>
      <el-table :data="types" stripe style="margin-top: 12px">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="code" label="编码" />
        <el-table-column prop="remark" label="备注" />
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button text @click="openTypeEdit(row)">编辑</el-button>
            <el-button text type="danger" @click="removeType(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <el-dialog v-model="typeFormVisible" :title="typeForm.id ? '编辑类型' : '新增类型'" width="460px">
      <el-form :model="typeForm" label-width="80px">
        <el-form-item label="名称"><el-input v-model="typeForm.name" /></el-form-item>
        <el-form-item label="编码"><el-input v-model="typeForm.code" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="typeForm.remark" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="typeFormVisible = false">取消</el-button>
        <el-button type="primary" @click="saveType">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  createDevice,
  createDeviceType,
  deleteDevice,
  deleteDeviceType,
  deviceStatistics,
  listDeviceTypes,
  listDevices,
  updateDevice,
  updateDeviceType
} from '../api/device'

const loading = ref(false)
const rows = ref([])
const total = ref(0)
const types = ref([])
const stats = ref({})
const dialogVisible = ref(false)
const typeDialogVisible = ref(false)
const typeFormVisible = ref(false)
const query = reactive({ page: 1, pageSize: 10, keyword: '', typeId: '', status: '' })
const form = reactive(emptyForm())
const typeForm = reactive({ id: 0, name: '', code: '', remark: '' })

function emptyForm() {
  return { id: 0, name: '', code: '', type_id: 1, status: 'offline', location: '', remark: '' }
}

const statusText = (value) => ({ online: '在线', offline: '离线', repair: '维修' }[value] || value)
const statusType = (value) => ({ online: 'success', offline: 'info', repair: 'warning' }[value] || 'info')

const load = async () => {
  loading.value = true
  try {
    const params = { ...query, typeId: query.typeId || 0 }
    const res = await listDevices(params)
    rows.value = res.items || []
    total.value = res.total || 0
    stats.value = await deviceStatistics()
  } finally {
    loading.value = false
  }
}

const loadTypes = async () => {
  types.value = await listDeviceTypes()
  if (types.value[0] && !form.type_id) form.type_id = types.value[0].id
}

const openCreate = () => {
  Object.assign(form, emptyForm(), { type_id: types.value[0]?.id || 1 })
  dialogVisible.value = true
}

const openEdit = (row) => {
  Object.assign(form, emptyForm(), row)
  dialogVisible.value = true
}

const save = async () => {
  if (form.id) await updateDevice(form.id, form)
  else await createDevice(form)
  ElMessage.success('保存成功')
  dialogVisible.value = false
  load()
}

const remove = async (row) => {
  await ElMessageBox.confirm(`确定删除设备 ${row.name}？`, '删除确认', { type: 'warning' })
  await deleteDevice(row.id)
  ElMessage.success('删除成功')
  load()
}

const openTypeCreate = () => {
  Object.assign(typeForm, { id: 0, name: '', code: '', remark: '' })
  typeFormVisible.value = true
}

const openTypeEdit = (row) => {
  Object.assign(typeForm, row)
  typeFormVisible.value = true
}

const saveType = async () => {
  if (typeForm.id) await updateDeviceType(typeForm.id, typeForm)
  else await createDeviceType(typeForm)
  ElMessage.success('保存成功')
  typeFormVisible.value = false
  loadTypes()
}

const removeType = async (row) => {
  await ElMessageBox.confirm(`确定删除类型 ${row.name}？`, '删除确认', { type: 'warning' })
  await deleteDeviceType(row.id)
  ElMessage.success('删除成功')
  loadTypes()
}

onMounted(async () => {
  await loadTypes()
  load()
})
</script>

<style scoped>
.pagination {
  justify-content: flex-end;
  margin-top: 14px;
}
</style>
