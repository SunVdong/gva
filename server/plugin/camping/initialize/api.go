package initialize

import (
	"context"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Api(ctx context.Context) {
	entities := []model.SysApi{
		{Path: "/camping/site/createSite", Description: "创建露营场地", ApiGroup: "露营", Method: "POST"},
		{Path: "/camping/site/deleteSite", Description: "删除场地", ApiGroup: "露营", Method: "DELETE"},
		{Path: "/camping/site/deleteSiteByIds", Description: "批量删除场地", ApiGroup: "露营", Method: "DELETE"},
		{Path: "/camping/site/updateSite", Description: "更新场地", ApiGroup: "露营", Method: "PUT"},
		{Path: "/camping/site/findSite", Description: "查询场地", ApiGroup: "露营", Method: "GET"},
		{Path: "/camping/site/getSiteList", Description: "场地列表", ApiGroup: "露营", Method: "GET"},
		{Path: "/camping/timeSlot/createTimeSlot", Description: "创建时段", ApiGroup: "露营", Method: "POST"},
		{Path: "/camping/timeSlot/deleteTimeSlot", Description: "删除时段", ApiGroup: "露营", Method: "DELETE"},
		{Path: "/camping/timeSlot/deleteTimeSlotByIds", Description: "批量删除时段", ApiGroup: "露营", Method: "DELETE"},
		{Path: "/camping/timeSlot/updateTimeSlot", Description: "更新时段", ApiGroup: "露营", Method: "PUT"},
		{Path: "/camping/timeSlot/findTimeSlot", Description: "查询时段", ApiGroup: "露营", Method: "GET"},
		{Path: "/camping/timeSlot/getTimeSlotList", Description: "时段列表", ApiGroup: "露营", Method: "GET"},
		{Path: "/camping/reservation/getReservation", Description: "查询预约", ApiGroup: "露营", Method: "GET"},
		{Path: "/camping/reservation/getReservationList", Description: "预约列表", ApiGroup: "露营", Method: "GET"},
		{Path: "/camping/reservation/verifyReservation", Description: "核销预约", ApiGroup: "露营", Method: "POST"},
		{Path: "/camping/reservation/verifyReservationByCode", Description: "按核销码核销", ApiGroup: "露营", Method: "POST"},
		{Path: "/camping/reservation/cancelReservation", Description: "取消预约", ApiGroup: "露营", Method: "POST"},
	}
	utils.RegisterApis(entities...)
}
