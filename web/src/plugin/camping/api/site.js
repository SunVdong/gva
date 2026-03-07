import service from '@/utils/request'

/** 场地 */
export const createSite = (data) => service({ url: '/camping/site/createSite', method: 'post', data })
export const deleteSite = (params) => service({ url: '/camping/site/deleteSite', method: 'delete', params })
export const deleteSiteByIds = (data) => service({ url: '/camping/site/deleteSiteByIds', method: 'delete', data })
export const updateSite = (data) => service({ url: '/camping/site/updateSite', method: 'put', data })
export const findSite = (params) => service({ url: '/camping/site/findSite', method: 'get', params })
export const getSiteList = (params) => service({ url: '/camping/site/getSiteList', method: 'get', params })
/** 公开 */
export const getSiteListPublic = () => service({ url: '/camping/site/getSiteListPublic', method: 'get' })
export const getSiteDetailPublic = (params) => service({ url: '/camping/site/getSiteDetailPublic', method: 'get', params })
