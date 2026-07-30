# Limited Activity Order Review

> 限时活动订单评价契约。对齐门票评价；独立表。2026-07-30 任务沉淀。

---

## Scenario: Activity order review create / delete / detail

### 1. Scope / Trigger

- Trigger：小程序活动订单评价 create/delete；后台/小程序订单详情聚合 `review`。
- 范围：`limited_activity_order_reviews` + mini 路由 + admin `findOrder` / `order.vue`。
- 非范围：评价列表管理、审核回复、活动页聚合分；勿复用门票 `order_reviews`。

### 2. Signatures

**DB / Model**（`server/plugin/limitedActivity/model/order_review.go`）：

| 列 | JSON | 约束 |
|----|------|------|
| order_id | orderId | `uniqueIndex:idx_la_order_review` |
| user_id | userId | 评价用户 |
| rating | rating | 1–5 |
| content | content | ≤50 字 |

表名：`limited_activity_order_reviews`。迁移：插件 `initialize/gorm.go` AutoMigrate。

**API**

| 方法 | 路径 | 鉴权 |
|------|------|------|
| POST | `/limitedActivity/mini/order/review/create` | JWT |
| POST | `/limitedActivity/mini/order/review/delete` | JWT（`id` JSON 或 query） |
| GET | `/limitedActivity/mini/order/detail` | JWT，已核销附 `review` |
| GET | `/limitedActivity/order/findOrder` | 后台，已核销附 `review` |

**Create body**

```json
{ "orderId": 1, "rating": 5, "content": "可选，50字内" }
```

### 3. Contracts

- `review` 载荷：`{ ID, rating, content, createdAt }`；无评价为 `null`；非 `status==2` 可不返回该字段（对齐门票）。
- 评价条件：`order.userId == JWT user` 且 `status==2` 且 `verifiedAt != nil`；一单一评。
- **删除必须硬删**：`Unscoped().Delete`。MySQL 下唯一索引包含软删行，软删会导致「删后重评」失败。

### 4. Validation & Error Matrix

| 条件 | 文案 |
|------|------|
| 订单不存在 | 订单不存在 |
| 非本人 | 无权对该订单评价 |
| status ≠ 2 | 仅已核销订单可评价 |
| verifiedAt 空 | 仅核销后的订单可评价 |
| 已有评价 | 该订单已评价过，可先删除再重新评价 |
| 删非本人评价 | 无权删除该评价 |

### 5. Good / Base / Bad Cases

- **Good**：已核销本人单 → create → detail/后台见 review → delete → 再 create。
- **Base**：核销中（status=1）不可评。
- **Bad**：软删评价（唯一索引挡住重评）；复用 `order_reviews`（与门票 order_id 撞车）。

### 6. Tests Required

- create / 非本人 / 未核销 / 重复评价 / delete 后 recreate。
- findOrder 与 mini detail 在 status=2 时 `review` 形状一致。

### 7. Wrong vs Correct

#### Wrong

```go
return global.GVA_DB.Delete(&review).Error // 软删 → 唯一索引挡重评
```

#### Correct

```go
return global.GVA_DB.Unscoped().Delete(&review).Error
```

独立表 `limited_activity_order_reviews`，勿加 type 复用门票评价表。

---

## Design Decisions

| 决策 | 理由 |
|------|------|
| 独立表而非 type 字段 | 与露营分表一致；避免跨业务 order_id 碰撞 |
| 删除硬删 | MySQL unique + 软删无法删后重评 |
| 规则对齐门票 | 产品确认；文案与 status=2 门槛一致 |
