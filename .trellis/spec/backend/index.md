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

---

## Golden Paths（优先对照）

| 主题 | 锚定文件 |
|------|----------|
| 基础模型 | `server/global/model.go` |
| 统一响应 | `server/model/common/response/response.go` |
| 核心 ApiGroup | `server/api/v1/enter.go`、`server/api/v1/system/enter.go` |
| 插件样板 | `server/plugin/announcement/`、`server/plugin/ticket/` |
| 事务与业务 | `server/plugin/ticket/service/order.go` |

---

## Source of Conventions

- `.claude/rules/project_rules.md`
- 仓库真实代码（尤其 `server/plugin/ticket/`）
