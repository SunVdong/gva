import service from '@/utils/request'

/** 限时活动管理 */
export const createActivity = (data) =>
  service({ url: '/limitedActivity/activity/createActivity', method: 'post', data })
export const deleteActivity = (params) =>
  service({ url: '/limitedActivity/activity/deleteActivity', method: 'delete', params })
export const deleteActivityByIds = (data) =>
  service({ url: '/limitedActivity/activity/deleteActivityByIds', method: 'delete', data })
export const updateActivity = (data) =>
  service({ url: '/limitedActivity/activity/updateActivity', method: 'put', data })
export const findActivity = (params) =>
  service({ url: '/limitedActivity/activity/findActivity', method: 'get', params })
export const getActivityList = (params) =>
  service({ url: '/limitedActivity/activity/getActivityList', method: 'get', params })
