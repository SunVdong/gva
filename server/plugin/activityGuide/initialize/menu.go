package initialize

import (
	"context"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Menu(ctx context.Context) {
	// 与 camping、announcement 等插件一致：不删除已有菜单，使用 FirstOrCreate 保留菜单 ID，避免角色-菜单关联失效
	entities := []model.SysBaseMenu{
		{
			ParentId:  0,
			MenuLevel: 0,
			Path:      "activityGuide",
			Name:      "activityGuide",
			Hidden:    false,
			Component: "plugin/activityGuide/view/guide.vue",
			Sort:      6,
			Meta:      model.Meta{Title: "活动指南", Icon: "document"},
		},
	}
	utils.RegisterMenus(entities...)
}
