package mini

import (
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model"
	campingRequest "github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model/request"
	"github.com/gin-gonic/gin"
)

type reservationApi struct{}

// Create 小程序-提交预约（需登录，请求头带 x-token）
// @Tags        小程序-露营
// @Summary     提交预约
// @Description 提交露营场地预约（预定人、手机号、日期、时段等），需先登录，请求头携带 x-token
// @Accept      json
// @Produce     json
// @Param       x-token header string false "小程序登录后返回的 token"
// @Param       data body request.CreateVenueReservationRequest true "预约信息"
// @Success     200 {object} response.Response{data=model.VenueReservation,msg=string}
// @Router      /camping/mini/reservation/create [post]
func (a *reservationApi) Create(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.FailWithMessage("请先登录", c)
		return
	}
	var req campingRequest.CreateVenueReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := svcReservation.CreateReservation(req, userID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(res, c)
}

// Update 小程序-修改预约（仅本人）
// @Tags        小程序-露营
// @Summary     修改预约
// @Description 用户修改自己的预约信息（日期、时段、联系人信息等）
// @Accept      json
// @Produce     json
// @Param       x-token header string false "小程序登录后返回的 token"
// @Param       data body request.UpdateVenueReservationRequest true "修改后的预约信息"
// @Success     200 {object} response.Response{data=model.VenueReservation,msg=string}
// @Router      /camping/mini/reservation/update [post]
func (a *reservationApi) Update(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.FailWithMessage("请先登录", c)
		return
	}
	var req campingRequest.UpdateVenueReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := svcReservation.UpdateReservation(req, userID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(res, c)
}

// MyList 小程序-我的预约列表
// @Tags        小程序-露营
// @Summary     我的预约列表
// @Description 按当前用户获取预约列表，分页。需先登录，请求头携带 x-token
// @Accept      json
// @Produce     json
// @Param       x-token  header string false "小程序登录后返回的 token"
// @Param       page     query int false "页码"
// @Param       pageSize query int false "每页条数"
// @Param       status   query int false "状态 0待核销 1已核销 2已取消 3已过期"
// @Success     200 {object} response.Response{data=response.PageResult,msg=string}
// @Router      /camping/mini/reservation/myList [get]
func (a *reservationApi) MyList(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.FailWithMessage("请先登录", c)
		return
	}
	var req campingRequest.VenueReservationSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = &userID
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 50 {
		req.PageSize = 20
	}
	list, total, err := svcReservation.GetReservationList(req)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	// 为每条预约附带场地名、时段展示（便于列表展示）
	items := make([]gin.H, 0, len(list))
	for _, r := range list {
		venueName := ""
		timeslotRange := ""
		if v, _ := svcVenue.GetVenue(r.VenueID); v.ID != 0 {
			venueName = v.Name
		}
		if s, _ := svcVenueTimeslot.GetVenueTimeslot(r.TimeslotID); s.ID != 0 {
			timeslotRange = s.StartTime.FormatHHMM() + "-" + s.EndTime.FormatHHMM()
		}
		items = append(items, gin.H{
			"id":             r.ID,
			"reservationNo":  r.ReservationNo,
			"venueId":        r.VenueID,
			"venueName":      venueName,
			"timeslotId":     r.TimeslotID,
			"timeslotRange":  timeslotRange,
			"reserveDate":    r.ReserveDate,
			"contactName":    r.ContactName,
			"contactPhone":   r.ContactPhone,
			"contactCount":   r.ContactCount,
			"status":         r.Status,
			"verifyCode":     r.VerifyCode,
			"createdAt":      r.CreatedAt,
		})
	}
	response.OkWithDetailed(response.PageResult{
		List: items, Total: total, Page: req.Page, PageSize: req.PageSize,
	}, "获取成功", c)
}

