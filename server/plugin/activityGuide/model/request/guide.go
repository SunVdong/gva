package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type GuideSearch struct {
	ID             *uint      `json:"ID" form:"ID"`
	Name           string     `json:"name" form:"name"`
	IsPreview      *bool      `json:"isPreview" form:"isPreview"`
	ShowStatus     *bool      `json:"showStatus" form:"showStatus"`
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	request.PageInfo
}
