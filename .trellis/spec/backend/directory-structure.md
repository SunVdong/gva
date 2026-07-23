# Directory Structure

> 后端代码如何组织。新增模块/插件时对齐此布局。

---

## Overview

GVA 采用前后端分离；后端在 `server/`。业务分两类：

1. **核心模块**：`api/v1`、`service`、`router`、`model` 下按域分包（`system`、`example`、`feedback`、`mini` 等）
2. **插件**：`server/plugin/<name>/`，自包含 api/service/router/model/initialize，经 plugin v2 注册

---

## Top-level Layout

```
server/
├── api/v1/           # HTTP 控制器（按域）
├── config/           # 配置结构体
├── core/             # 启动与核心引导
├── docs/             # Swagger
├── global/           # GVA_DB / GVA_LOG / GVA_MODEL
├── initialize/       # DB、路由、插件、Redis 等初始化
├── middleware/       # JWT、Casbin、操作记录等
├── model/            # 领域模型 + common/request|response
├── plugin/           # 业务插件（ticket、camping、announcement…）
├── router/           # 路由注册
├── service/          # 业务逻辑
├── source/           # 初始化数据
├── utils/            # 工具 + plugin/v2 接口
├── mcp/              # MCP 辅助
└── main.go
```

---

## Layer Responsibilities

| 层 | 职责 | 禁止 |
|----|------|------|
| Model | GORM 结构体、request DTO | 业务编排、HTTP |
| Service | CRUD 与业务规则，返回 `(result, error)` | `gin.Context`、写 HTTP 响应 |
| API | 绑定参数、调 Service、`response.*` | 直接 `GVA_DB` |
| Router | 路径、中间件、映射到 Api 方法 | 业务逻辑 |
| Initialize（插件） | AutoMigrate、挂路由、菜单/API 注册 | 散落在随意文件 |

依赖链：`Router → API → Service → Model`。

---

## `enter.go` Group Pattern

核心与插件都用组入口避免循环引用。

**核心三层聚合：**

| 层 | 入口 | 全局变量 |
|----|------|----------|
| Service | `server/service/enter.go` | `ServiceGroupApp` |
| API | `server/api/v1/enter.go` | `ApiGroupApp` |
| Router | `server/router/enter.go` | `RouterGroupApp` |

系统域再拆一层，例如：

- `server/api/v1/system/enter.go`：`var xxxService = service.ServiceGroupApp.SystemServiceGroup.Xxx`
- `server/router/system/enter.go`：`var xxxApi = api.ApiGroupApp.SystemApiGroup.XxxApi`

**插件（ticket 现状）** 使用小写结构体 + 包级变量：

- `server/plugin/ticket/api/enter.go` → `var Api`
- `server/plugin/ticket/service/enter.go` → `var Service`
- `server/plugin/ticket/router/enter.go` → `var Router`

新插件优先对齐已有插件风格；核心新模块对齐 `system` 的 `XxxGroupApp` 风格。

---

## Plugin Layout

参考：`server/plugin/ticket/`、`server/plugin/announcement/`。

```
server/plugin/<name>/
├── plugin.go              # 实现 plugin v2，interfaces.Register
├── api/
│   ├── enter.go
│   ├── *.go               # 后台 API
│   └── mini/              # 可选：小程序/公开端 API
├── service/
│   ├── enter.go
│   └── *.go
├── router/
│   ├── enter.go
│   └── *.go
├── model/
│   ├── *.go
│   └── request/
├── initialize/
│   ├── gorm.go            # AutoMigrate
│   ├── router.go          # public/private 分组
│   ├── menu.go
│   └── api.go
└── config/                # 可选（announcement 有，ticket 无）
```

插件入口要点（见 `server/plugin/ticket/plugin.go`）：

- 实现 `server/utils/plugin/v2` 的 `Plugin` 接口
- `init` 中 `interfaces.Register(Plugin)`
- `Register` 内调用 `initialize.Api/Menu/Gorm/Router`
- 本体用匿名 import 激活：`server/plugin/register.go` 中 `_ ".../plugin/<name>"`

路由：`initialize/router.go` 通常拆 `public`（可匿名）与 `private`（JWT + Casbin）；写操作常挂 `middleware.OperationRecord()`。

经典样板：`server/plugin/announcement/`。

---

## Module Wiring Cheatsheet

- API → Service：`var xxxService = service.ServiceGroupApp....` 或插件包内 `Service.Xxx`
- Router → API：`api.ApiGroupApp....` 或插件 `Api.Xxx`
- Initialize → Router：`router.RouterGroupApp....` 或插件 `Router.Init(...)`

---

## Anti-patterns

- API 直接 `global.GVA_DB` 查改数据
- Service 依赖 `*gin.Context`
- 新文件散落在错误层目录却不在对应 `enter.go` 注册
- 插件不经 `interfaces.Register` / `register.go` 匿名 import，导致未加载
