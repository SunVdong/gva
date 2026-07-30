# Design: 限时活动订单评价

## Approach

镜像门票插件 `ticket` 的 `OrderReview` 全链路，落在 `limitedActivity` 插件内，表名与模型独立，避免与门票评价表冲突。

## Data Model

**表** `limited_activity_order_reviews`

| 字段 | 说明 |
|------|------|
| id / created_at / updated_at / deleted_at | GVA_MODEL |
| order_id | 活动订单 ID，`uniqueIndex` |
| user_id | 评价用户 ID（小程序 `users.id`） |
| rating | 1–5 |
| content | ≤50 字 |

模型：`server/plugin/limitedActivity/model/order_review.go`  
请求：`model/request/order_review.go`（`CreateActivityOrderReviewRequest`）

## Service Rules

`service/order_review.go`（对齐 `ticket/service/order_review.go`）：

1. `CreateReview`：查 `ActivityOrder` → 本人 → `status==2` → `VerifiedAt!=nil` → 未存在评价 → 创建
2. `DeleteReview`：存在 → 本人 → **硬删**（`Unscoped`；MySQL 下 `order_id` 唯一索引含软删行，软删会导致删后重评失败）
3. `GetByOrderID`：详情聚合用

错误文案与门票一致：`无权对该订单评价` / `仅已核销订单可评价` / `仅核销后的订单可评价` / `该订单已评价过，可先删除再重新评价` 等。

## API Surface

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| POST | `/limitedActivity/mini/order/review/create` | JWT | 发表评价 |
| POST | `/limitedActivity/mini/order/review/delete` | JWT | 删除评价 |
| GET | `/limitedActivity/mini/order/detail` | JWT | 已核销时附 `review` |
| GET | `/limitedActivity/order/findOrder` | 后台 | 已核销时附 `review` |

`review` 载荷形状与门票一致：

```json
{ "ID": 1, "rating": 5, "content": "...", "createdAt": "..." }
```

未评价时为 `null`；非已核销订单可不返回 `review` 字段（与门票 `Find` / mini `Detail` 行为一致）。

## Wiring

- `initialize/gorm.go`：AutoMigrate 新模型
- `service/enter.go`：挂 `OrderReview`
- `api/enter.go`、`api/mini/enter.go`：注入 service
- `router/mini.go`：注册 create / delete
- 后台 `api/order.go` `Find`、小程序 `api/mini/order.go` `Detail`：聚合 review
- 后台 `web/.../limitedActivity/view/order.vue`：详情区展示评价

小程序评价接口为 public+JWT，不强制写入 `initialize/api.go` 的后台 Casbin 列表（与门票 mini review 一致）；若项目惯例要求登记 swagger，仅补注释即可。

## Compatibility / Rollback

- 新增表与只读字段聚合，对旧客户端向后兼容（多返回 `review`）
- 回滚：停用路由 + 忽略新表即可；已写入评价数据可保留

## Risks

- 表名必须独立，禁止复用 `order_reviews`
- 评价条件依赖「全部核销完成」；部分核销（status=1）不可评，与门票一致
