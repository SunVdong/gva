import service from '@/utils/request'

/** 场地（Venue） */
export const createSite = (data) => service({ url: '/camping/site/createSite', method: 'post', data })
export const deleteSite = (params) => service({ url: '/camping/site/deleteSite', method: 'delete', params })
export const deleteSiteByIds = (data) => service({ url: '/camping/site/deleteSiteByIds', method: 'delete', data })
export const updateSite = (data) => service({ url: '/camping/site/updateSite', method: 'put', data })
export const findSite = (params) => service({ url: '/camping/site/findSite', method: 'get', params })
export const getSiteList = (params) => service({ url: '/camping/site/getSiteList', method: 'get', params })
/** 公开 */
export const getSiteListPublic = () => service({ url: '/camping/site/getSiteListPublic', method: 'get' })
export const getSiteDetailPublic = (params) => service({ url: '/camping/site/getSiteDetailPublic', method: 'get', params })

/** 场地开放时间 */
export const getVenueOpenTimeByVenue = (params) => service({ url: '/camping/venueOpenTime/getByVenue', method: 'get', params })
export const saveVenueOpenTime = (data) => service({ url: '/camping/venueOpenTime/save', method: 'post', data })

/** 场地日历 */
export const getVenueCalendarByVenue = (params) => service({ url: '/camping/venueCalendar/getByVenue', method: 'get', params })
export const setVenueCalendar = (data) => service({ url: '/camping/venueCalendar/set', method: 'post', data })
