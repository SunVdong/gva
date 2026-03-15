package mini

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/user"
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
		"user":  miniUserResp(user),
	}, c)
}

// LoginByPhone 本机号一键登录（getPhoneNumber 返回的 code 换手机号，按手机号登录/注册）
// @Tags        小程序
// @Summary     本机号一键登录
// @Description 前端 wx.getPhoneNumber 用户同意后拿到 code，后端向微信换手机号，在 users 表按手机号查找或创建用户并签发 JWT
// @Accept      json
// @Produce     json
// @Param       data body object true "请求体" example({"code":"getPhoneNumber 返回的 code"})
// @Success     200 {object} response.Response{data=object,msg=string} "data 含 token、user(id,phone,nickname,avatarUrl)"
// @Router      /mini/loginByPhone [post]
func (a *AuthApi) LoginByPhone(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请传入 code", c)
		return
	}
	token, user, err := mini.MiniLoginByPhone(req.Code)
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
