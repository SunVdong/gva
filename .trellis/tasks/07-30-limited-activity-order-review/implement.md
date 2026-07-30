# Implement: 限时活动订单评价

## Checklist

1. [x] 新增 `model/order_review.go` + `model/request/order_review.go`
2. [x] 新增 `service/order_review.go`，`enter.go` 注册 `OrderReview`
3. [x] `initialize/gorm.go` AutoMigrate 新表
4. [x] 小程序 `CreateReview` / `DeleteReview` API + `router/mini.go` 路由
5. [x] `api/mini/order.go` Detail 已核销时附带 `review`；`mini/enter.go` 注入
6. [x] 后台 `api/order.go` Find 已核销时附带 `review`；`api/enter.go` 注入
7. [x] 后台 `web/src/plugin/limitedActivity/view/order.vue` 详情展示评价块
8. [x] 本地编译通过：`go build ./plugin/limitedActivity/...`

## Validation

```bash
cd server && go build ./plugin/limitedActivity/...
```

手动（有环境时）：

1. 已核销活动订单，本人 token → create 成功
2. 同单再 create → 已评价过
3. detail / 后台 findOrder 可见 review
4. delete 后可再 create
5. 非本人 / status≠2 → 对应错误

## Review Gates

- 行为与文案对齐门票 `ticket` 评价
- 不改门票 / 露营评价代码
- 表名 `limited_activity_order_reviews`

## Rollback

回退本任务提交；新表可留空。
