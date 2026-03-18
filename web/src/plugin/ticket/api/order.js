import service from '@/utils/request'

/** 订单管理 */
export const getOrderList = (params) => service({ url: '/ticket/order/getOrderList', method: 'get', params })
export const findOrder = (params) => service({ url: '/ticket/order/findOrder', method: 'get', params })

/** 门票订单 - H5 核销相关公开接口 */
export const getTicketOrderByCodePublic = (params) =>
  service({ url: '/ticket/order/getOrderByCodePublic', method: 'get', params })

export const verifyTicketOrderByCodePublic = (params) =>
  service({ url: '/ticket/order/verifyOrderByCodePublic', method: 'post', params })
