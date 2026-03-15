package mini

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service/mini"
	"github.com/gin-gonic/gin"
)

type AuthApi struct{}

// Login 小程序登录（微信 code2session + 本系统 JWT）
// @Tags        小程序
// @Summary     小程序登录
// @Description 使用 wx.login 获得的 code 换取 openid，在 users 表创建/绑定用户并签发本系统 JWT，请求头带 x-token 后可选注入 x-user-id
// @Accept      json
// @Produce     json
// @Param       data body object true "请求体" example({"code":"微信 wx.login 返回的 code"})
// @Success     200 {object} response.Response{data=object,msg=string} "data 含 token、user(id,openid,nickname,avatarUrl)"
// @Router      /mini/login [post]
func (a *AuthApi) Login(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请传入 code", c)
		return
	}
	token, user, err := mini.MiniLogin(req.Code)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(gin.H{
		"token": token,
		"user": gin.H{
			"id":        user.ID,
			"openid":    user.OpenID,
			"nickname":  user.Nickname,
			"avatarUrl": user.AvatarURL,
		},
	}, c)
}