// MyDetail 小程序-预约详情（含核销码，仅本人）
// @Tags        小程序-露营
// @Summary     预约详情
// @Description 获取预约详情含场地名、时段、核销码，仅限本人。需先登录，请求头携带 x-token
// @Accept      json
// @Produce     json
// @Param       x-token header string false "小程序登录后返回的 token"
// @Param       id query int true "预约ID"
// @Success     200 {object} response.Response{data=object,msg=string}
// @Router      /camping/mini/reservation/myDetail [get]
func (a *reservationApi) MyDetail(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.FailWithMessage("请先登录", c)
		return
	}
	idStr := c.Query("id")
	if idStr == "" {
		response.FailWithMessage("请传入预约 id", c)
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("id 格式错误", c)
		return
	}
	res, err := svcReservation.GetReservation(uint(id))
	if err != nil || res.ID == 0 {
		response.FailWithMessage("预约不存在", c)
		return
	}
	if res.UserID != userID {
		response.FailWithMessage("无权查看该预约", c)
		return
	}
	venueName := ""
	timeslotRange := ""
	if v, _ := svcVenue.GetVenue(res.VenueID); v.ID != 0 {
		venueName = v.Name
	}
	if s, _ := svcVenueTimeslot.GetVenueTimeslot(res.TimeslotID); s.ID != 0 {
		timeslotRange = s.StartTime.FormatHHMM() + "-" + s.EndTime.FormatHHMM()
	}

	canChange := false
	var lastChangeTime *time.Time

	// 仅在预约为未取消状态时计算退改能力和最晚时间
	if res.Status == 0 {
		var venue model.Venue
		if err := global.GVA_DB.Where("id = ?", res.VenueID).First(&venue).Error; err == nil && venue.ID != 0 && venue.RefundChangeHours > 0 {
			var slot model.VenueTimeslot
			if err := global.GVA_DB.Where("id = ?", res.TimeslotID).First(&slot).Error; err == nil {
				if startTime, err := combineDateAndTimeOnly(res.ReserveDate, slot.StartTime); err == nil {
					t := startTime.Add(-time.Duration(venue.RefundChangeHours) * time.Hour)
					lastChangeTime = &t
					now := time.Now()
					if now.Before(t) || now.Equal(t) {
						canChange = true
					}
				}
			}
		}
	}

	var lastChangeAt interface{}
	if lastChangeTime != nil {
		// 格式化为 "YYYY-MM-DD HH:MM"
		lastChangeAt = lastChangeTime.Format("2006-01-02 15:04")
	} else {
		lastChangeAt = nil
	}

	detail := gin.H{
		"id":           res.ID,
		"reservationNo": res.ReservationNo,
		"venueId":      res.VenueID,
		"venueName":    venueName,
		"timeslotId":   res.TimeslotID,
		"timeslotRange": timeslotRange,
		"reserveDate":  res.ReserveDate,
		"contactName":  res.ContactName,
		"contactPhone": res.ContactPhone,
		"contactCount": res.ContactCount,
		"status":       res.Status,
		"verifyCode":   res.VerifyCode,
		"createdAt":    res.CreatedAt,
		"canChange":    canChange,
		"lastChangeAt": lastChangeAt,
	}
	// 已核销时附带评价信息（有则返回，无则 null）
	if res.Status == 1 {
		review, _ := svcReservationReview.GetByReservationID(res.ID)
		if review.ID != 0 {
			detail["review"] = gin.H{
				"id":        review.ID,
				"rating":    review.Rating,
				"content":   review.Content,
				"createdAt": review.CreatedAt,
			}
		} else {
			detail["review"] = nil
		}
	}
	response.OkWithData(detail, c)
}

// Cancel 小程序-取消预约（仅本人）
// @Tags        小程序-露营
// @Summary     取消预约
// @Description 用户取消自己的预约。需先登录，请求头携带 x-token
// @Accept      json
// @Produce     json
// @Param       x-token header string false "小程序登录后返回的 token"
// @Param       id query int true "预约ID"
// @Success     200 {object} response.Response{msg=string}
// @Router      /camping/mini/reservation/cancel [post]
func (a *reservationApi) Cancel(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.FailWithMessage("请先登录", c)
		return
	}
	var idReq struct {
		ID uint `form:"id" json:"id" binding:"required"`
	}
	_ = c.ShouldBindJSON(&idReq)
	if idReq.ID == 0 {
		_ = c.ShouldBindQuery(&idReq)
	}
	if idReq.ID == 0 {
		response.FailWithMessage("请传入预约 id", c)
		return
	}
	res, err := svcReservation.GetReservation(idReq.ID)
	if err != nil || res.ID == 0 {
		response.FailWithMessage("预约不存在", c)
		return
	}
	if res.UserID != userID {
		response.FailWithMessage("无权取消该预约", c)
		return
	}
	if res.Status == 2 {
		response.FailWithMessage("该预约已取消", c)
		return
	}
	if err := svcReservation.CancelReservation(idReq.ID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("取消成功", c)
}

