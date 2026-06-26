import request from '../utils/request'

export const listUsers = (params) => request.get('/users', { params })
export const createUser = (data) => request.post('/users', data)
export const updateUser = (id, data) => request.put(`/users/${id}`, data)
export const deleteUser = (id) => request.delete(`/users/${id}`)
export const updateUserStatus = (id, data) => request.put(`/users/${id}/status`, data)
export const resetUserPassword = (id, data) => request.put(`/users/${id}/password`, data)
export const assignUserRoles = (id, data) => request.put(`/users/${id}/roles`, data)

export const listRoles = (params) => request.get('/roles', { params })
export const createRole = (data) => request.post('/roles', data)
export const updateRole = (id, data) => request.put(`/roles/${id}`, data)
export const deleteRole = (id) => request.delete(`/roles/${id}`)
export const assignRoleMenus = (id, data) => request.put(`/roles/${id}/menus`, data)

export const menuTree = () => request.get('/menus/tree')
export const createMenu = (data) => request.post('/menus', data)
export const updateMenu = (id, data) => request.put(`/menus/${id}`, data)
export const deleteMenu = (id) => request.delete(`/menus/${id}`)

export const deptTree = () => request.get('/depts/tree')
export const createDept = (data) => request.post('/depts', data)
export const updateDept = (id, data) => request.put(`/depts/${id}`, data)
export const deleteDept = (id) => request.delete(`/depts/${id}`)
