package mini

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type CommonApi struct{}

// Ping 小程序健康检查/连通性检测，无需鉴权
// @Tags        小程序
// @Summary     健康检查
// @Description 小程序端连通性检测，无需鉴权
// @Accept      json
// @Produce     json
// @Success     200 {object} response.Response{data=object,msg=string}
// @Router      /mini/ping [get]
func (a *CommonApi) Ping(c *gin.Context) {
	response.OkWithData(map[string]interface{}{
		"message": "pong",
	}, c)
}
