import service from '@/utils/request'

/** 订单管理 */
export const getOrderList = (params) => service({ url: '/ticket/order/getOrderList', method: 'get', params })
export const findOrder = (params) => service({ url: '/ticket/order/findOrder', method: 'get', params })