// CancelRule 小程序-取消规则说明
// @Tags        小程序-露营
// @Summary     取消规则
// @Description 根据场地返回取消规则文案（用于弹窗提示）
// @Accept      json
// @Produce     json
// @Param       venueId query int true "场地ID"
// @Success     200 {object} response.Response{data=object,msg=string}
// @Router      /camping/mini/reservation/cancelRule [get]
func (a *reservationApi) CancelRule(c *gin.Context) {
	var idReq struct {
		VenueID uint `form:"venueId" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage("请传入场地 id", c)
		return
	}
	venue, err := svcVenue.GetVenue(idReq.VenueID)
	if err != nil || venue.ID == 0 {
		response.FailWithMessage("场地不存在", c)
		return
	}
	var rule string
	if venue.RefundChangeHours <= 0 {
		rule = "不支持退改，取消后将不可恢复。"
	} else {
		rule = "使用前 " + strconv.Itoa(venue.RefundChangeHours) + " 小时可免费取消，逾期不可取消。"
	}
	response.OkWithData(gin.H{"rule": rule, "refundChangeHours": venue.RefundChangeHours}, c)
}

// CreateReview 小程序-发布预约评价（仅核销后的预约，一单一评）
// @Tags        小程序-露营
// @Summary     发布预约评价
// @Description 对已核销的预约进行评价（评分1-5、50字内内容），每个预约只能评价一次
// @Accept      json
// @Produce     json
// @Param       x-token header string false "小程序登录后返回的 token"
// @Param       data body request.CreateReservationReviewRequest true "评价内容"
// @Success     200 {object} response.Response{data=model.VenueReservationReview,msg=string}
// @Router      /camping/mini/reservation/review/create [post]
func (a *reservationApi) CreateReview(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.FailWithMessage("请先登录", c)
		return
	}
	var req campingRequest.CreateReservationReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	review, err := svcReservationReview.CreateReview(req, userID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(review, c)
}

// DeleteReview 小程序-删除预约评价（仅本人）
// @Tags        小程序-露营
// @Summary     删除预约评价
// @Description 删除自己对该预约的评价
// @Accept      json
// @Produce     json
// @Param       x-token header string false "小程序登录后返回的 token"
// @Param       id query int true "评价ID"
// @Success     200 {object} response.Response{msg=string}
// @Router      /camping/mini/reservation/review/delete [post]
func (a *reservationApi) DeleteReview(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.FailWithMessage("请先登录", c)
		return
	}
	var idReq struct {
		ID uint `form:"id" json:"id" binding:"required"`
	}
	_ = c.ShouldBindJSON(&idReq)
	if idReq.ID == 0 {
		_ = c.ShouldBindQuery(&idReq)
	}
	if idReq.ID == 0 {
		response.FailWithMessage("请传入评价 id", c)
		return
	}
	if err := svcReservationReview.DeleteReview(idReq.ID, userID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func getUserID(c *gin.Context) (uint, bool) {
	uid, exists := c.Get("x-user-id")
	if !exists {
		return 0, false
	}
	u, ok := uid.(uint)
	return u, ok
}

// combineDateAndTimeOnly 将日期与 TimeOnly 组合为本地时区的 time.Time
func combineDateAndTimeOnly(date time.Time, t model.TimeOnly) (time.Time, error) {
	s := string(t)
	loc := time.Local
	if s == "" {
		return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc), nil
	}
	var parsed time.Time
	var err error
	if len(s) >= 8 {
		parsed, err = time.ParseInLocation("15:04:05", s, loc)
	} else {
		parsed, err = time.ParseInLocation("15:04", s, loc)
	}
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(date.Year(), date.Month(), date.Day(), parsed.Hour(), parsed.Minute(), parsed.Second(), 0, loc), nil
}
