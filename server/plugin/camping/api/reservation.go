package api

import (
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	campingRequest "github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model/request"
	"github.com/gin-gonic/gin"
)

var Reservation = new(reservationApi)

type reservationApi struct{}

// CreateReservation 创建预约（公开接口，前台提交预约）
// @Tags CampingReservation
// @Summary 提交预约(公开)
// @accept application/json
// @Produce application/json
// @Param data body request.CreateCampingReservationRequest true "预约信息"
// @Success 200 {object} response.Response{data=model.CampingReservation,msg=string} "预约成功"
// @Router /camping/reservation/createReservation [post]
func (a *reservationApi) CreateReservation(c *gin.Context) {
	var req campingRequest.CreateCampingReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := serviceResv.CreateReservation(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(res, c)
}

// GetReservation 根据ID查询预约
// @Tags CampingReservation
// @Summary 根据ID查询预约
// @Security ApiKeyAuth
// @Param id query int true "预约ID"
// @Success 200 {object} response.Response{data=model.CampingReservation,msg=string} "查询成功"
// @Router /camping/reservation/getReservation [get]
func (a *reservationApi) GetReservation(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := serviceResv.GetReservation(idReq.ID)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(res, c)
}

// GetReservationByVerifyCodePublic 公开：根据核销码查询预约详情（用于展示二维码信息）
// @Tags CampingReservation
// @Summary 根据核销码查询预约(公开)
// @Param code query string true "核销码"
// @Success 200 {object} response.Response{data=model.CampingReservation,msg=string} "查询成功"
// @Router /camping/reservation/getReservationByVerifyCodePublic [get]
func (a *reservationApi) GetReservationByVerifyCodePublic(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.FailWithMessage("核销码不能为空", c)
		return
	}
	res, err := serviceResv.GetReservationByVerifyCode(code)
	if err != nil {
		response.FailWithMessage("预约不存在或已失效", c)
		return
	}
	response.OkWithData(res, c)
}

// GetReservationList 分页获取预约列表
// @Tags CampingReservation
// @Summary 分页获取预约列表
// @Security ApiKeyAuth
// @Param data query request.CampingReservationSearch true "查询参数"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /camping/reservation/getReservationList [get]
func (a *reservationApi) GetReservationList(c *gin.Context) {
	var req campingRequest.CampingReservationSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceResv.GetReservationList(req)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// VerifyReservation 核销（根据ID）
// @Tags CampingReservation
// @Summary 核销预约
// @Security ApiKeyAuth
// @Param id query int true "预约ID"
// @Success 200 {object} response.Response{msg=string} "核销成功"
// @Router /camping/reservation/verifyReservation [post]
func (a *reservationApi) VerifyReservation(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceResv.VerifyReservation(idReq.ID); err != nil {
		response.FailWithMessage("核销失败，请检查状态", c)
		return
	}
	response.OkWithMessage("核销成功", c)
}

// VerifyReservationByCode 核销（根据核销码）
// @Tags CampingReservation
// @Summary 根据核销码核销
// @Security ApiKeyAuth
// @Param code query string true "核销码"
// @Success 200 {object} response.Response{msg=string} "核销成功"
// @Router /camping/reservation/verifyReservationByCode [post]
func (a *reservationApi) VerifyReservationByCode(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.FailWithMessage("核销码不能为空", c)
		return
	}
	if err := serviceResv.VerifyReservationByCode(code); err != nil {
		response.FailWithMessage("核销失败，请检查核销码或状态", c)
		return
	}
	response.OkWithMessage("核销成功", c)
}

// CancelReservation 取消预约
// @Tags CampingReservation
// @Summary 取消预约
// @Security ApiKeyAuth
// @Param id query int true "预约ID"
// @Success 200 {object} response.Response{msg=string} "取消成功"
// @Router /camping/reservation/cancelReservation [post]
func (a *reservationApi) CancelReservation(c *gin.Context) {
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceResv.CancelReservation(idReq.ID); err != nil {
		response.FailWithMessage("取消失败", c)
		return
	}
	response.OkWithMessage("取消成功", c)
}

// GetReservedSlotIdsPublic 公开：获取某场地某日已预约的时段ID列表
// @Tags CampingReservation
// @Summary 获取已预约时段(公开)
// @Param siteId query int true "场地ID"
// @Param reserveDate query string true "日期 2006-01-02"
// @Success 200 {object} response.Response{data=[]uint,msg=string} "获取成功"
// @Router /camping/reservation/getReservedSlotIdsPublic [get]
func (a *reservationApi) GetReservedSlotIdsPublic(c *gin.Context) {
	siteIDStr := c.Query("siteId")
	reserveDateStr := c.Query("reserveDate")
	if siteIDStr == "" || reserveDateStr == "" {
		response.FailWithMessage("场地ID和预约日期不能为空", c)
		return
	}
	var siteID uint
	if id, err := strconv.ParseUint(siteIDStr, 10, 32); err != nil {
		response.FailWithMessage("场地ID格式错误", c)
		return
	} else {
		siteID = uint(id)
	}
	reserveDate, err := time.ParseInLocation("2006-01-02", reserveDateStr, time.Local)
	if err != nil {
		response.FailWithMessage("日期格式错误，请使用 2006-01-02", c)
		return
	}
	ids, err := serviceResv.GetReservedSlotIds(siteID, reserveDate)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(ids, c)
}
