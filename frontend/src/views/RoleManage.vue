<template>
  <div class="page">
    <div class="panel">
      <div class="toolbar">
        <div class="toolbar-left">
          <el-input v-model="query.keyword" clearable placeholder="搜索角色名称/编码" style="width: 240px" @keyup.enter="load" />
          <el-button @click="load">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
        </div>
        <el-button type="primary" @click="openCreate">
          <el-icon><Plus /></el-icon>
          新增角色
        </el-button>
      </div>
    </div>
    <div class="panel">
      <el-table :data="rows" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="角色名称" min-width="140" />
        <el-table-column prop="code" label="角色编码" min-width="140" />
        <el-table-column prop="remark" label="备注" min-width="220" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="230" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button text @click="openMenus(row)">授权</el-button>
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

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑角色' : '新增角色'" width="520px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="角色名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="角色编码"><el-input v-model="form.code" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="form.remark" /></el-form-item>
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

    <el-dialog v-model="menuDialogVisible" title="分配菜单权限" width="520px">
      <el-tree ref="menuTreeRef" :data="menus" node-key="id" show-checkbox default-expand-all :props="{ label: 'name', children: 'children' }" />
      <template #footer>
        <el-button @click="menuDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveMenus">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { nextTick, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { assignRoleMenus, createRole, deleteRole, listRoles, menuTree, updateRole } from '../api/system'

const loading = ref(false)
const saving = ref(false)
const rows = ref([])
const total = ref(0)
const menus = ref([])
const dialogVisible = ref(false)
const menuDialogVisible = ref(false)
const menuTreeRef = ref()
const currentRole = ref(null)
const query = reactive({ page: 1, pageSize: 10, keyword: '' })
const form = reactive({ id: 0, name: '', code: '', remark: '', status: 1 })

const load = async () => {
  loading.value = true
  try {
    const res = await listRoles(query)
    rows.value = res.items || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

const loadMenus = async () => {
  menus.value = await menuTree()
}

const openCreate = () => {
  Object.assign(form, { id: 0, name: '', code: '', remark: '', status: 1 })
  dialogVisible.value = true
}

const openEdit = (row) => {
  Object.assign(form, row)
  dialogVisible.value = true
}

const save = async () => {
  saving.value = true
  try {
    if (form.id) await updateRole(form.id, form)
    else await createRole(form)
    ElMessage.success('保存成功')
    dialogVisible.value = false
    load()
  } finally {
    saving.value = false
  }
}

const openMenus = async (row) => {
  currentRole.value = row
  menuDialogVisible.value = true
  await nextTick()
  menuTreeRef.value?.setCheckedKeys(row.menu_ids || [])
}

const saveMenus = async () => {
  const checked = menuTreeRef.value.getCheckedKeys()
  const half = menuTreeRef.value.getHalfCheckedKeys()
  await assignRoleMenus(currentRole.value.id, { menu_ids: [...checked, ...half] })
  ElMessage.success('授权成功')
  menuDialogVisible.value = false
  load()
}

const remove = async (row) => {
  await ElMessageBox.confirm(`确定删除角色 ${row.name}？`, '删除确认', { type: 'warning' })
  await deleteRole(row.id)
  ElMessage.success('删除成功')
  load()
}

onMounted(() => {
  load()
  loadMenus()
})
</script>

<style scoped>
.pagination {
  justify-content: flex-end;
  margin-top: 14px;
}
</style>
