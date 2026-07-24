# 限时活动 — 技术设计

## Architecture & Boundaries

```
小程序(仓外) --mini API--> limitedActivity 插件
管理端 web  --admin API--> limitedActivity 插件
公开 H5     --public API--> limitedActivity 插件
支付/退款回调 --/mini/pay/*--> pay.go 分发到 limitedActivity 订单处理
```

- **新插件包名**：`limitedActivity`（避免与 `activityGuide` 冲突）。
- **职责内**：活动 CRUD、Banner CRUD、报名订单、名额、核销记录、后台按比例退款、mini/公开 API、管理端页面、H5 type 扩展。
- **职责外复用**：`service/mini` 微信支付；现有 `#/h5/verify` 页；用户体系 `users` + `x-token`。
- **分层**：Router → API → Service → Model；API 不直接写 DB。

## Data Model

### Activity（限时活动）

| 字段 | 说明 |
|------|------|
| name / detail | 名称、详情 |
| start_time / end_time | 活动时间窗（报名与「进行中」判断） |
| market_price / price | 市场价、实际价 |
| quota / sold | 总名额、已占用（含待支付占用） |
| cover_image / group_qr / service_qr | 封面、群二维码、客服二维码 |
| status | 显示状态（如 1 显示 / 0 隐藏） |
| 审计 | 复用 `GVA_MODEL`（含 CreatedAt/UpdatedAt）；创建人/更新人按项目惯例字段写入 |

### ActivityOrder（参与订单）

状态对齐 ticket：`0待支付 1待核销 2已核销 5已关闭 6已退款 7退款中`（本业务可不引入 3/4，或预留）。

| 字段 | 说明 |
|------|------|
| order_no | 前缀 `A` + 时间戳 + 随机 |
| user_id | 小程序用户 |
| activity_id + 活动快照名 | 下单时快照名称等 |
| contact_name / contact_phone | 联系人 |
| quantity | 人次 |
| unit_price / pay_amount | 实际单价、应付总额 |
| total_use_times / verified_times | 可核销/已核销 |
| pay_time / wx_transaction_id | 支付信息 |
| refund_no / wx_refund_id / refund_time / refund_amount | 退款单与实退金额（分或元与 ticket 一致用元+分运算） |

### ActivityOrderVerifyRecord

`order_id`, `verify_no`, `verified_at`（对齐 ticket `OrderVerifyRecord`）。

### Banner（活动列表轮播）

表名：`limited_activity_banners`。

| 字段 | 说明 |
|------|------|
| title | 后台识别用标题（小程序可不展示） |
| image | 轮播图 URL（必填） |
| detail_image | 详情长图 URL（必填；点击轮播图后展示） |
| sort | 排序，越小越靠前 |
| status | `1` 显示 / `0` 隐藏 |
| 审计 | `GVA_MODEL` + created_by / updated_by |

**无** 跳转类型、activity_id、link_path、start_time / end_time。

## Key Flows

### 报名占用

1. 校验活动显示、时间窗、`quota - sold >= quantity`。
2. 事务：`sold += quantity`，创建 `status=0` 订单，`total_use_times=quantity`。
3. 待支付超时关闭：`status=5`，`sold -= quantity`（可复用 ticket 超时任务模式或同库定时逻辑）。

### 支付

1. 小程序：`POST /mini/pay/create`，`orderType=limitedActivity`，`orderId`。
2. `pay.go` Create 分支：校验本人、status=0、金额转分 → `CreateJSAPI`；`out_trade_no = orderNo_unix`。
3. Notify：订单号前缀 `A` → 验金额与 transaction_id → `status=1`，幂等。

### 核销

- 公开：`GET/POST` 按 orderNo 查询与核销（命名对齐 ticket `*ByCodePublic`）。
- 每次：`verified_times++`，写记录；达到 `total_use_times` → `status=2`。
- H5：`#/h5/verify?type=limitedActivity&code=<orderNo>`，扩展 `h5Verify.vue`。

### 后台按比例退款

1. 条件：`status=1` 且 `remaining = total_use_times - verified_times > 0`，且有 `wx_transaction_id`，无进行中冲突退款单。
2. `refundFen = Round(pay_amount*100 * remaining / total_use_times)`；`totalFen = Round(pay_amount*100)`。
3. 调 `mini.CreateRefund`；标记 `status=7`；回调成功 → `status=6`，写 `refund_amount`，`sold -= remaining`；此后不可核销。
4. **不**对用户开放 `/mini/pay/refund` 活动分支（避免自助退）；仅管理端 API。

### Banner 下发

1. Admin CRUD；保存时校验 `image`、`detail_image` 均非空。
2. Mini 列表：`status=1`，按 `sort ASC, id ASC`。
3. 返回 `image`、`detailImage` 等字段；点击后展示长图由小程序完成，本仓不实现 UI。

## API Surface（摘要）

| 端 | 能力 |
|----|------|
| Admin | 活动 CRUD；Banner CRUD；订单列表/详情/核销记录；后台退款 |
| Mini | Banner 列表；活动列表/详情；下单；我的订单列表/详情（含核销进度、群二维码、客服二维码） |
| Public | 按 code 查询、核销 |
| Pay | Create/Notify 增加 `limitedActivity`；RefundNotify 增加活动订单成功处理 |

## Compatibility & Trade-offs

- 按比例退款与 ticket「有核销不可退」不同：活动订单退款逻辑独立，勿复用 ticket 用户退款函数。
- `sold` 含待支付占用：简单，需严格超时释放；备选「仅已支付占名额」会超卖风险更高，不采用。
- 小程序 UI 仓外：以 Swagger/接口约定交付；管理端与 H5 本仓交付。

## Rollback

- 功能开关：插件未注册 / 菜单隐藏即可停用前台入口。
- 表由 AutoMigrate 创建；回滚以停用插件为主，不强制删表。
