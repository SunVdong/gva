import service from '@/utils/request'

/** 限时活动订单管理 */
export const getOrderList = (params) =>
  service({ url: '/limitedActivity/order/getOrderList', method: 'get', params })
export const findOrder = (params) =>
  service({ url: '/limitedActivity/order/findOrder', method: 'get', params })
export const refundOrder = (params) =>
  service({ url: '/limitedActivity/order/refundOrder', method: 'post', params })

/** H5 核销公开接口 */
export const getActivityOrderByCodePublic = (params) =>
  service({ url: '/limitedActivity/order/getOrderByCodePublic', method: 'get', params })
export const verifyActivityOrderByCodePublic = (params) =>
  service({ url: '/limitedActivity/order/verifyOrderByCodePublic', method: 'post', params })
