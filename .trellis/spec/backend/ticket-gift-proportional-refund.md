# Ticket Gift Times & Proportional Refund

> 门票赠送次数、次数×数量快照、按付费次数比例整单退款。退款仅小程序用户本人自助发起（无后台退款入口）。

---

## Scenario: 多次票赠送 + 整单按 A 比例退

### 1. Scope / Trigger

- Trigger: 改门票下单次数、核销次数语义、用户退款金额或日历 `sold` 回退
- 改动面：`server/plugin/ticket/`（model/service/api）+ `server/api/v1/mini/pay.go`（仅 ticket 退款委托）+ 管理端 `product.vue`/`order.vue`（展示；无退款操作）+ H5 `h5Verify.vue` 门票分支展示
- 勿拆按张凭证表；勿在 limitedActivity 复制本公式（活动仍按 `remaining/totalUseTimes`）；勿改 limitedActivity 后台退款
- 对齐模式：`limitedActivity` 的「先 status=7 再调微信」，但分母用付费次数 A，不是 total；门票退款入口仅小程序

### 2. Signatures

| 层 | 签名 / 字段 |
|----|-------------|
| DB `ticket_sku` | `gift_use_times` int default 0（单张赠送 p；仅 `ticket_type=2` 有意义） |
| DB `orders` | `paid_use_times`（A）、`gift_use_times`（B 快照）、`refund_amount` decimal；`total_use_times`=A+B |
| Service | `RemainingPaidUseTimes(order) int` |
| Service | `CalcRefundFen(order) (refundFen, totalFen, remainingPaid int, err error)` |
| Service | `RequestRefund(orderID, userID uint) error`（**强制 userID>0** 且 `order.UserID == userID`；否则无权） |
| Service | `ApplyRefundSuccessByRefundNo` / `ReleaseRefundRequested` |
| 用户退款 | `POST /mini/pay/refund` → `ticket.Order.RequestRefund(orderId, userID)`（唯一发起入口） |

下单快照（`CreateOrder`）：

```go
m := sku.UseTimes; if m <= 0 { m = 1 }
p := sku.GiftUseTimes; if sku.TicketType != 2 { p = 0 }
// ticketType==2 时 quantity 必须为 1
A := quantity * m
B := quantity * p
PaidUseTimes, GiftUseTimes, TotalUseTimes = A, B, A+B
```

### 3. Contracts

**次数语义**

| 符号 | 含义 |
|------|------|
| A (`paidUseTimes`) | 付费可核销次数 = `quantity × useTimes` |
| B (`giftUseTimes`) | 赠送可核销次数 = `quantity × giftUseTimes` |
| total | A+B；核销上限 |
| 核销顺序 | 共用 `verifiedTimes` 累加；语义上先消耗 A 再消耗 B |
| 核销码 | 仍为 `order_no`；不拆码 |

**退款公式（部分核销仍按剩余付费比例）**

```
remainingPaid = A - min(verifiedTimes, A)
refundFen = Round(payAmount * 100 * remainingPaid / A)  // remainingPaid>0
```

- 未核销（`verifiedTimes==0`）：全额退
- 部分核销（`0 < verifiedTimes < A`，且 `status==1`）：按 `remainingPaid/A` 比例退
- 付费次数用尽（`remainingPaid==0`）：不可退
- 赠送次数不进入分子/分母

**sold 回退（退款成功）**

| 票种 | release |
|------|---------|
| 多次票 `ticket_type==2` | `1`（每单仅 1 张） |
| 单次票 | `remainingPaid`（且 ≤ `quantity`） |
| SKU 缺失兜底 | `quantity==1 && A>1` 视为多次票 |

**API 展示字段**

- 小程序订单详情：可返回 `paidUseTimes`、`giftUseTimes`、`canRefund`、`refundAmountFen` / `refundAmount`（预估）
- 管理端 `findOrder`：返回 `paidUseTimes`、`giftUseTimes`、`remainingTimes` 等；**不**返回 `canRefund` / `refundAmountFen` / `refundAmount` 预估

