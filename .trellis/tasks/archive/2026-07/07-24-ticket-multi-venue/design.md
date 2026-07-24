# Design: 门票多次票多场合核销

## Architecture / Boundaries

- 仅改 `ticket` 插件 + 共享 H5 `h5Verify.vue` 的 `type=ticket` 分支。
- 不改支付、退款、日历库存逻辑。
- 场合为代码内固定枚举，不建场合表。

## Data Model

### `ticket_sku`
| 字段 | 类型 | 说明 |
|------|------|------|
| `support_multi_venue` | bool / tinyint, default false | 是否支持多场合；仅多次票有意义 |

### `orders`（TicketOrder）
| 字段 | 类型 | 说明 |
|------|------|------|
| `support_multi_venue` | bool / tinyint, default false | 下单快照；核销以此为准 |

### `order_verify_records`
| 字段 | 类型 | 说明 |
|------|------|------|
| `venue` | string(32), default '' | 该次核销场合 code；不支持多场合时为空 |

### 场合枚举（存 code，UI 展示中文）

| code | label |
|------|-------|
| `zhongshanling` | 中山陵 |
| `zhaozhao` | 爪爪 |
| `lululand` | lululand |
| `hongshan` | 红山 |

后端集中定义校验函数（合法 code 白名单）；前后端 label 映射保持一致。

## Data Flow / Contracts

### 1. SKU 保存
- 管理端 `product.vue`：`ticketType===2` 时展示 checkbox「是否支持多场合」。
- `ticketType!==2` 时强制 `supportMultiVenue=false` 再提交（对齐 `useTimes` 处理）。

### 2. 下单快照
- `CreateOrder` 写订单时：`SupportMultiVenue = sku.TicketType==2 && sku.SupportMultiVenue`。

### 3. H5 查询
- `GET getOrderByCodePublic` 返回的 `order.supportMultiVenue` 驱动是否展示场合选择。
- 可选：同响应附带 `venueOptions: [{code,label}]`，避免前端写死；也可前端常量，与后端白名单一致即可。

### 4. H5 核销
- `POST verifyOrderByCodePublic`：`code` 仍可 query；**body JSON** 增加可选 `venue`（string）。
- 服务端：
  - `order.SupportMultiVenue==true`：`venue` 必填且 ∈ 白名单，写入本次 `OrderVerifyRecord.Venue`。
  - 否则：不要求 `venue`，记录 `Venue=""`（忽略传入值更安全，防止脏数据）。
- 签名扩展：`VerifyOrder(orderID, venue string)`；`VerifyOrderByOrderNoPublic(orderNo, venue)`。

### 5. 后台列表筛选
- `TicketOrderSearch` 增加：
  - `ticketType *int`：join `ticket_sku`，`ticket_sku.ticket_type = ?`（与现网列表展示票种一致，读当前 SKU）。
  - `venue string`：`EXISTS (SELECT 1 FROM order_verify_records r WHERE r.order_id = orders.id AND r.venue = ?)`。
- 列表 UI：筛选项「票种」「核销场合」；不加场合列。

### 6. 后台详情
- 已有 `verifyRecords`；表格增加「场合」列（空则显示 —）。

## Compatibility / Migration

- GORM AutoMigrate 加列即可；旧数据默认 `false` / `''`，行为与现网一致。
- 公开核销 API 对旧客户端：不传 `venue` 时，仅非多场合订单仍可核销；多场合订单返回明确错误文案（如「请选择核销场合」）。

## Trade-offs

| 决策 | 选择 | 原因 |
|------|------|------|
| 快照 vs 实时 SKU | 订单快照 | 已售票行为稳定 |
| 场合存 code | 稳定、可改文案 | label 仅展示 |
| 票种筛选用 join SKU | 不新增订单字段 | 与现网票种展示同源；SKU 改票种属运维边缘情况 |
| 核销 API 用 body 传 venue | 不破坏现有 query `code` | 与扩展字段兼容 |

## Rollback

- 回滚代码后新列可保留；忽略字段不影响旧逻辑。
- 若需强制关场合：将 SKU/订单 `support_multi_venue` 置 false（已产生的记录场合仍保留可查）。
