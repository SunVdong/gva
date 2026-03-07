package initialize

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Menu(ctx context.Context) {
	// 一级菜单：露营预约（parentId=0），下挂 4 个子菜单
	entities := []model.SysBaseMenu{
		// 父菜单：露营预约（一级）
		{
			ParentId:  0,
			Path:      "camping",
			Name:      "camping",
			Hidden:    false,
			Component: "view/routerHolder.vue",
			Sort:      10,
			Meta:      model.Meta{Title: "露营预约", Icon: "office-building"},
		},
		// 子菜单
		{
			Path:      "campingSite",
			Name:      "campingSite",
			Hidden:    false,
			Component: "plugin/camping/view/site.vue",
			Sort:      1,
			Meta:      model.Meta{Title: "露营场地管理", Icon: "office-building"},
		},
		{
			Path:      "campingTimeSlot",
			Name:      "campingTimeSlot",
			Hidden:    false,
			Component: "plugin/camping/view/timeSlot.vue",
			Sort:      2,
			Meta:      model.Meta{Title: "预约时段管理", Icon: "timer"},
		},
		{
			Path:      "campingReservation",
			Name:      "campingReservation",
			Hidden:    false,
			Component: "plugin/camping/view/reservation.vue",
			Sort:      3,
			Meta:      model.Meta{Title: "预约列表", Icon: "list"},
		},
		{
			Path:      "campingVerify",
			Name:      "campingVerify",
			Hidden:    false,
			Component: "plugin/camping/view/verify.vue",
			Sort:      4,
			Meta:      model.Meta{Title: "核销", Icon: "check"},
		},
	}
	utils.RegisterMenus(entities...)
	// 将已存在的 4 个子菜单的父级统一改为「露营预约」（RegisterMenus 里 FirstOrCreate 不会改已有记录的 parent_id）
	var campingMenu model.SysBaseMenu
	if err := global.GVA_DB.Where("name = ?", "camping").First(&campingMenu).Error; err == nil && campingMenu.ID > 0 {
		global.GVA_DB.Model(&model.SysBaseMenu{}).
			Where("name IN ?", []string{"campingSite", "campingTimeSlot", "campingReservation", "campingVerify"}).
			Update("parent_id", campingMenu.ID)
	}
}
