package mini

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// GenToken 生成测试 JWT（等同 server/t/gen_token，用于调试）
// @Tags        小程序
// @Summary     生成测试 JWT
// @Description 生成一个包含测试用户信息的 JWT，用于本地调试，可直接作为 x-token 使用
// @Produce     json
// @Success     200 {object} response.Response{data=object,msg=string} "data 含 token"
// @Router      /mini/genToken [get]
func (a *CommonApi) GenToken(c *gin.Context) {
	j := utils.NewJWT()
	claims := j.CreateClaims(systemReq.BaseClaims{
		UUID:        uuid.New(),
		ID:          1,
		Username:    "test_user",
		NickName:    "测试用户",
		AuthorityId: 1,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		response.FailWithMessage("生成 token 失败", c)
		return
	}
	response.OkWithData(gin.H{
		"token": token,
	}, c)
}
