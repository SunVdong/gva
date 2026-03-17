package mini

import (
	"errors"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/user"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/google/uuid"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"gorm.io/gorm"
)

// MiniLogin 小程序登录：code2session -> 在 users 表中按 openid 查找/创建用户 -> 签发 JWT（与后台 sys_users 分离）
func MiniLogin(code string) (token string, u user.User, err error) {
	cfg := &miniConfig.Config{
		AppID:     global.GVA_CONFIG.Miniprogram.AppID,
		AppSecret: global.GVA_CONFIG.Miniprogram.AppSecret,
		Cache:     cache.NewMemory(),
	}
	if cfg.AppID == "" || cfg.AppSecret == "" {
		return "", u, fmt.Errorf("小程序 AppID/AppSecret 未配置")
	}
	mini := miniprogram.NewMiniProgram(cfg)
	res, err := mini.GetAuth().Code2Session(code)
	if err != nil {
		return "", u, fmt.Errorf("微信登录失败: %w", err)
	}
	if res.ErrCode != 0 {
		return "", u, fmt.Errorf("微信登录失败: %s", res.ErrMsg)
	}
	if res.OpenID == "" {
		return "", u, fmt.Errorf("微信未返回 openid")
	}

	err = global.GVA_DB.Where("openid = ?", res.OpenID).First(&u).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", u, fmt.Errorf("查询用户失败: %w", err)
		}
		u = user.User{
			OpenID:     res.OpenID,
			UnionID:    res.UnionID,
			SessionKey: res.SessionKey,
			Nickname:   "微信用户",
		}
		if err = global.GVA_DB.Create(&u).Error; err != nil {
			return "", u, fmt.Errorf("创建用户失败: %w", err)
		}
	} else {
		u.SessionKey = res.SessionKey
		u.UnionID = res.UnionID
		global.GVA_DB.Model(&u).Updates(map[string]interface{}{
			"session_key": res.SessionKey,
			"unionid":     res.UnionID,
		})
	}

	j := utils.NewJWT()
	claims := j.CreateClaims(systemReq.BaseClaims{
		UUID:        uuid.NewSHA1(uuid.NameSpaceOID, []byte(u.OpenID)),
		ID:          u.ID,
		Username:    u.OpenID,
		NickName:    u.Nickname,
		AuthorityId: 0,
	})
	token, err = j.CreateToken(claims)
	if err != nil {
		return "", u, fmt.Errorf("签发 token 失败: %w", err)
	}
	return token, u, nil
}

// MiniLoginByPhone 本机号一键登录：getPhoneNumber 的 code 换手机号 -> 在 users 表按手机号查找或创建 -> 签发 JWT
func MiniLoginByPhone(code string) (token string, u user.User, err error) {
	cfg := &miniConfig.Config{
		AppID:     global.GVA_CONFIG.Miniprogram.AppID,
		AppSecret: global.GVA_CONFIG.Miniprogram.AppSecret,
		Cache:     cache.NewMemory(),
	}
	if cfg.AppID == "" || cfg.AppSecret == "" {
		return "", u, fmt.Errorf("小程序 AppID/AppSecret 未配置")
	}
	mini := miniprogram.NewMiniProgram(cfg)
	res, err := mini.GetAuth().GetPhoneNumber(code)
	if err != nil {
		return "", u, fmt.Errorf("获取手机号失败: %w", err)
	}
	if res.ErrCode != 0 {
		return "", u, fmt.Errorf("获取手机号失败: %s", res.ErrMsg)
	}
	phone := res.PhoneInfo.PurePhoneNumber
	if phone == "" {
		phone = res.PhoneInfo.PhoneNumber
	}
	if phone == "" {
		return "", u, fmt.Errorf("未获取到手机号")
	}

	err = global.GVA_DB.Where("phone = ?", phone).First(&u).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", u, fmt.Errorf("查询用户失败: %w", err)
		}
		u = user.User{
			Phone:    &phone,
			Nickname: "微信用户",
		}
		if err = global.GVA_DB.Create(&u).Error; err != nil {
			return "", u, fmt.Errorf("创建用户失败: %w", err)
		}
	}

	j := utils.NewJWT()
	claims := j.CreateClaims(systemReq.BaseClaims{
		UUID:        uuid.NewSHA1(uuid.NameSpaceOID, []byte(phone)),
		ID:          u.ID,
		Username:    phone,
		NickName:    u.Nickname,
		AuthorityId: 0,
	})
	token, err = j.CreateToken(claims)
	if err != nil {
		return "", u, fmt.Errorf("签发 token 失败: %w", err)
	}
	return token, u, nil
}

