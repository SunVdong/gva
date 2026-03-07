import service from '@/utils/request'

/** 时段 */
export const createTimeSlot = (data) => service({ url: '/camping/timeSlot/createTimeSlot', method: 'post', data })
export const deleteTimeSlot = (params) => service({ url: '/camping/timeSlot/deleteTimeSlot', method: 'delete', params })
export const deleteTimeSlotByIds = (data) => service({ url: '/camping/timeSlot/deleteTimeSlotByIds', method: 'delete', data })
export const updateTimeSlot = (data) => service({ url: '/camping/timeSlot/updateTimeSlot', method: 'put', data })
export const findTimeSlot = (params) => service({ url: '/camping/timeSlot/findTimeSlot', method: 'get', params })
export const getTimeSlotList = (params) => service({ url: '/camping/timeSlot/getTimeSlotList', method: 'get', params })
/** 公开 */
export const getAllTimeSlotsPublic = () => service({ url: '/camping/timeSlot/getAllTimeSlotsPublic', method: 'get' })
