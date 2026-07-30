# Design: 门票订单按核销场合统计

## Scope / Boundaries

- 改动面：`server/plugin/ticket/`（API / Service / request / router / initialize/api）+ `web/src/plugin/ticket/`（`api/order.js`、`view/order.vue`）
- 不改：H5 核销、SKU、`order_verify_records` 表结构、camping / limitedActivity
- 数据源：现有 `order_verify_records`（`venue` + `verified_at`）

## Architecture

```
order.vue 汇总条
  → GET /ticket/order/getVenueVerifyStats?month=YYYY-MM
  → api.GetVenueVerifyStats
  → service.GetVenueVerifyStats(month)
  → GROUP BY venue on order_verify_records
  → 补齐白名单场合缺项为 0
```

## Contracts

### Request

| 字段 | 位置 | 约束 |
|------|------|------|
| `month` | query | 可选，格式 `YYYY-MM`；缺省或空 → 服务端用当前自然月 |

非法 `month`（无法解析）→ 返回错误信息，不静默回退。

### Response `data`

```json
{
  "month": "2026-07",
  "items": [
    { "venue": "zhongshanling", "label": "中山陵", "count": 12 },
    { "venue": "zhaozhao", "label": "爪爪", "count": 3 },
    { "venue": "lululand", "label": "lululand", "count": 0 },
    { "venue": "hongshan", "label": "红山", "count": 5 }
  ]
}
```

- `items` 顺序与 `model.VenueOptions()` 一致。
- 只含白名单四场合；SQL 侧 `venue IN (...)` 且 `venue <> ''`，空 venue / 未知 code 不返回。

### Service 查询要点

- 时间窗：`verified_at >= monthStart AND verified_at < nextMonthStart`（半开区间，避免月末边界漏计）。
- `deleted_at IS NULL`（软删记录不计）。
- `GROUP BY venue` + `COUNT(*)`。
- 结果映射到固定四场合，缺省补 0。

### API / 权限

- `GET /ticket/order/getVenueVerifyStats`，挂 private order 组（与 `getOrderList` 同级，无需 OperationRecord）。
- `initialize/api.go` 注册新 API 描述，便于权限同步。

## Frontend

- `order.js` 新增 `getVenueVerifyStats({ month })`。
- `order.vue`：搜索区与表格之间增加汇总条：
  - `el-date-picker` type=`month`，默认当前月，绑定 `statsMonth`
  - 四个场合计数展示（可用简单横向标签 / 数字，避免卡片堆砌）
  - `onMounted` + 月份变更时拉取统计；与 `getTableData` 独立调用

## Compatibility / Rollback

- 无 DB 迁移；回滚删 API + 前端汇总条即可。
- 历史空 venue 数据不影响展示（被忽略）。

## Trade-offs

| 选择 | 原因 |
|------|------|
| 按月而非任意日期 | 产品明确；参数更简单 |
| 统计与列表筛选解耦 | 列表是订单维度，统计是核销记录维度 |
| 服务端补齐 0 | 前端不必硬编码缺项逻辑，契约稳定 |
