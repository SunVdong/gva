package mini

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/user"
	"github.com/flipped-aurora/gin-vue-admin/server/service/mini"
	"github.com/gin-gonic/gin"
)

type AuthApi struct{}

// Login 小程序登录（组合登录：wx.login code + getPhoneNumber code）
// @Tags        小程序
// @Summary     小程序组合登录
// @Description 前端先 wx.login 获取 login_code，再在按钮回调中通过 wx.getPhoneNumber 获取 phone_code，后端同时使用两个 code 换取 openid 和手机号，在 users 表中绑定 openid + phone 并签发本系统 JWT，请求头带 x-token 后可选注入 x-user-id
// @Accept      json
// @Produce     json
// @Param       data body object true "请求体" example({"login_code":"wx.login 返回的 code","phone_code":"getPhoneNumber 返回的 code"})
// @Success     200 {object} response.Response{data=object,msg=string} "data 含 token、user(id,openid,nickname,avatarUrl)"
// @Router      /mini/login [post]
func (a *AuthApi) Login(c *gin.Context) {
	var req struct {
		LoginCode string `json:"login_code" binding:"required"`
		PhoneCode string `json:"phone_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请传入 login_code 和 phone_code", c)
		return
	}
	token, user, err := mini.MiniLoginWithPhone(req.LoginCode, req.PhoneCode)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(gin.H{
		"token": token,
		"user":  miniUserResp(user),
	}, c)
}

func miniUserResp(user user.User) gin.H {
	resp := gin.H{
		"id":        user.ID,
		"openid":    user.OpenID,
		"nickname":  user.Nickname,
		"avatarUrl": user.AvatarURL,
	}
	if user.Phone != nil {
		resp["phone"] = *user.Phone
	} else {
		resp["phone"] = ""
	}
	return resp
}
