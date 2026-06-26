<template>
  <div class="page">
    <div class="panel">
      <div class="toolbar">
        <div class="toolbar-left">
          <el-input v-model="query.keyword" clearable placeholder="搜索用户名/昵称/手机号" style="width: 240px" @keyup.enter="load" />
          <el-button @click="load">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
        </div>
        <el-button type="primary" @click="openCreate">
          <el-icon><Plus /></el-icon>
          新增用户
        </el-button>
      </div>
    </div>
    <div class="panel">
      <el-table :data="rows" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="nickname" label="昵称" min-width="120" />
        <el-table-column prop="phone" label="手机号" min-width="130" />
        <el-table-column prop="email" label="邮箱" min-width="180" />
        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <el-switch :model-value="row.status === 1" @change="(value) => changeStatus(row, value)" />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" min-width="170" />
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button text @click="openPassword(row)">重置密码</el-button>
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

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑用户' : '新增用户'" width="560px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="用户名">
          <el-input v-model="form.username" :disabled="Boolean(form.id)" />
        </el-form-item>
        <el-form-item v-if="!form.id" label="密码">
          <el-input v-model="form.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="form.nickname" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="form.phone" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="form.email" />
        </el-form-item>
        <el-form-item label="部门">
          <el-tree-select v-model="form.dept_id" :data="deptOptions" node-key="id" check-strictly :props="{ label: 'name', value: 'id' }" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="form.role_ids" multiple style="width: 100%">
            <el-option v-for="role in roles" :key="role.id" :label="role.name" :value="role.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio-button :label="1">启用</el-radio-button>
            <el-radio-button :label="0">禁用</el-radio-button>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  assignUserRoles,
  createUser,
  deleteUser,
  deptTree,
  listRoles,
  listUsers,
  resetUserPassword,
  updateUser,
  updateUserStatus
} from '../api/system'

const loading = ref(false)
const saving = ref(false)
const rows = ref([])
const total = ref(0)
const roles = ref([])
const deptOptions = ref([])
const dialogVisible = ref(false)
const query = reactive({ page: 1, pageSize: 10, keyword: '' })
const form = reactive(emptyForm())

function emptyForm() {
  return { id: 0, username: '', password: '123456', nickname: '', phone: '', email: '', dept_id: 1, status: 1, role_ids: [] }
}

const resetForm = (data = emptyForm()) => {
  Object.assign(form, emptyForm(), data)
}

const load = async () => {
  loading.value = true
  try {
    const res = await listUsers(query)
    rows.value = res.items || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

const loadOptions = async () => {
  const roleRes = await listRoles({ page: 1, pageSize: 100 })
  roles.value = roleRes.items || []
  deptOptions.value = await deptTree()
}

const openCreate = () => {
  resetForm()
  dialogVisible.value = true
}

const openEdit = (row) => {
  resetForm({ ...row, password: '', role_ids: row.role_ids || [] })
  dialogVisible.value = true
}

const save = async () => {
  saving.value = true
  try {
    if (form.id) {
      await updateUser(form.id, form)
      await assignUserRoles(form.id, { role_ids: form.role_ids })
    } else {
      await createUser(form)
    }
    ElMessage.success('保存成功')
    dialogVisible.value = false
    load()
  } finally {
    saving.value = false
  }
}

const changeStatus = async (row, value) => {
  await updateUserStatus(row.id, { status: value ? 1 : 0 })
  row.status = value ? 1 : 0
}

const openPassword = async (row) => {
  const { value } = await ElMessageBox.prompt(`重置 ${row.username} 的密码`, '重置密码', {
    inputValue: '123456',
    inputType: 'password',
    inputPattern: /^.{6,}$/,
    inputErrorMessage: '密码至少 6 位'
  })
  await resetUserPassword(row.id, { password: value })
  ElMessage.success('密码已重置')
}

const remove = async (row) => {
  await ElMessageBox.confirm(`确定删除用户 ${row.username}？`, '删除确认', { type: 'warning' })
  await deleteUser(row.id)
  ElMessage.success('删除成功')
  load()
}

onMounted(() => {
  load()
  loadOptions()
})
</script>

<style scoped>
.pagination {
  justify-content: flex-end;
  margin-top: 14px;
}
</style>
