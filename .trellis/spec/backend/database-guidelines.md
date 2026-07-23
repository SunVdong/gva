# Database Guidelines

> GORM、迁移、查询与命名——以仓库现状为准。

---

## Overview

- ORM：GORM
- 全局 DB：`global.GVA_DB`（多库：`GVA_DBList` + `GetGlobalDBByDBName` / `MustGetGlobalDBByDBName`）
- 初始化：`server/initialize/gorm.go` 及各驱动文件
- 插件表：各自 `initialize/gorm.go` 里 `AutoMigrate`

---

## Base Model

所有持久化实体嵌入 `global.GVA_MODEL`（`server/global/model.go`）：

- `ID uint`，JSON 字段名为 **`ID`（大写）**
- `CreatedAt` / `UpdatedAt` / `DeletedAt`（软删）

示例：`server/model/system/sys_user.go`、`server/plugin/ticket/model/order.go`。

字段须有清晰的 `json` 与 `gorm` 标签。跨 model/request/response 的**同一业务字段类型必须一致**（尤其状态、ID、枚举、时间）。数据模型用指针（如 `*string`）而 DTO 用值类型时，在 Service 做 nil 安全转换。

---

## Request DTOs

- 路径：`model/request/` 或插件 `plugin/<name>/model/request/`
- 绑定标签：`json` + `form`（Gin）
- 列表查询：搜索结构体嵌入 `request.PageInfo`（`server/model/common/request/common.go`），可用 `Paginate()` scope

示例：`server/plugin/ticket/model/request/order.go` → `TicketOrderSearch`。

---

## Migrations

- 核心与插件均以 **`AutoMigrate`** 为主（插件：`server/plugin/ticket/initialize/gorm.go`）
- 个别手工 SQL 可并存（如 `server/plugin/ticket/migrate_order_items.sql`），新变更优先跟现有 AutoMigrate 路径，避免 silently 只改 SQL

---

## Transactions & Queries

多步写操作使用：

```go
global.GVA_DB.Transaction(func(tx *gorm.DB) error {
    // 全部用 tx，不要混用 global.GVA_DB
    return nil
})
```

参考：`server/service/system/sys_user.go`、`server/plugin/ticket/service/order.go`。

库存/订单等并发敏感路径可使用行锁：

`tx.Clauses(clause.Locking{Strength: "UPDATE"})`（见 ticket order service）。

虚拟展示字段用 `gorm:"-"`（如 order 上的 `ProductName`），不落库。

---

## Naming Reality

| 区域 | 现状 |
|------|------|
| 系统表 | 常 `sys_` + 复数，如 `sys_users`（`TableName()`） |
| Ticket 等插件 | 混用：`orders`、`ticket_sku` 等，以各 model 的 `TableName()` 为准 |
| JSON | camelCase（如 `orderNo`）；主键 JSON 为 `ID` |
| DB column | 插件里常见显式 `column:snake_case` |

**新表**：同插件内保持一致；系统域继续 `sys_` 风格。不要为了“理想统一”去改已有表名。

业务状态多为 `int` 魔法数字 + model 注释说明（如订单 0–7），新增状态先读现有注释再扩展。

---

## Anti-patterns

- 事务回调内仍写 `global.GVA_DB`（应用 `tx`）
- model 与 request 同名字段类型不一致
- 漏写 `TableName()` 导致表名与既有数据不一致
- 在 API 层拼装复杂查询绕过 Service
