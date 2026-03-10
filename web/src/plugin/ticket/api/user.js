import service from '@/utils/request'

/** 用户管理（C 端用户） */
export const getUserList = (params) => service({ url: '/ticket/user/getUserList', method: 'get', params })
export const findUser = (params) => service({ url: '/ticket/user/findUser', method: 'get', params })
export const updateUser = (data) => service({ url: '/ticket/user/updateUser', method: 'put', data })
