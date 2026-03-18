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

// CreateReservation 创建预约（公开接口，前台提交预约；可选带 x-token 关联当前用户）
// @Tags CampingReservation
// @Summary 提交预约(公开)
// @Param x-token header string false "可选，小程序登录后返回的 token，携带时预约关联当前用户"
// @Param data body request.CreateVenueReservationRequest true "预约信息"
// @Success 200 {object} response.Response{data=model.VenueReservation,msg=string} "预约成功"
// @Router /camping/reservation/createReservation [post]
func (a *reservationApi) CreateReservation(c *gin.Context) {
	var req campingRequest.CreateVenueReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userID := uint(0)
	if uid, exists := c.Get("x-user-id"); exists {
		if u, ok := uid.(uint); ok {
			userID = u
		}
	}
	res, err := serviceResv.CreateReservation(req, userID)
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
// @Success 200 {object} response.Response{data=model.VenueReservation,msg=string} "查询成功"
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

// GetReservationByVerifyCodePublic 公开：根据核销码查询预约详情
// @Tags CampingReservation
// @Summary 根据核销码查询预约(公开)
// @Param code query string true "核销码"
// @Success 200 {object} response.Response{data=model.VenueReservation,msg=string} "查询成功"
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
// @Param data query request.VenueReservationSearch true "查询参数"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /camping/reservation/getReservationList [get]
func (a *reservationApi) GetReservationList(c *gin.Context) {
	var req campingRequest.VenueReservationSearch
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

// VerifyReservationByCodePublic 公开核销（根据核销码，无需登录，用于 H5 扫码核销）
// @Tags CampingReservation
// @Summary 根据核销码核销(公开)
// @Param code query string true "核销码"
// @Success 200 {object} response.Response{data=model.VenueReservation,msg=string} "核销成功"
// @Router /camping/reservation/verifyReservationByCodePublic [post]
func (a *reservationApi) VerifyReservationByCodePublic(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.FailWithMessage("核销码不能为空", c)
		return
	}
	if err := serviceResv.VerifyReservationByCode(code); err != nil {
		response.FailWithMessage("核销失败，请检查核销码或状态", c)
		return
	}
	res, _ := serviceResv.GetReservationByVerifyCode(code)
	response.OkWithData(res, c)
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

// GetReservedSlotIdsPublic 公开：获取某场地某日已约满的时段ID列表
// @Tags CampingReservation
// @Summary 获取已约满时段(公开)
// @Param venueId query int true "场地ID"
// @Param reserveDate query string true "日期 2006-01-02"
// @Success 200 {object} response.Response{data=[]uint,msg=string} "获取成功"
// @Router /camping/reservation/getReservedSlotIdsPublic [get]
func (a *reservationApi) GetReservedSlotIdsPublic(c *gin.Context) {
	venueIDStr := c.Query("venueId")
	reserveDateStr := c.Query("reserveDate")
	if venueIDStr == "" || reserveDateStr == "" {
		response.FailWithMessage("场地ID和预约日期不能为空", c)
		return
	}
	venueID, err := strconv.ParseUint(venueIDStr, 10, 32)
	if err != nil {
		response.FailWithMessage("场地ID格式错误", c)
		return
	}
	reserveDate, err := time.ParseInLocation("2006-01-02", reserveDateStr, time.Local)
	if err != nil {
		response.FailWithMessage("日期格式错误，请使用 2006-01-02", c)
		return
	}
	ids, err := serviceResv.GetReservedTimeslotIds(uint(venueID), reserveDate)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(ids, c)
}
