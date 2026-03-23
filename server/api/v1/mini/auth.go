package mini

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/user"
	"github.com/flipped-aurora/gin-vue-admin/server/service/mini"
	"github.com/gin-gonic/gin"
)

type AuthApi struct{}

// Login 小程序登录（code + encryptedData + iv）
// @Tags        小程序
// @Summary     小程序一键登录
// @Description 前端 getPhoneNumber 拿到 encryptedData/iv 后立即调用 wx.login 获取 code，三者一起发给后端；后端用 code 换 session_key，再用 session_key 解密手机号，完成注册/登录并签发 JWT
// @Accept      json
// @Produce     json
// @Param       data body object true "请求体" example({"code":"wx.login 返回的 code","encryptedData":"加密数据","iv":"初始化向量"})
// @Success     200 {object} response.Response{data=object,msg=string} "data 含 token、user(id,openid,nickname,avatarUrl,phone)"
// @Router      /mini/login [post]
func (a *AuthApi) Login(c *gin.Context) {
	var req struct {
		Code          string `json:"code" binding:"required"`
		EncryptedData string `json:"encryptedData" binding:"required"`
		IV            string `json:"iv" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请传入 code、encryptedData 和 iv", c)
		return
	}
	token, user, err := mini.MiniLoginDecrypt(req.Code, req.EncryptedData, req.IV)
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
