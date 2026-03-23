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

// MiniLogin 小程序登录（仅 openid，不获取手机号）：code2session -> 按 openid 查找/创建用户 -> 签发 JWT
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
			OpenID:   res.OpenID,
			UnionID:  res.UnionID,
			Nickname: "微信用户",
		}
		if err = global.GVA_DB.Create(&u).Error; err != nil {
			return "", u, fmt.Errorf("创建用户失败: %w", err)
		}
	} else if u.UnionID == "" && res.UnionID != "" {
		global.GVA_DB.Model(&u).Update("unionid", res.UnionID)
		u.UnionID = res.UnionID
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

// MiniLoginDecrypt 生产级小程序登录（code + encryptedData + iv）：
// 1. code → jscode2session → openid + session_key（session_key 仅用于解密，不落库）
// 2. session_key + iv 解密 encryptedData → 手机号
// 3. 按 phone 或 openid 查找/创建用户，同时绑定 openid + phone
// 4. 签发 JWT
func MiniLoginDecrypt(code, encryptedData, iv string) (token string, u user.User, err error) {
	cfg := &miniConfig.Config{
		AppID:     global.GVA_CONFIG.Miniprogram.AppID,
		AppSecret: global.GVA_CONFIG.Miniprogram.AppSecret,
		Cache:     cache.NewMemory(),
	}
	if cfg.AppID == "" || cfg.AppSecret == "" {
		return "", u, fmt.Errorf("小程序 AppID/AppSecret 未配置")
	}
	mini := miniprogram.NewMiniProgram(cfg)

	// 1. jscode2session：用 code 换取 openid + session_key
	sess, err := mini.GetAuth().Code2Session(code)
	if err != nil {
		return "", u, fmt.Errorf("微信登录失败: %w", err)
	}
	if sess.ErrCode != 0 {
		return "", u, fmt.Errorf("微信登录失败: %s", sess.ErrMsg)
	}
	if sess.OpenID == "" {
		return "", u, fmt.Errorf("微信未返回 openid")
	}

	// 2. 用 session_key 解密 encryptedData 得到手机号（session_key 不存库）
	plainData, err := mini.GetEncryptor().Decrypt(sess.SessionKey, encryptedData, iv)
	if err != nil {
		return "", u, fmt.Errorf("解密手机号失败: %w", err)
	}
	phone := plainData.PurePhoneNumber
	if phone == "" {
		phone = plainData.PhoneNumber
	}
	if phone == "" {
		return "", u, fmt.Errorf("未获取到手机号")
	}

	// 3. 按 phone 或 openid 查找/创建用户
	err = global.GVA_DB.Where("phone = ? OR openid = ?", phone, sess.OpenID).First(&u).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", u, fmt.Errorf("查询用户失败: %w", err)
		}
		u = user.User{
			OpenID:   sess.OpenID,
			UnionID:  sess.UnionID,
			Nickname: "微信用户",
		}
		u.Phone = &phone
		if err = global.GVA_DB.Create(&u).Error; err != nil {
			return "", u, fmt.Errorf("创建用户失败: %w", err)
		}
	} else {
		updateMap := map[string]interface{}{}
		if u.OpenID == "" {
			updateMap["openid"] = sess.OpenID
			u.OpenID = sess.OpenID
		}
		if u.UnionID == "" && sess.UnionID != "" {
			updateMap["unionid"] = sess.UnionID
			u.UnionID = sess.UnionID
		}
		if u.Phone == nil || *u.Phone == "" {
			updateMap["phone"] = phone
			u.Phone = &phone
		}
		if len(updateMap) > 0 {
			if err = global.GVA_DB.Model(&u).Updates(updateMap).Error; err != nil {
				return "", u, fmt.Errorf("更新用户失败: %w", err)
			}
		}
	}

	// 4. 签发 JWT（手机号作为主标识）
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
