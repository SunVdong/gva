import service from '@/utils/request'

export const deleteFeedback = (data) => {
  return service({
    url: '/feedback/deleteFeedback',
    method: 'delete',
    data
  })
}

export const deleteFeedbackByIds = (data) => {
  return service({
    url: '/feedback/deleteFeedbackByIds',
    method: 'delete',
    data
  })
}

export const getFeedbackList = (params) => {
  return service({
    url: '/feedback/getFeedbackList',
    method: 'get',
    params
  })
}

export const findFeedback = (params) => {
  return service({
    url: '/feedback/findFeedback',
    method: 'get',
    params
  })
}
