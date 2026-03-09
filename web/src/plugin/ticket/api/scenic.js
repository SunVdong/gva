import service from '@/utils/request'

/** 景区 */
export const createScenic = (data) => service({ url: '/ticket/scenic/createScenic', method: 'post', data })
export const deleteScenic = (params) => service({ url: '/ticket/scenic/deleteScenic', method: 'delete', params })
export const deleteScenicByIds = (data) => service({ url: '/ticket/scenic/deleteScenicByIds', method: 'delete', data })
export const updateScenic = (data) => service({ url: '/ticket/scenic/updateScenic', method: 'put', data })
export const findScenic = (params) => service({ url: '/ticket/scenic/findScenic', method: 'get', params })
export const getScenicList = (params) => service({ url: '/ticket/scenic/getScenicList', method: 'get', params })

/** 景区开放时间 */
export const getScenicOpenTimeByScenic = (params) => service({ url: '/ticket/scenicOpenTime/getByScenic', method: 'get', params })
export const saveScenicOpenTime = (data) => service({ url: '/ticket/scenicOpenTime/save', method: 'post', data })
