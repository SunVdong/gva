import service from '@/utils/request'

/** 门票商品 */
export const createProduct = (data) => service({ url: '/ticket/product/createProduct', method: 'post', data })
export const deleteProduct = (params) => service({ url: '/ticket/product/deleteProduct', method: 'delete', params })
export const deleteProductByIds = (data) => service({ url: '/ticket/product/deleteProductByIds', method: 'delete', data })
export const updateProduct = (data) => service({ url: '/ticket/product/updateProduct', method: 'put', data })
export const findProduct = (params) => service({ url: '/ticket/product/findProduct', method: 'get', params })
export const getProductList = (params) => service({ url: '/ticket/product/getProductList', method: 'get', params })

/** 门票 SKU */
export const createSku = (data) => service({ url: '/ticket/sku/createSku', method: 'post', data })
export const deleteSku = (params) => service({ url: '/ticket/sku/deleteSku', method: 'delete', params })
export const deleteSkuByIds = (data) => service({ url: '/ticket/sku/deleteSkuByIds', method: 'delete', data })
export const updateSku = (data) => service({ url: '/ticket/sku/updateSku', method: 'put', data })
export const findSku = (params) => service({ url: '/ticket/sku/findSku', method: 'get', params })
export const getSkuList = (params) => service({ url: '/ticket/sku/getSkuList', method: 'get', params })

/** 门票规则 */
export const getRuleByProduct = (params) => service({ url: '/ticket/rule/getByProduct', method: 'get', params })
export const saveRule = (data) => service({ url: '/ticket/rule/save', method: 'post', data })


/** 日历库存 */
export const getCalendarBySku = (params) => service({ url: '/ticket/calendar/getBySku', method: 'get', params })
export const setCalendar = (data) => service({ url: '/ticket/calendar/set', method: 'post', data })
