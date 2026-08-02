# Design: 门票赠送次数与整单比例退款

## Architecture / Boundaries

- **支付/退款入口**：`server/api/v1/mini/pay.go`（用户退款）；后台退款走 ticket API → service（对齐 limitedActivity `AdminRefund`）。
- **业务规则**：集中在 `server/plugin/ticket/service/order.go`（下单快照、核销、`CalcRefundFen`、退款成功落库与 sold 回退）。
- **配置面**：SKU `giftUseTimes`；管理端 `product.vue`；订单展示 `order.vue` / mini API / H5 核销展示字段。
- **不改**：核销码体系（仍 `order_no`）；不引入凭证子表。

## Data Model

### TicketSku
- 新增 `GiftUseTimes int`（`gift_use_times`，default 0；仅多次票有意义）。

### TicketOrder
- 新增 `PaidUseTimes int`（付费总次数 A 快照）
- 新增 `GiftUseTimes int`（赠送总次数 B 快照；注意与 SKU 字段同名 JSON 时订单侧为订单快照）
- 新增 `RefundAmount float64`（实退金额，元）
- `TotalUseTimes` 语义改为 A+B（下单时写入 `quantity×(m+p)`）

AutoMigrate 已有：`server/plugin/ticket/initialize/gorm.go`。

### 兼容
- 读单计算时：若 `PaidUseTimes<=0`，则 `A=TotalUseTimes`（或 `max(TotalUseTimes,1)`），`B=0`。

## Contracts / Formulas

```
m = sku.UseTimes; if m<=0 { m=1 }
p = sku.GiftUseTimes; if ticketType!=2 { p=0 }
A = quantity * m
B = quantity * p
totalUseTimes = A + B

paidConsumed = min(verifiedTimes, A)
remainingPaid = A - paidConsumed
refundFen = Round(payAmount*100 * remainingPaid / A)  // A>0 && remainingPaid>0
```

### sold 回退（退款成功事务内）
- 单次票（订单快照 m 可由 `A/quantity` 得，或下单额外不存 m）：`release = remainingPaid`（因 m=1）
- 多次票：`release = 1`（quantity 恒为 1）
- 实现建议：下单可冗余不存 m；用 `SkuTicketType` 或 `A==quantity`（单次）/ `ticketType` join 判断；更稳妥是退款时若 `quantity==1 && A>1` 或 join SKU type==2 → 多次；否则 `release=remainingPaid`（单次 m=1）。

更稳约定：
- `ticketType==2` → `sold -= 1`
- 否则 → `sold -= remainingPaid`（且 `release` 不超过 `quantity`）

## Data Flow

### 下单 CreateOrder
1. 多次票强制 `quantity==1`
2. 计算 A/B/total，写入订单
3. 日历 `sold += quantity`（不变）

### 核销 VerifyOrder
1. 仍按 `totalUseTimes` 上限累加（含赠送）
2. 退款中 status=7 不可核销（与 limitedActivity 一致，避免比例漂移）

### 用户退款 `/mini/pay/refund`
1. 锁单 → `CalcRefundFen` → status=7 + refund_no
2. `CreateRefund(txId, refundNo, totalFen, refundFen)`
3. SUCCESS 则 `ApplyRefundSuccess`（status=6、refundAmount、sold 回退）

### 后台退款
- 删除/替换 `RefundPendingVerifyMultiTicket`「仅改状态」；改为与用户端同一 service 方法（可抽 `RequestRefund(orderID)` 供两端调用；用户端额外校验 userID）。

## Compatibility / Migration

- 历史订单：无 A/B 时兼容公式；比例退对老多次票等价「按 total 剩余退」（B=0）。
- 历史多次票 quantity>1（若存在）：本需求起禁止；存量不强制拆分，按订单快照计算。

## Trade-offs

| 选择 | 取舍 |
|------|------|
| 不拆码 | 单次多张共享次数池，现场一码；实现简单 |
| 退款只看 A | 赠送核销不增加可退额度；A 用尽不可退 |
| 对齐 limitedActivity | 复用锁单→微信→回调模式，降低双轨风险 |

## Rollback

- 配置层：赠送次数置 0 即无赠送。
- 代码回滚：恢复全额/未核销才退；`totalUseTimes` 不乘数量（需评估已产生新单数据）。
- 微信退款中订单依赖 status=7 释放逻辑，回滚时保留该路径。
