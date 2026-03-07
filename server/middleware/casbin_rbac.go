package middleware

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

// CasbinHandler 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		waitUse, _ := utils.GetClaims(c)
		path := c.Request.URL.Path
		// 去掉路由前缀后再做权限匹配（库里存的 path 多为 /user/xxx，不含 /api）
		prefix := strings.Trim(global.GVA_CONFIG.System.RouterPrefix, "/")
		if prefix != "" && strings.HasPrefix(path, "/") {
			path = strings.TrimPrefix(path, "/"+prefix)
		}
		obj := path
		if obj == "" {
			obj = "/"
		}
		act := c.Request.Method
		sub := strconv.Itoa(int(waitUse.AuthorityId))
		e := utils.GetCasbin()
		success, _ := e.Enforce(sub, obj, act)
		if !success {
			response.FailWithDetailed(gin.H{}, "权限不足", c)
			c.Abort()
			return
		}
		c.Next()
	}
}
