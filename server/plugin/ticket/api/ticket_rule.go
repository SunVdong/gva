package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model"
	ticketRequest "github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
	"github.com/gin-gonic/gin"
)

var Rule = new(ticketRuleApi)

type ticketRuleApi struct{}

func (a *ticketRuleApi) GetByProduct(c *gin.Context) {
	var idReq struct {
		ProductID uint `form:"productId" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage("商品ID不能为空", c)
		return
	}
	list, err := serviceRule.GetByProduct(idReq.ProductID)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(list, c)
}

func (a *ticketRuleApi) Save(c *gin.Context) {
	var body struct {
		ProductID uint                      `json:"productId" binding:"required"`
		List      []ticketRequest.TicketRuleItem `json:"list"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list := make([]model.TicketRule, 0, len(body.List))
	for _, item := range body.List {
		list = append(list, model.TicketRule{
			ProductID: body.ProductID,
			Title:     item.Title,
			Content:   item.Content,
			Sort:      item.Sort,
		})
	}
	if err := serviceRule.SaveByProduct(body.ProductID, list); err != nil {
		response.FailWithMessage("保存失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("保存成功", c)
}
