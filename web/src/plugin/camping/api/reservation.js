import service from '@/utils/request'

/** 预约 */
export const getReservation = (params) => service({ url: '/camping/reservation/getReservation', method: 'get', params })
export const getReservationList = (params) => service({ url: '/camping/reservation/getReservationList', method: 'get', params })
export const verifyReservation = (params) => service({ url: '/camping/reservation/verifyReservation', method: 'post', params })
export const verifyReservationByCode = (params) => service({ url: '/camping/reservation/verifyReservationByCode', method: 'post', params })
export const cancelReservation = (params) => service({ url: '/camping/reservation/cancelReservation', method: 'post', params })
/** 公开（H5 核销等） */
export const getReservationByVerifyCodePublic = (params) => service({ url: '/camping/reservation/getReservationByVerifyCodePublic', method: 'get', params })
export const verifyReservationByCodePublic = (params) => service({ url: '/camping/reservation/verifyReservationByCodePublic', method: 'post', params })
export const getReservedSlotIdsPublic = (params) => service({ url: '/camping/reservation/getReservedSlotIdsPublic', method: 'get', params })
