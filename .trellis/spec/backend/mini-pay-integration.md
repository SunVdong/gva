# Mini Pay Plugin Integration

> 新业务插件如何接入公共微信支付（`/mini/pay/*`）。限时活动（`limitedActivity`）落地后沉淀。

---

## Scenario: 插件订单接入公共支付

### 1. Scope / Trigger

- Trigger: 新插件需要小程序 JSAPI 支付 / 支付回调 /（可选）退款回调
- 改动面：`server/api/v1/mini/pay.go` + 插件订单 Service；底层勿重写 `server/service/mini/pay.go`

### 2. Signatures

| 能力 | 入口 | 说明 |
|------|------|------|
| 调起支付 | `POST /mini/pay/create` | body: `orderType`, `orderId` |
| 支付回调 | `POST /mini/pay/notify` | 微信服务器；按订单号前缀分发 |
| 用户退款 | `POST /mini/pay/refund` | 仅已接入且允许自助退的业务（当前仅 ticket） |
| 退款回调 | `POST /mini/pay/refund/notify` | 按退款单/订单归属分发 |

订单号约定：

- ticket：`T` + 时间戳 + 随机
- limitedActivity：`A` + 时间戳 + 随机
- 预支付 `out_trade_no` = `{orderNo}_{unix}`；回调时取第一个 `_` 之前为业务订单号

### 3. Contracts

**Create 请求**

| 字段 | 类型 | 约束 |
|------|------|------|
| orderType | string | 必填；已注册业务键，如 `ticket` / `limitedActivity` |
| orderId | uint | 必填；当前登录用户本人订单 |

**Create 成功 data**：小程序 `wx.requestPayment` 参数（`timeStamp/nonceStr/package/signType/paySign`）。

**业务订单最低字段**：`order_no`、`user_id`、`pay_amount`、`status`（0 待支付）、支付后 `pay_time` / `wx_transaction_id`；若支持退款另需 `refund_no` / `wx_refund_id` / `refund_time`（及可选 `refund_amount`）。

**环境（已有 mini 配置）**：商户号、APIv3、证书、`notify-url` / `refund-notify-url`（见 `server/config/mini.go`）。

### 4. Validation & Error Matrix

| 条件 | 行为 |
|------|------|
| 未登录 / 无 openid | Fail：「请先登录」/「用户未绑定微信」 |
| 订单不存在或非本人 | Fail：「订单不存在或无权支付」 |
| status ≠ 待支付 | Fail：「订单状态不允许支付」 |
| pay_amount 分 ≤ 0 | Fail：「订单金额异常」 |
| 未知 orderType | Fail：「不支持的订单类型」 |
| 回调金额与订单不一致 | Notify FAIL；不改状态 |
| 重复同一 transaction_id | 幂等成功 |
| 不同 transaction_id 重复付 | 按业务拒绝或记冲突（对齐 ticket 断言） |

### 5. Good / Base / Bad Cases

- **Good**：`orderType=limitedActivity`，待支付单，Create → 用户支付 → Notify 入账 `status=待核销`
- **Base**：ticket 既有路径保持不变
- **Bad**：新业务只写插件下单、不扩展 `pay.go` Create/Notify → 无法支付或回调 FAIL「未知的订单类型」

### 6. Tests Required

- Create：本人/非本人、错误状态、金额 0、未知 orderType
- Notify：金额一致幂等；金额不一致失败；前缀路由到正确插件
- 退款（若后台退）：先锁单再调微信，避免核销插队；失败释放退款中状态
- 断言点：`status`、`wx_transaction_id`、名额/库存回退次数

### 7. Wrong vs Correct

#### Wrong

```go
// 在插件内再封装一套微信下单；或 Notify 不认新前缀
case "ticket":
    // ...
default:
    // 新业务永远走到这里
```

#### Correct

```go
switch req.OrderType {
case "ticket":
    // 既有
case "limitedActivity":
    // 查本插件订单 → CreateJSAPI(outTradeNo, fen, desc, openID, ip)
default:
    response.FailWithMessage("不支持的订单类型", c)
}
// Notify：orderNo[0]=='A' → 活动入账；=='T' → 门票入账
```

---

## Design Decision: 用户自助退 vs 后台退

**Context**：限时活动要求加客服后由后台退，且支持按未核销比例退款。

**Decision**：

- ticket：可走 `/mini/pay/refund`（有核销次数则拒）
- limitedActivity：**不**在 `/mini/pay/refund` 增加分支；仅管理端按  
  `refundFen = Round(payAmount*100 * remaining/total)` 调 `mini.CreateRefund`，退款回调仍接入 `RefundNotify`

**Why**：避免用户自助全额退与业务规则冲突；比例退款只在后台可控。

---

## Gotchas

1. **先标退款中再调微信**：否则核销可插队，导致退款比例/名额回退不一致。
2. **`sold` 含待支付占用**：超时关单必须释放名额（见 `server/initialize/timer.go`）。
3. **`pay.go` 直接查库**是历史公共入口模式；新业务优先把入账/退款成功逻辑放在插件 Service，API 薄封装调用。
