package mini

import (
	"fmt"
	"sort"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
	"github.com/gin-gonic/gin"
)

// 星期 1-7 对应中文
var weekDayNames = []string{"", "周一", "周二", "周三", "周四", "周五", "周六", "周日"}

// buildOpenTimeRemark 根据开放时间列表生成说明：全配且一致则「全年开放 每天 xx:xx-xx:xx」，否则按时间分组「周一，周二 开放时间 xx:xx-xx:xx，…」
func buildOpenTimeRemark(openTimes []model.ScenicOpenTime) string {
	if len(openTimes) == 0 {
		return ""
	}
	toHHMM := func(t model.TimeOnly) string {
		s := string(t)
		if len(s) > 5 {
			return s[:5]
		}
		return s
	}
	// 按 (open, close) 分组，记录该时间段对应的星期几
	type key struct{ open, close string }
	group := make(map[key][]int)
	for _, o := range openTimes {
		openStr := toHHMM(o.OpenTime)
		closeStr := toHHMM(o.CloseTime)
		if openStr == "" || closeStr == "" {
			continue
		}
		k := key{openStr, closeStr}
		group[k] = append(group[k], o.WeekDay)
	}
	if len(group) == 0 {
		return ""
	}
	// 仅一组且覆盖 1-7 且时间一致 → 全年开放 每天 xx:xx-xx:xx
	if len(group) == 1 {
		for k, days := range group {
			seen := make(map[int]bool)
			for _, d := range days {
				if d >= 1 && d <= 7 {
					seen[d] = true
				}
			}
			if len(seen) == 7 {
				return fmt.Sprintf("全年开放 每天 %s-%s", k.open, k.close)
			}
			break
		}
	}
	// 多组或未全配：按时间分组输出「周一，周二 开放时间 xx:xx-xx:xx，…」
	type part struct {
		minDay int
		text   string
	}
	var parts []part
	for k, days := range group {
		sort.Ints(days)
		var names []string
		for _, d := range days {
			if d >= 1 && d <= 7 {
				names = append(names, weekDayNames[d])
			}
		}
		if len(names) == 0 {
			continue
		}
		parts = append(parts, part{
			minDay: days[0],
			text:   fmt.Sprintf("%s 开放时间 %s-%s", strings.Join(names, "，"), k.open, k.close),
		})
	}
	sort.Slice(parts, func(i, j int) bool { return parts[i].minDay < parts[j].minDay })
	var texts []string
	for _, p := range parts {
		texts = append(texts, p.text)
	}
	return strings.Join(texts, "\n")
}

var Scenic = new(miniScenicApi)

type miniScenicApi struct{}

// List 小程序-景区列表（仅启用）
// @Tags        小程序-景点
// @Summary     景区列表
// @Description 小程序端获取已启用的景区列表，分页
// @Accept      json
// @Produce     json
// @Param       page     query    int false "页码"
// @Param       pageSize query    int false "每页条数"
// @Success     200      {object} response.Response{data=response.PageResult,msg=string}
// @Router      /ticket/mini/scenic/list [get]
func (a *miniScenicApi) List(c *gin.Context) {
	status := 1
	req := request.ScenicSearch{Status: &status}
	req.Page = 1
	req.PageSize = 20
	_ = c.ShouldBindQuery(&req)
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}
	list, total, err := svcScenic.GetList(req)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: req.Page, PageSize: req.PageSize}, "获取成功", c)
}

// Detail 小程序-景区详情（仅启用时返回）
// @Tags        小程序-景点
// @Summary     景区详情
// @Description 小程序端获取景区详情，仅启用时返回
// @Accept      json
// @Produce     json
// @Param       id query int true "景区ID"
// @Success     200 {object} response.Response{data=object,msg=string}
// @Router      /ticket/mini/scenic/detail [get]
func (a *miniScenicApi) Detail(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	res, err := svcScenic.Get(idReq.ID)
	if err != nil {
		response.FailWithMessage("景区不存在", c)
		return
	}
	if res.Status != 1 {
		response.FailWithMessage("景区已下架", c)
		return
	}
	openTimes, _ := svcOpenTime.GetByScenic(idReq.ID)
	response.OkWithData(gin.H{
		"scenic":         res,
		"openTimes":      openTimes,
		"openTimeRemark": buildOpenTimeRemark(openTimes),
	}, c)
}
