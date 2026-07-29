# Limited Activity Address Field

> 限时活动 `Activity.address`（活动地点）字段契约。2026-07-29 任务沉淀。

---

## Scenario: Activity optional address

### 1. Scope / Trigger

- Trigger：活动实体增删改查 / 小程序活动列表与详情响应变更；DB 经插件 AutoMigrate 加列。
- 范围：`limited_activities.address` + 管理端活动表单/列表。
- 非范围：订单快照、核销 H5、地图/经纬度、地点搜索。

### 2. Signatures

**DB / Model**（`server/plugin/limitedActivity/model/activity.go`）：

```go
Address string `json:"address" form:"address" gorm:"column:address;comment:活动地点;size:256;"`
```

- 表：`limited_activities`（`TableName()`）
- 迁移：插件 `initialize/gorm.go` → `AutoMigrate(Activity)`，无手工 SQL

**Service Update map**（`service/activity.go`）须显式包含：

```go
"address": m.Address
```

Create 走整模 `Create`；Create/Update 均对 `Address` 做 `strings.TrimSpace`。

**Admin UI**：`web/src/plugin/limitedActivity/view/activity.vue` — 表单选填、列表空显示 `-`、payload 带 `address`。

### 3. Contracts

| 方向 | 字段 | 类型 | 约束 |
|------|------|------|------|
| Request/Response | `address` | string | 选填；最长约 256；空串合法 |
| JSON | camelCase `address` | — | 勿用 `location` |
| 订单 | — | — | **不**快照 address；订单仍仅 `activityName` |

管理端与小程序活动 List/Detail/Create/Update 凡返回 `model.Activity` 即带 `address`。

### 4. Validation & Error Matrix

| 条件 | 行为 |
|------|------|
| `address` 缺省 / `""` | 允许保存 |
| 仅空白 | Trim 后按空串存储 |
| 名称空等既有校验 | 与 address 无关，仍按原规则报错 |

### 5. Good / Base / Bad Cases

- **Good**：填写「景区东门集合点」，保存后列表与详情回显。
- **Base**：不填 address，旧活动照常编辑保存，JSON 为 `""`。
- **Bad**：字段命名成 `location` 或 Update map 漏写 `address`（编辑无法落库）。

### 6. Tests Required

- 手工/接口：Create 带 `address` → Find 回显；Update 清空 → 为空串。
- Update 后 DB 列 `address` 与 JSON 一致（防 Updates map 漏字段）。
- 回归：订单创建/核销路径不依赖 address。

### 7. Wrong vs Correct

#### Wrong

```go
// Update 漏 address → 前端改了也不落库
updates := map[string]any{"name": m.Name, "detail": m.Detail}
```

```json
{ "location": "东门" }
```

#### Correct

```go
updates := map[string]any{
  "name":    m.Name,
  "address": m.Address,
  // ...
}
```

```json
{ "address": "东门" }
```

---

## Design Decision: 选填且不入订单快照

**Context**：需展示活动地点，但历史活动与订单模型要少侵入。

**Decision**：活动侧选填文本 `address`；订单不写快照，需要时按 `activityId` 读当前活动。

**Why**：与现有「仅快照活动名」一致；改址后订单侧看到最新地点即可满足本轮需求。
