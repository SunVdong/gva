package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type FeedbackSearch struct {
	request.PageInfo
	Content string `json:"content" form:"content"`
	UserID  uint   `json:"userId" form:"userId"`
}
