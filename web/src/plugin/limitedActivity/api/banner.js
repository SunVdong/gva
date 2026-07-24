import service from '@/utils/request'

/** 限时活动 Banner 管理 */
export const createBanner = (data) =>
  service({ url: '/limitedActivity/banner/createBanner', method: 'post', data })
export const deleteBanner = (params) =>
  service({ url: '/limitedActivity/banner/deleteBanner', method: 'delete', params })
export const deleteBannerByIds = (data) =>
  service({ url: '/limitedActivity/banner/deleteBannerByIds', method: 'delete', data })
export const updateBanner = (data) =>
  service({ url: '/limitedActivity/banner/updateBanner', method: 'put', data })
export const findBanner = (params) =>
  service({ url: '/limitedActivity/banner/findBanner', method: 'get', params })
export const getBannerList = (params) =>
  service({ url: '/limitedActivity/banner/getBannerList', method: 'get', params })
