import service from '@/utils/request'

/** 订单管理 */
export const getOrderList = (params) => service({ url: '/ticket/order/getOrderList', method: 'get', params })
export const findOrder = (params) => service({ url: '/ticket/order/findOrder', method: 'get', params })
/** 核销场合白名单（公开，后端唯一来源） */
export const getVenueOptions = () =>
  service({ url: '/ticket/order/getVenueOptionsPublic', method: 'get' })
/** 核销场合按月统计 */
export const getVenueVerifyStats = (params) =>
  service({ url: '/ticket/order/getVenueVerifyStats', method: 'get', params })

/** 门票订单 - H5 核销相关公开接口 */
export const getTicketOrderByCodePublic = (params) =>
  service({ url: '/ticket/order/getOrderByCodePublic', method: 'get', params })

export const verifyTicketOrderByCodePublic = ({ code, venue } = {}) =>
  service({
    url: '/ticket/order/verifyOrderByCodePublic',
    method: 'post',
    params: { code },
    ...(venue ? { data: { venue } } : {})
  })
