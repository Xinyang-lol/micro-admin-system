import request from '../utils/request'

export const uploadFile = (data) => request.post('/files/upload', data, {
  headers: { 'Content-Type': 'multipart/form-data' }
})
export const listFiles = (params) => request.get('/files', { params })
export const deleteFile = (id) => request.delete(`/files/${id}`)
export const downloadFile = (id) => request.get(`/files/${id}/download`, { responseType: 'blob' })