**老单兼容**：`paidUseTimes<=0` 时 `A=totalUseTimes`（或 1），`B=0`。

### 4. Validation & Error Matrix

| 条件 | 行为 / 错误 |
|------|-------------|
| 多次票 `quantity≠1` | 下单拒绝：「多次票每次仅限购买一张」 |
| SKU 单次票 | 保存时强制 `giftUseTimes=0` |
| `userID==0` 或 `order.UserID != userID` | 「订单不存在或无权操作」 |
| `status≠1` | 不可退 / 不可核销（含 status=7 退款中） |
| `remainingPaid<=0` | 「付费次数已用尽，不可退款」 |
| 无 `wx_transaction_id` | 「订单缺少微信支付信息」 |
| 已有 `refund_no` | 「退款处理中或已退款」 |
| 微信 CreateRefund 失败 | `ReleaseRefundRequested` 恢复 status=1 |

### 5. Good / Base / Bad Cases

- **Good**：多次票 m=5,p=1，quantity=1 → total=6；核销 2 次后退款 → `refundFen=Round(pay*100*3/5)`；sold-=1
- **Base**：单次票 quantity=3 → A=3,B=0；核销 1 次 → 退 2/3；sold-=2
- **Bad**：用 `totalUseTimes`（含赠送）做退款分母；或只改 status=6 不调微信；或 `verifiedTimes>0` 一律拒退；或后台/匿名（userID=0）发起退款

### 6. Tests Required

- CreateOrder：多次 quantity=2 失败；quantity=1 且 gift 写入 A/B/total
- CalcRefundFen：k=0 全额；0&lt;k&lt;A 比例；k≥A 错误；分位 Round
- RequestRefund：仅本人（userID>0）；非本人/userID=0 拒绝；先 7 再微信；失败释放
- ApplyRefundSuccess：refund_amount；多次 sold-=1；单次 sold-=remainingPaid
- 断言：小程序 `canRefund` 与 Detail 预估金额一致；管理端无 refundOrder；limitedActivity `/mini/pay` 与后台退款分支不受影响

### 7. Wrong vs Correct

#### Wrong

```go
// 有核销就拒；或按 total（含赠送）比例退；后台只 Updates status=6；或允许 userID==0
if order.VerifiedTimes > 0 { return errors.New("不可退") }
refundFen = Round(pay * 100 * (total-verified) / total)
_ = RequestRefund(orderID, 0) // 后台旁路
```

#### Correct

```go
remainingPaid := RemainingPaidUseTimes(order) // 只看 A
refundFen, totalFen, _, err := CalcRefundFen(order)
// 先 status=7，再 mini.CreateRefund(..., totalFen, refundFen, ...)
// 仅小程序：RequestRefund(orderID, userID) 且 userID>0 && order.UserID==userID
```

---

## Design Decision: 不拆码 + 多次票限购 1

**Context**：曾考虑按张凭证多码；产品改为不拆码，并用「多次票每单 1 张」降低共享次数池歧义。

**Decision**：一单一码（`order_no`）；单次多张共享次数池；多次票强制 `quantity=1`；退款整单一笔微信单；发起入口仅小程序本人。

**Why**：支付/退款/列表模型保持简单；赠送与比例退仍可在订单级表达；管理端只读审计（实退金额/状态），不代退。

---

## Gotchas

1. **退款中禁止核销**：`VerifyOrder` 仅 `status==1`，避免 status=7 窗口内继续核销导致比例漂移。
2. **小程序确认金额**：用分单位 `Math.round` 再 /100，与 `CalcRefundFen` 一致，勿直接用浮点金额展示。
3. **公共 pay.go**：ticket 退款只委托 Service；勿在 `pay.go` 再写一套金额公式；改 limitedActivity 分支时勿连带改坏 ticket。
4. **无后台退款**：门票管理端不注册 `refundOrder`；limitedActivity 后台退款保持独立，勿误删。
