<template>
  <div class="page">
    <div class="panel toolbar">
      <el-button @click="load">
        <el-icon><Refresh /></el-icon>
        刷新
      </el-button>
      <el-button type="primary" @click="openCreate({ id: 0 })">
        <el-icon><Plus /></el-icon>
        新增部门
      </el-button>
    </div>
    <div class="panel">
      <el-table :data="rows" row-key="id" default-expand-all>
        <el-table-column prop="name" label="部门名称" min-width="220" />
        <el-table-column prop="sort" label="排序" width="120" />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '禁用' }}</el-tag>
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

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑部门' : '新增部门'" width="520px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="上级部门">
          <el-tree-select v-model="form.parent_id" :data="parentOptions" node-key="id" check-strictly :props="{ label: 'name', value: 'id' }" />
        </el-form-item>
        <el-form-item label="部门名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" /></el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio-button :label="1">启用</el-radio-button>
            <el-radio-button :label="0">禁用</el-radio-button>
          </el-radio-group>
        </el-form-item>
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
import { createDept, deleteDept, deptTree, updateDept } from '../api/system'

const rows = ref([])
const dialogVisible = ref(false)
const form = reactive(emptyForm())
const parentOptions = computed(() => [{ id: 0, name: '根部门', children: rows.value }])

function emptyForm() {
  return { id: 0, parent_id: 0, name: '', sort: 0, status: 1 }
}

const load = async () => {
  rows.value = await deptTree()
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
  if (form.id) await updateDept(form.id, form)
  else await createDept(form)
  ElMessage.success('保存成功')
  dialogVisible.value = false
  load()
}

const remove = async (row) => {
  await ElMessageBox.confirm(`确定删除部门 ${row.name}？`, '删除确认', { type: 'warning' })
  await deleteDept(row.id)
  ElMessage.success('删除成功')
  load()
}

onMounted(load)
</script>
