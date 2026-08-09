# 门票退款仅保留小程序端

## Goal

门票订单退款发起入口只保留小程序用户自助退票；管理后台不再提供退款 API、退款操作，以及「预计退款」展示。用户只能退本人订单；部分核销仍按剩余付费次数比例退款。

## Background

- 用户退款：`POST /mini/pay/refund` → `ticket.Order.RequestRefund(orderId, userID)`（校验本人）
- 后台退款（待删除）：`POST /ticket/order/refundOrder` → `AdminRefund` → `RequestRefund(orderID, 0)`
- 管理端 `web/src/plugin/ticket/view/order.vue` 列表/详情有退款按钮与「预计退款」；详情 API 返回 `canRefund` / `refundAmount*`
- 现有公式已支持部分核销：`remainingPaid = A - min(verifiedTimes, A)`，`refundFen = Round(payAmount*100*remainingPaid/A)`；`status==1` 且 `remainingPaid>0` 即可退（含已核销部分次数）

## Requirements

1. 删除后台门票退款 API：`POST /ticket/order/refundOrder`（路由、handler、插件 API 注册）
2. 删除管理端门票订单退款操作（列表/详情按钮、确认文案、`refundOrder` 封装）
3. 管理端去掉「预计退款」只读展示及相关前端状态字段；**保留**已退款订单的「实退金额」与状态标签（已退款/退款中）
4. 后台 `findOrder` 不再返回退款预估字段（`canRefund` / `refundAmountFen` / `refundAmount`）；小程序详情若仍依赖这些字段则保留
5. 删除仅服务后台的 `AdminRefund`；`RequestRefund` 仅供小程序路径，且 **强制 `userID > 0` 并校验订单归属本人**（他人订单拒绝）
6. **保留**按付费次数剩余比例退款：未核销全额退；部分核销按 `remainingPaid/A` 比例退；付费次数用尽不可退；赠送次数不计入分母
7. 小程序退款入口、微信回调、日历 `sold` 回退逻辑保持现有行为
8. 更新 `.trellis/spec/backend/ticket-gift-proportional-refund.md`：退款仅小程序本人发起；明确部分核销比例退
9. 不改 limitedActivity 后台退款能力

## Out of Scope

- 改动退款金额公式本身（沿用现有比例逻辑）
- 景点 `RefundChangeHours` 时间窗
- 限时活动退款

## Acceptance Criteria

- [ ] `POST /ticket/order/refundOrder` 不再注册
- [ ] 管理端订单页无退款按钮、无「预计退款」展示；已退款仍可见实退金额与状态
- [ ] 后台 `findOrder` 响应不含 `canRefund` / `refundAmountFen` / `refundAmount`
- [ ] 小程序仅能对**本人**订单调用 `/mini/pay/refund`；非本人返回无权/不存在类错误
- [ ] `AdminRefund` 已移除；`RequestRefund` 要求 `userID > 0` 且校验 `order.UserID == userID`
- [ ] 部分核销（`status=1` 且 `0 < verifiedTimes < A`）仍可退，金额为付费剩余比例；`remainingPaid==0` 不可退
- [ ] spec 已改为「仅小程序本人发起 + 部分核销比例退」；limitedActivity 后台退款未受影响

## Notes

- Lightweight：PRD-only；实现面主要为入口清理 + 收紧 `RequestRefund` 归属校验 + 管理端 UI + spec。
