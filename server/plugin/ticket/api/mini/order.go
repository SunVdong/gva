package mini

import (
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model/request"
	"github.com/gin-gonic/gin"
)

var Order = new(miniOrderApi)

type miniOrderApi struct{}

// Create 小程序-提交订单（需登录，请求头带 x-token）
// @Tags        小程序-景点
// @Summary     提交订单
// @Description 小程序端提交门票订单，需先登录，请求头携带 x-token。请求体不需传 userId，用户身份由 x-token 解析注入。创建成功后请携带 x-token 调用公共接口 POST /mini/pay/create，body 传 {"orderType":"ticket","orderId": 订单ID}，获取支付参数后调 wx.requestPayment 完成支付。
// @Accept      json
// @Produce     json
// @Param       x-token header string true "小程序登录后返回的 token（必填，用于识别用户）"
// @Param       data body request.MiniOrderCreate true "订单信息（bookerName、bookerPhone、items）"
// @Success     200  {object} response.Response{data=object,msg=string}
// @Router      /ticket/mini/order/create [post]
func (a *miniOrderApi) Create(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok || userID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}
	var req request.MiniOrderCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	order, err := svcOrder.CreateOrder(userID, req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(order, c)
}

// MyList 小程序-我的订单列表（需登录）
// @Tags        小程序-景点
// @Summary     我的订单列表
// @Description 按类型筛选：待支付、待核销、已完成（已完成含已核销/已取消/已过期/已关闭）；不传 orderType 返回全部。列表中每条订单带 statusLabel 表明状态。
// @Accept      json
// @Produce     json
// @Param       x-token   header string false "小程序登录后返回的 token"
// @Param       orderType query  string false "可选值：pending_payment|pending_verify|completed，代表，待支付|待核销|已完成， 不传默认全部"
// @Param       page      query  int    false "页码"
// @Param       pageSize  query  int    false "每页条数"
// @Success     200       {object} response.Response{data=response.PageResult,msg=string}
// @Router      /ticket/mini/order/myList [get]
func (a *miniOrderApi) MyList(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok || userID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}
	var req request.TicketOrderSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = userID
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 50 {
		req.PageSize = 20
	}
	list, total, err := svcOrder.GetList(req)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	// 批量查订单项最晚游玩日，用于生成 statusLabel
	orderIDs := make([]uint, 0, len(list))
	for _, o := range list {
		orderIDs = append(orderIDs, o.ID)
	}
	maxVisitMap, _ := svcOrder.GetMaxVisitDateByOrderIDs(orderIDs)
	skuNamesMap, _ := svcOrder.GetSkuNamesByOrderIDs(orderIDs)
	productNamesMap, _ := svcOrder.GetProductNamesByOrderIDs(orderIDs)
	items := make([]gin.H, 0, len(list))
	for _, o := range list {
		maxVisit := maxVisitMap[o.ID]
		skuNames := skuNamesMap[o.ID]
		productNames := productNamesMap[o.ID]
		items = append(items, gin.H{
			"id":          o.ID,
			"orderNo":     o.OrderNo,
			"userId":      o.UserID,
			"bookerName":  o.BookerName,
			"bookerPhone": o.BookerPhone,
			"totalAmount": o.TotalAmount,
			"payAmount":   o.PayAmount,
			"status":      o.Status,
			"payTime":     o.PayTime,
			"verifiedAt":  o.VerifiedAt,
			"createdAt":   o.CreatedAt,
			"statusLabel": svcOrder.OrderStatusLabel(&o, maxVisit),
			"skuNames":    skuNames,
			"skuNameText": strings.Join(skuNames, "、"),
			"productNames":    productNames,
			"productNameText": strings.Join(productNames, "、"),
		})
	}
	response.OkWithDetailed(response.PageResult{
		List: items, Total: total, Page: req.Page, PageSize: req.PageSize,
	}, "获取成功", c)
}

// Detail 小程序-订单详情（含订单项，仅本人）
// @Tags        小程序-景点
// @Summary     订单详情
// @Description 小程序端获取订单详情及订单项，仅限当前登录用户自己的订单
// @Accept      json
// @Produce     json
// @Param       x-token header string false "小程序登录后返回的 token"
// @Param       id query int true "订单ID"
// @Success     200 {object} response.Response{data=object,msg=string}
// @Router      /ticket/mini/order/detail [get]
func (a *miniOrderApi) Detail(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok || userID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}
	var idReq struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	order, items, err := svcOrder.GetByID(idReq.ID)
	if err != nil || order.ID == 0 {
		response.FailWithMessage("订单不存在", c)
		return
	}
	if order.UserID != userID {
		response.FailWithMessage("无权查看该订单", c)
		return
	}
	data := gin.H{"order": order, "items": items}
	// 已核销时附带评价信息（有则返回，无则 null）
	if order.Status == 2 && order.VerifiedAt != nil {
		review, _ := svcOrderReview.GetByOrderID(order.ID)
		if review.ID != 0 {
			data["review"] = gin.H{
				"id":        review.ID,
				"rating":    review.Rating,
				"content":   review.Content,
				"createdAt": review.CreatedAt,
			}
		} else {
			data["review"] = nil
		}
	}
	response.OkWithData(data, c)
}

// CreateReview 小程序-发布订单评价（仅核销后的订单，一单一评）
// @Tags        小程序-景点
// @Summary     发布订单评价
// @Description 对已核销的门票订单进行评价（评分1-5、50字内内容），每个订单只能评价一次
// @Accept      json
// @Produce     json
// @Param       x-token header string false "小程序登录后返回的 token"
// @Param       data body request.CreateOrderReviewRequest true "评价内容"
// @Success     200 {object} response.Response{data=model.OrderReview,msg=string}
// @Router      /ticket/mini/order/review/create [post]
func (a *miniOrderApi) CreateReview(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok || userID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}
	var req request.CreateOrderReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	review, err := svcOrderReview.CreateReview(req, userID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(review, c)
}

// DeleteReview 小程序-删除订单评价（仅本人）
// @Tags        小程序-景点
// @Summary     删除订单评价
// @Description 删除自己对该订单的评价
// @Accept      json
// @Produce     json
// @Param       x-token header string false "小程序登录后返回的 token"
// @Param       id query int true "评价ID"
// @Success     200 {object} response.Response{msg=string}
// @Router      /ticket/mini/order/review/delete [post]
func (a *miniOrderApi) DeleteReview(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok || userID == 0 {
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
	if err := svcOrderReview.DeleteReview(idReq.ID, userID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// getUserID 从上下文中获取 OptionalJWTAuth 注入的用户 ID
func getUserID(c *gin.Context) (uint, bool) {
	uid, exists := c.Get("x-user-id")
	if !exists || uid == nil {
		return 0, false
	}
	u, ok := uid.(uint)
	return u, ok
}
