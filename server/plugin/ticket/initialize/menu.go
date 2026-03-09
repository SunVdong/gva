package initialize

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Menu(ctx context.Context) {
	entities := []model.SysBaseMenu{
		{
			ParentId:  0,
			Path:      "ticket",
			Name:      "ticket",
			Hidden:    false,
			Component: "view/routerHolder.vue",
			Sort:      11,
			Meta:      model.Meta{Title: "景点门票", Icon: "ticket"},
		},
		{
			Path:      "ticketScenic",
			Name:      "ticketScenic",
			Hidden:    false,
			Component: "plugin/ticket/view/scenic.vue",
			Sort:      1,
			Meta:      model.Meta{Title: "景区管理", Icon: "place"},
		},
		{
			Path:      "ticketProduct",
			Name:      "ticketProduct",
			Hidden:    false,
			Component: "plugin/ticket/view/product.vue",
			Sort:      2,
			Meta:      model.Meta{Title: "门票商品", Icon: "goods"},
		},
		{
			Path:      "ticketCalendar",
			Name:      "ticketCalendar",
			Hidden:    false,
			Component: "plugin/ticket/view/calendar.vue",
			Sort:      3,
			Meta:      model.Meta{Title: "日历库存", Icon: "calendar"},
		},
	}
	utils.RegisterMenus(entities...)
	var ticketMenu model.SysBaseMenu
	if err := global.GVA_DB.Where("name = ?", "ticket").First(&ticketMenu).Error; err == nil && ticketMenu.ID > 0 {
		global.GVA_DB.Model(&model.SysBaseMenu{}).
			Where("name IN ?", []string{"ticketScenic", "ticketProduct", "ticketCalendar"}).
			Update("parent_id", ticketMenu.ID)
	}
}
