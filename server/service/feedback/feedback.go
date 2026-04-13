package feedback

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	feedbackModel "github.com/flipped-aurora/gin-vue-admin/server/model/feedback"
	feedbackReq "github.com/flipped-aurora/gin-vue-admin/server/model/feedback/request"
)

type FeedbackService struct{}

func (s *FeedbackService) CreateFeedback(f feedbackModel.Feedback) (err error) {
	return global.GVA_DB.Create(&f).Error
}

func (s *FeedbackService) DeleteFeedback(f feedbackModel.Feedback) (err error) {
	return global.GVA_DB.Delete(&f).Error
}

func (s *FeedbackService) DeleteFeedbackByIds(ids request.IdsReq) (err error) {
	return global.GVA_DB.Delete(&[]feedbackModel.Feedback{}, "id in (?)", ids.Ids).Error
}

func (s *FeedbackService) GetFeedback(id uint) (f feedbackModel.Feedback, err error) {
	err = global.GVA_DB.Preload("User").Where("id = ?", id).First(&f).Error
	return
}

func (s *FeedbackService) GetFeedbackList(info feedbackReq.FeedbackSearch) (list []feedbackModel.Feedback, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&feedbackModel.Feedback{})
	if info.Content != "" {
		db = db.Where("content LIKE ?", "%"+info.Content+"%")
	}
	if info.UserID != 0 {
		db = db.Where("user_id = ?", info.UserID)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("id desc").Preload("User").Find(&list).Error
	return
}

// ListFeedbackByUserID 小程序端：仅查询指定用户的反馈记录（分页）
func (s *FeedbackService) ListFeedbackByUserID(userID uint, pageInfo request.PageInfo) (list []feedbackModel.Feedback, total int64, err error) {
	db := global.GVA_DB.Model(&feedbackModel.Feedback{}).Where("user_id = ?", userID)
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Scopes(pageInfo.Paginate()).Order("id desc").Find(&list).Error
	return
}
