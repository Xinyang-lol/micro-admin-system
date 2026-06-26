<template>
  <div class="page">
    <div class="panel toolbar">
      <div class="toolbar-left">
        <el-button @click="load">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
      <el-button type="primary" @click="openCreate({ id: 0 })">
        <el-icon><Plus /></el-icon>
        新增菜单
      </el-button>
    </div>
    <div class="panel">
      <el-table :data="rows" row-key="id" default-expand-all>
        <el-table-column prop="name" label="名称" min-width="170" />
        <el-table-column prop="path" label="路由" min-width="160" />
        <el-table-column prop="permission" label="权限标识" min-width="190" />
        <el-table-column prop="component" label="组件" min-width="130" />
        <el-table-column prop="sort" label="排序" width="90" />
        <el-table-column label="类型" width="90">
          <template #default="{ row }">
            <el-tag :type="row.type === 0 ? 'primary' : 'success'">{{ row.type === 0 ? '菜单' : '按钮' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" @click="openCreate(row)">新增</el-button>
            <el-button text @click="openEdit(row)">编辑</el-button>
            <el-button text type="danger" @click="remove(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑菜单' : '新增菜单'" width="560px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="上级菜单">
          <el-tree-select v-model="form.parent_id" :data="parentOptions" node-key="id" check-strictly :props="{ label: 'name', value: 'id' }" />
        </el-form-item>
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="路由"><el-input v-model="form.path" /></el-form-item>
        <el-form-item label="组件"><el-input v-model="form.component" /></el-form-item>
        <el-form-item label="权限标识"><el-input v-model="form.permission" /></el-form-item>
        <el-form-item label="图标"><el-input v-model="form.icon" /></el-form-item>
        <el-form-item label="类型">
          <el-radio-group v-model="form.type">
            <el-radio-button :label="0">菜单</el-radio-button>
            <el-radio-button :label="1">按钮</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { createMenu, deleteMenu, menuTree, updateMenu } from '../api/system'

const rows = ref([])
const dialogVisible = ref(false)
const form = reactive(emptyForm())
const parentOptions = computed(() => [{ id: 0, name: '根目录', children: rows.value }])

function emptyForm() {
  return { id: 0, parent_id: 0, name: '', path: '', component: '', permission: '', icon: '', type: 0, sort: 0 }
}

const load = async () => {
  rows.value = await menuTree()
}

const openCreate = (parent) => {
  Object.assign(form, emptyForm(), { parent_id: parent.id || 0 })
  dialogVisible.value = true
}

const openEdit = (row) => {
  Object.assign(form, emptyForm(), row)
  dialogVisible.value = true
}

const save = async () => {
  if (form.id) await updateMenu(form.id, form)
  else await createMenu(form)
  ElMessage.success('保存成功')
  dialogVisible.value = false
  load()
}

const remove = async (row) => {
  await ElMessageBox.confirm(`确定删除菜单 ${row.name}？`, '删除确认', { type: 'warning' })
  await deleteMenu(row.id)
  ElMessage.success('删除成功')
  load()
}

onMounted(load)
</script>
