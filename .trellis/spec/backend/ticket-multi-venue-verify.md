# Ticket Multi-Venue Verify

> 门票多次票「多场合」核销契约。`ticket-multi-venue` 任务落地后沉淀。

---

## Scenario: 多次票多场合核销

### 1. Scope / Trigger

- Trigger: 门票核销需要按点位记录场合，或扩展公开核销 API / 订单筛选
- 改动面：`server/plugin/ticket/` + 管理端 `web/src/plugin/ticket/` + 共享 H5 `h5Verify.vue` 的 **`type=ticket` 分支 only**
- 勿把场合逻辑并入 camping / limitedActivity；勿实时读 SKU 决定是否必选场合

### 2. Signatures

| 层 | 签名 / 表字段 |
|----|----------------|
| DB `ticket_sku` | `support_multi_venue` bool default false（仅 `ticket_type=2` 有意义） |
| DB `orders` | `support_multi_venue` bool default false（**下单快照**） |
| DB `order_verify_records` | `venue` string(32) default ''（按次；非多场合为空） |
| Service | `VerifyOrder(orderID uint, venue string) error` |
| Service | `VerifyOrderByOrderNoPublic(orderNo, venue string) error` |
| API | `POST /ticket/order/verifyOrderByCodePublic?code=` + body `{ "venue": "<code>" }` |
| 场合白名单 | `model/venue.go`：`zhongshanling` / `zhaozhao` / `lululand` / `hongshan` |

下单快照（`CreateOrder`）：

```go
SupportMultiVenue: sku.TicketType == 2 && sku.SupportMultiVenue
```

列表筛选（`TicketOrderSearch`）：

- `ticketType`：`EXISTS` join 当前 `ticket_sku.ticket_type`
- `venue`：`EXISTS` 核销记录 `r.venue = ?`（至少一次匹配）

### 3. Contracts

**公开核销请求**

| 字段 | 位置 | 约束 |
|------|------|------|
| code | query | 必填，订单号 |
| venue | body JSON | 当 `order.supportMultiVenue==true` 时必填且 ∈ 白名单；否则可省略，服务端忽略并写空 |

**公开查询**：`GET getOrderByCodePublic` 返回的 `order.supportMultiVenue` 驱动 H5 是否展示场合选择。

**管理端展示**：详情 `verifyRecords` 展示场合；列表**不加**场合列，仅筛选项。

**前端场合 options**：可与后端重复常量，但 **code 必须与 `venue.go` 一致**。

### 4. Validation & Error Matrix

| 条件 | 行为 / 错误 |
|------|-------------|
| 订单非待核销 / 次数用尽 | 既有核销错误（与多场合无关） |
| `supportMultiVenue==true` 且 venue 空或非法 | `请选择核销场合` |
| `supportMultiVenue==false` | 不校验 venue；记录 `Venue=""` |
| SKU `ticketType!=2` 保存 | 强制 `SupportMultiVenue=false` |
| 历史订单无快照列默认 | `false`，行为与现网一致 |

### 5. Good / Base / Bad Cases

- **Good**：多次票开多场合 → 下单快照 true → H5 选 `hongshan` 核销成功 → 记录含 venue；二次可选其他场合
- **Base**：多次票未开 / 单次票 → H5 无场合控件；不传 venue 可核销
- **Bad**：多场合订单不传 venue / 传未知 code → 核销失败；改 SKU 开关不应改变已下单订单的必选行为

### 6. Tests Required

- 单测/手工：快照 true 无 venue 失败；合法 venue 写入记录
- 快照 false 传入脏 venue → 记录仍为空
- 列表：`ticketType=2`、`venue=hongshan` 仅命中至少一次该场合的订单
- H5：确认 camping / limitedActivity 分支不受影响

### 7. Wrong vs Correct

#### Wrong

```go
// 核销时实时读 SKU，导致运营改开关影响已售票
sku := loadSku(order.SkuID)
if sku.SupportMultiVenue { ... }
```

#### Correct

```go
// 以订单下单快照为准
if order.SupportMultiVenue {
    if !model.IsValidVenue(venue) {
        return fmt.Errorf("请选择核销场合")
    }
}
```

---

## Design Decisions

### 下单快照 vs 实时 SKU

**Decision**：订单字段快照。改 SKU 只影响新单。

### 场合存 code 不存中文

**Decision**：稳定枚举；label 仅展示。扩展场合时先改 `venue.go`，再对齐管理端 / H5 常量。

---

## Scenario: 管理端按核销场合月度统计

> `ticket-venue-stats` 任务落地。复用 `order_verify_records`，无新表。

### 1. Scope / Trigger

- Trigger: 订单管理页需要按场合看某月核销次数
- 改动面：`server/plugin/ticket/` + `web/src/plugin/ticket/view/order.vue`（列表上方汇总条）
- 与订单列表筛选解耦；勿统计空 `venue`

### 2. Signatures

| 层 | 签名 |
|----|------|
| Service | `GetVenueVerifyStats(month string) (monthOut string, items []TicketVenueVerifyStatsItem, err error)` |
| API | `GET /ticket/order/getVenueVerifyStats?month=YYYY-MM` |
| Request | `month` query 可选；空 → 当前自然月 |
| Response item | `{ venue, label, count }`；顺序 = `VenueOptions()` |

### 3. Contracts

- 时间：`verified_at >= monthStart AND verified_at < nextMonthStart`（半开区间）
- 过滤：`venue <> '' AND venue IN (白名单)`；软删不计
- 缺项补 `count=0`；不返回「未指定场合」
- 前端：月份选择器默认当前月、不可清空；与列表 `searchInfo` 独立

### 4. Validation & Error Matrix

| 条件 | 行为 |
|------|------|
| `month` 空 | 用当前自然月 |
| `month` 非法（非 `YYYY-MM`） | 错误「月份格式错误，应为 YYYY-MM」，不静默回退 |

### 5. Good / Base / Bad Cases

- **Good**：当月有 `hongshan` 核销 → items 中红山 count>0，其它场合可为 0
- **Base**：某月无任何白名单核销 → 四场合均为 0
- **Bad**：传 `2026-13` → 接口失败，不返回空统计冒充成功

### 6. Tests Required

- 手工：默认当前月；切换月份数字变化；空 venue 记录不计入
- 改列表筛选时汇总不变（直至改月份）

### 7. Wrong vs Correct

#### Wrong

```go
// 闭区间月末 23:59:59，易漏跨秒 / 时区边界
db.Where("verified_at BETWEEN ? AND ?", start, endOfMonth)
```

#### Correct

```go
db.Where("verified_at >= ? AND verified_at < ?", monthStart, monthStart.AddDate(0, 1, 0))
```
