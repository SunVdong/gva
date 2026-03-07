package initialize

import (
	"context"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Api(ctx context.Context) {
	entities := []model.SysApi{
		{Path: "/activityGuide/createGuide", Description: "创建活动指南", ApiGroup: "活动指南", Method: "POST"},
		{Path: "/activityGuide/deleteGuide", Description: "删除活动指南", ApiGroup: "活动指南", Method: "DELETE"},
		{Path: "/activityGuide/deleteGuideByIds", Description: "批量删除活动指南", ApiGroup: "活动指南", Method: "DELETE"},
		{Path: "/activityGuide/updateGuide", Description: "更新活动指南", ApiGroup: "活动指南", Method: "PUT"},
		{Path: "/activityGuide/findGuide", Description: "根据ID获取活动指南", ApiGroup: "活动指南", Method: "GET"},
		{Path: "/activityGuide/getGuideList", Description: "获取活动指南列表", ApiGroup: "活动指南", Method: "GET"},
	}
	utils.RegisterApis(entities...)
}
