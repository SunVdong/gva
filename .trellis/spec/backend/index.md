# Backend Development Guidelines

> gin-vue-admin（GVA）后端约定。描述仓库**现状**，供实现/检查子代理对齐。

---

## Overview

后端根目录：`server/`。技术栈：Go + Gin + GORM + Zap + Casbin + Viper。

分层依赖必须单向：`Router → API → Service → Model`。API 禁止直接操作 DB；Service 禁止使用 `gin.Context`。

---

## Guidelines Index

| Guide | Description | Status |
|-------|-------------|--------|
| [Directory Structure](./directory-structure.md) | 目录布局、enter.go、插件结构 | Filled |
| [Database Guidelines](./database-guidelines.md) | GORM、迁移、事务、命名 | Filled |
| [Error Handling](./error-handling.md) | Service error → response 包 | Filled |
| [Logging Guidelines](./logging-guidelines.md) | `global.GVA_LOG` / zap | Filled |
| [Quality Guidelines](./quality-guidelines.md) | 分层、Swagger、禁止项 | Filled |
| [Mini Pay Integration](./mini-pay-integration.md) | 插件接入 `/mini/pay`、订单前缀、退款边界 | Filled |
| [Ticket Multi-Venue Verify](./ticket-multi-venue-verify.md) | 多次票多场合快照、公开核销 venue、列表筛选 | Filled |
| [Limited Activity Address](./limited-activity-address.md) | 活动选填 `address`、AutoMigrate、不入订单快照 | Filled |

---

## Golden Paths（优先对照）

| 主题 | 锚定文件 |
|------|----------|
| 基础模型 | `server/global/model.go` |
| 统一响应 | `server/model/common/response/response.go` |
| 核心 ApiGroup | `server/api/v1/enter.go`、`server/api/v1/system/enter.go` |
| 插件样板 | `server/plugin/announcement/`、`server/plugin/ticket/`、`server/plugin/limitedActivity/` |
| 事务与业务 | `server/plugin/ticket/service/order.go`、`server/plugin/limitedActivity/service/order.go` |
| 公共支付分发 | `server/api/v1/mini/pay.go` |

---

## Source of Conventions

- `.claude/rules/project_rules.md`
- 仓库真实代码（尤其 `server/plugin/ticket/`）
