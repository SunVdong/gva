package mini

import (
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/user"
	"github.com/flipped-aurora/gin-vue-admin/server/service/mini"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
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

// UserInfo 小程序-获取当前用户信息
// @Tags        小程序
// @Summary     获取当前用户信息
// @Description 获取当前登录用户信息（含昵称、头像），需先登录，请求头携带 x-token
// @Accept      json
// @Produce     json
// @Param       x-token header string false "小程序登录后返回的 token"
// @Success     200 {object} response.Response{data=object,msg=string}
// @Router      /mini/user/info [get]
func (a *AuthApi) UserInfo(c *gin.Context) {
	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}
	var u user.User
	if err := global.GVA_DB.Where("id = ?", userID).First(&u).Error; err != nil || u.ID == 0 {
		response.FailWithMessage("用户不存在", c)
		return
	}
	response.OkWithData(miniUserResp(u), c)
}

// UpdateProfile 小程序-修改当前用户昵称与头像
// @Tags        小程序
// @Summary     修改用户昵称与头像
// @Description 支持同时修改昵称和头像，也支持只修改其中一个字段，需先登录，请求头携带 x-token
// @Accept      json
// @Produce     json
// @Param       x-token header string false "小程序登录后返回的 token"
// @Param       data body object true "请求体" example({"nickname":"新的昵称","avatarUrl":"https://xxx/avatar.png"})
// @Success     200 {object} response.Response{data=object,msg=string}
// @Router      /mini/user/updateProfile [post]
func (a *AuthApi) UpdateProfile(c *gin.Context) {
	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}

	var req struct {
		Nickname  *string `json:"nickname"`
		AvatarURL *string `json:"avatarUrl"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数格式错误", c)
		return
	}
	if req.Nickname == nil && req.AvatarURL == nil {
		response.FailWithMessage("nickname 和 avatarUrl 至少传一个", c)
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != nil {
		nickname := strings.TrimSpace(*req.Nickname)
		if nickname == "" {
			response.FailWithMessage("nickname 不能为空", c)
			return
		}
		updates["nickname"] = nickname
	}
	if req.AvatarURL != nil {
		avatarURL := strings.TrimSpace(*req.AvatarURL)
		if avatarURL == "" {
			response.FailWithMessage("avatarUrl 不能为空", c)
			return
		}
		updates["avatar_url"] = avatarURL
	}

	if err := global.GVA_DB.Model(&user.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		response.FailWithMessage("更新失败", c)
		return
	}

	var u user.User
	if err := global.GVA_DB.Where("id = ?", userID).First(&u).Error; err != nil || u.ID == 0 {
		response.FailWithMessage("用户不存在", c)
		return
	}
	response.OkWithDetailed(miniUserResp(u), "更新成功", c)
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
