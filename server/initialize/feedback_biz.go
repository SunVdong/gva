package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

// RegisterFeedbackApis 与 activityGuide 等插件一致：用 RegisterApis 将接口写入 sys_apis（FirstOrCreate），
// 已有数据库在每次服务启动后也会在「API 管理 / 角色权限」里看到可勾选项；不依赖仅执行一次的 source/system 初始化种子。
func RegisterFeedbackApis() {
	if global.GVA_DB == nil {
		return
	}
	entities := []system.SysApi{
		{Path: "/feedback/deleteFeedback", Description: "删除意见反馈", ApiGroup: "意见反馈", Method: "DELETE"},
		{Path: "/feedback/deleteFeedbackByIds", Description: "批量删除意见反馈", ApiGroup: "意见反馈", Method: "DELETE"},
		{Path: "/feedback/findFeedback", Description: "根据ID获取意见反馈", ApiGroup: "意见反馈", Method: "GET"},
		{Path: "/feedback/getFeedbackList", Description: "获取意见反馈列表", ApiGroup: "意见反馈", Method: "GET"},
	}
	utils.RegisterApis(entities...)
}
