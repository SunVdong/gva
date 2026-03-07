package initialize

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Menu(ctx context.Context) {
	// 删除已有的活动指南菜单，以便用新标题重新创建
	_ = global.GVA_DB.Where("name = ?", "activityGuide").Delete(&model.SysBaseMenu{}).Error

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
