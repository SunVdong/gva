package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	ticketRequest "github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
	"github.com/gin-gonic/gin"
)

var ScenicOpenTime = new(scenicOpenTimeApi)

type scenicOpenTimeApi struct{}

func (a *scenicOpenTimeApi) GetByScenic(c *gin.Context) {
	var idReq struct {
		ScenicID uint `form:"scenicId" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage("景区ID不能为空", c)
		return
	}
	list, err := serviceScenicOpenTime.GetByScenic(idReq.ScenicID)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(list, c)
}

func (a *scenicOpenTimeApi) Save(c *gin.Context) {
	var req ticketRequest.ScenicOpenTimeSave
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceScenicOpenTime.Save(req); err != nil {
		response.FailWithMessage("保存失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("保存成功", c)
}
