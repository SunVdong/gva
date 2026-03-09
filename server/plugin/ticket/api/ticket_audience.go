package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model"
	"github.com/gin-gonic/gin"
)

var Audience = new(ticketAudienceApi)

type ticketAudienceApi struct{}

func (a *ticketAudienceApi) GetBySku(c *gin.Context) {
	var idReq struct {
		SkuID uint `form:"skuId" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage("SKU ID不能为空", c)
		return
	}
	list, err := serviceAudience.GetBySku(idReq.SkuID)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(list, c)
}

func (a *ticketAudienceApi) Save(c *gin.Context) {
	var body struct {
		SkuID uint `json:"skuId" binding:"required"`
		List  []struct {
			AudienceType string `json:"audienceType"`
			Description  string `json:"description"`
		} `json:"list"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list := make([]model.TicketAudience, 0, len(body.List))
	for _, item := range body.List {
		list = append(list, model.TicketAudience{
			SkuID:        body.SkuID,
			AudienceType: item.AudienceType,
			Description:  item.Description,
		})
	}
	if err := serviceAudience.SaveBySku(body.SkuID, list); err != nil {
		response.FailWithMessage("保存失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("保存成功", c)
}
