import service from '@/utils/request'

export const createGuide = (data) => {
  return service({
    url: '/activityGuide/createGuide',
    method: 'post',
    data
  })
}

export const deleteGuide = (params) => {
  return service({
    url: '/activityGuide/deleteGuide',
    method: 'delete',
    params
  })
}

export const deleteGuideByIds = (params) => {
  return service({
    url: '/activityGuide/deleteGuideByIds',
    method: 'delete',
    params
  })
}

export const updateGuide = (data) => {
  return service({
    url: '/activityGuide/updateGuide',
    method: 'put',
    data
  })
}

export const findGuide = (params) => {
  return service({
    url: '/activityGuide/findGuide',
    method: 'get',
    params
  })
}

export const getGuideList = (params) => {
  return service({
    url: '/activityGuide/getGuideList',
    method: 'get',
    params
  })
}
