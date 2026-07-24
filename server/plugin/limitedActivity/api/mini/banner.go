package mini

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

var Banner = new(miniBannerApi)

type miniBannerApi struct{}

// List 小程序-Banner 列表（仅显示中）
// @Tags        小程序-限时活动
// @Summary     Banner 列表
// @Description 返回显示中的 Banner（status=1），含 image 与 detailImage，按 sort 升序
// @Accept      json
// @Produce     json
// @Success     200 {object} response.Response{data=[]object,msg=string}
// @Router      /limitedActivity/mini/banner/list [get]
func (a *miniBannerApi) List(c *gin.Context) {
	list, err := svcBanner.GetMiniList()
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(list, c)
}