// MiniLoginWithPhone 组合登录：同时使用 wx.login 的 code 和 getPhoneNumber 的 code
// 1. loginCode -> Code2Session，获取 openid/unionid/session_key
// 2. phoneCode -> GetPhoneNumber，获取手机号
// 3. 在 users 表中按手机号或 openid 查找/创建用户，并同时绑定 openid + phone
// 4. 签发 JWT
func MiniLoginWithPhone(loginCode, phoneCode string) (token string, u user.User, err error) {
	cfg := &miniConfig.Config{
		AppID:     global.GVA_CONFIG.Miniprogram.AppID,
		AppSecret: global.GVA_CONFIG.Miniprogram.AppSecret,
		Cache:     cache.NewMemory(),
	}
	if cfg.AppID == "" || cfg.AppSecret == "" {
		return "", u, fmt.Errorf("小程序 AppID/AppSecret 未配置")
	}
	mini := miniprogram.NewMiniProgram(cfg)
	auth := mini.GetAuth()

	// 1. Code2Session
	sess, err := auth.Code2Session(loginCode)
	if err != nil {
		return "", u, fmt.Errorf("微信登录失败: %w", err)
	}
	if sess.ErrCode != 0 {
		return "", u, fmt.Errorf("微信登录失败: %s", sess.ErrMsg)
	}
	if sess.OpenID == "" {
		return "", u, fmt.Errorf("微信未返回 openid")
	}

	// 2. GetPhoneNumber
	phoneRes, err := auth.GetPhoneNumber(phoneCode)
	if err != nil {
		return "", u, fmt.Errorf("获取手机号失败: %w", err)
	}
	if phoneRes.ErrCode != 0 {
		return "", u, fmt.Errorf("获取手机号失败: %s", phoneRes.ErrMsg)
	}
	phone := phoneRes.PhoneInfo.PurePhoneNumber
	if phone == "" {
		phone = phoneRes.PhoneInfo.PhoneNumber
	}
	if phone == "" {
		return "", u, fmt.Errorf("未获取到手机号")
	}

	// 3. 查找或创建用户（按手机号或 openid）
	err = global.GVA_DB.Where("phone = ? OR openid = ?", phone, sess.OpenID).First(&u).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", u, fmt.Errorf("查询用户失败: %w", err)
		}
		// 不存在则创建
		u = user.User{
			OpenID:     sess.OpenID,
			UnionID:    sess.UnionID,
			SessionKey: sess.SessionKey,
			Nickname:   "微信用户",
		}
		u.Phone = &phone
		if err = global.GVA_DB.Create(&u).Error; err != nil {
			return "", u, fmt.Errorf("创建用户失败: %w", err)
		}
	} else {
		// 已存在则补齐 openid/phone/unionid/session_key
		updateMap := map[string]interface{}{
			"session_key": sess.SessionKey,
			"unionid":     sess.UnionID,
		}
		if u.OpenID == "" {
			updateMap["openid"] = sess.OpenID
			u.OpenID = sess.OpenID
		}
		if u.Phone == nil || *u.Phone == "" {
			updateMap["phone"] = phone
			u.Phone = &phone
		}
		u.SessionKey = sess.SessionKey
		u.UnionID = sess.UnionID
		if err = global.GVA_DB.Model(&u).Updates(updateMap).Error; err != nil {
			return "", u, fmt.Errorf("更新用户失败: %w", err)
		}
	}

	// 4. 签发 JWT（优先使用手机号作为 Username/UUID 基础）
	idKey := phone
	if idKey == "" {
		idKey = u.OpenID
	}
	j := utils.NewJWT()
	claims := j.CreateClaims(systemReq.BaseClaims{
		UUID:        uuid.NewSHA1(uuid.NameSpaceOID, []byte(idKey)),
		ID:          u.ID,
		Username:    idKey,
		NickName:    u.Nickname,
		AuthorityId: 0,
	})
	token, err = j.CreateToken(claims)
	if err != nil {
		return "", u, fmt.Errorf("签发 token 失败: %w", err)
	}
	return token, u, nil
}
