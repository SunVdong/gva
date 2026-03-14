package mini

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type CommonApi struct{}

// Ping 小程序健康检查/连通性检测，无需鉴权
func (a *CommonApi) Ping(c *gin.Context) {
	response.OkWithData(map[string]interface{}{
		"message": "pong",
	}, c)
}
