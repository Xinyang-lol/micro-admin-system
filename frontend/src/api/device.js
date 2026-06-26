import request from '../utils/request'

export const listDevices = (params) => request.get('/devices', { params })
export const createDevice = (data) => request.post('/devices', data)
export const updateDevice = (id, data) => request.put(`/devices/${id}`, data)
export const deleteDevice = (id) => request.delete(`/devices/${id}`)
export const deviceStatistics = () => request.get('/devices/statistics')

export const listDeviceTypes = () => request.get('/device-types')
export const createDeviceType = (data) => request.post('/device-types', data)
export const updateDeviceType = (id, data) => request.put(`/device-types/${id}`, data)
export const deleteDeviceType = (id) => request.delete(`/device-types/${id}`)
