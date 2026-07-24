package initialize

import (
	"context"

	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Api(ctx context.Context) {
	entities := []model.SysApi{
		{Path: "/limitedActivity/activity/createActivity", Description: "创建限时活动", ApiGroup: "限时活动", Method: "POST"},
		{Path: "/limitedActivity/activity/deleteActivity", Description: "删除限时活动", ApiGroup: "限时活动", Method: "DELETE"},
		{Path: "/limitedActivity/activity/deleteActivityByIds", Description: "批量删除限时活动", ApiGroup: "限时活动", Method: "DELETE"},
		{Path: "/limitedActivity/activity/updateActivity", Description: "更新限时活动", ApiGroup: "限时活动", Method: "PUT"},
		{Path: "/limitedActivity/activity/findActivity", Description: "查询限时活动", ApiGroup: "限时活动", Method: "GET"},
		{Path: "/limitedActivity/activity/getActivityList", Description: "限时活动列表", ApiGroup: "限时活动", Method: "GET"},
		{Path: "/limitedActivity/order/getOrderList", Description: "活动订单列表", ApiGroup: "限时活动", Method: "GET"},
		{Path: "/limitedActivity/order/findOrder", Description: "活动订单详情", ApiGroup: "限时活动", Method: "GET"},
		{Path: "/limitedActivity/order/refundOrder", Description: "活动订单退款", ApiGroup: "限时活动", Method: "POST"},
		{Path: "/limitedActivity/banner/createBanner", Description: "创建Banner", ApiGroup: "限时活动", Method: "POST"},
		{Path: "/limitedActivity/banner/deleteBanner", Description: "删除Banner", ApiGroup: "限时活动", Method: "DELETE"},
		{Path: "/limitedActivity/banner/deleteBannerByIds", Description: "批量删除Banner", ApiGroup: "限时活动", Method: "DELETE"},
		{Path: "/limitedActivity/banner/updateBanner", Description: "更新Banner", ApiGroup: "限时活动", Method: "PUT"},
		{Path: "/limitedActivity/banner/findBanner", Description: "查询Banner", ApiGroup: "限时活动", Method: "GET"},
		{Path: "/limitedActivity/banner/getBannerList", Description: "Banner列表", ApiGroup: "限时活动", Method: "GET"},
	}
	utils.RegisterApis(entities...)
}
