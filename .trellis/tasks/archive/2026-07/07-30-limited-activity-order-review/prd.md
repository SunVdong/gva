# 限时活动订单评价与删除

## Goal

为限时活动订单补齐与门票一致的评价能力：小程序端可对已核销完成的本人订单发表/删除评价；后台订单详情可查看评价内容。

## Requirements

1. **评价规则（对齐门票）**
   - 仅订单本人可评价 / 删除评价
   - 仅订单状态为已核销（`status = 2`）且 `verifiedAt` 不为空时可评价
   - 一单一评；已评价需先删除再重新评价
   - 评分 `rating`：1–5（必填）；内容 `content`：最多 50 字（可选）

2. **小程序接口**
   - `POST /limitedActivity/mini/order/review/create`：创建评价，body `{ orderId, rating, content }`，需登录
   - `POST /limitedActivity/mini/order/review/delete`：删除评价，传评价 `id`，需登录
   - `GET /limitedActivity/mini/order/detail`：已核销订单详情附带 `review`（无则 `null`）

3. **后台展示**
   - `GET /limitedActivity/order/findOrder`：已核销订单详情附带 `review`
   - 后台订单详情抽屉展示评分、内容、评价时间（对齐门票订单页）

4. **数据**
   - 独立评价表（勿复用门票 `order_reviews`），按活动订单 ID 唯一索引
   - 服务启动 AutoMigrate 建表

## Out of Scope

- 评价列表管理页、审核、回复
- 修改评价（仅删后重评）
- 限时活动活动详情页聚合评分展示
- 改动门票 / 露营评价逻辑

## Acceptance Criteria

- [ ] 本人对已核销完成的活动订单可成功评价一次
- [ ] 非本人 / 未核销完成 / 已评价再评，返回与门票同等语义的错误文案
- [ ] 本人可删除自己的评价；删除后可重新评价
- [ ] 小程序订单详情在已核销时返回 `review`
- [ ] 后台订单详情在已核销时展示评价信息（无评价显示 `-`）
- [ ] 重启服务后评价表自动迁移成功

## Notes

- 用户确认：规则对齐门票；后台也展示评价。
