package initialize

import (
	"context"

	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Menu(ctx context.Context) {
	entities := []model.SysBaseMenu{
		{
			ParentId:  0,
			Path:      "limitedActivity",
			Name:      "limitedActivity",
			Hidden:    false,
			Component: "view/routerHolder.vue",
			Sort:      12,
			Meta:      model.Meta{Title: "限时活动", Icon: "present"},
		},
		{
			Path:      "limitedActivityManage",
			Name:      "limitedActivityManage",
			Hidden:    false,
			Component: "plugin/limitedActivity/view/activity.vue",
			Sort:      1,
			Meta:      model.Meta{Title: "活动管理", Icon: "flag"},
		},
		{
			Path:      "limitedActivityOrder",
			Name:      "limitedActivityOrder",
			Hidden:    false,
			Component: "plugin/limitedActivity/view/order.vue",
			Sort:      2,
			Meta:      model.Meta{Title: "活动订单", Icon: "list"},
		},
	}
	utils.RegisterMenus(entities...)
}
