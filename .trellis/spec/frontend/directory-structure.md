# Directory Structure

> 前端代码组织方式。

---

## Overview

应用在 `web/src/`。业务分：

1. **核心**：`api/`、`view/`、`components/`、`pinia/`、`hooks/`
2. **插件**：`web/src/plugin/<name>/`（api + view，可选 components/form）

---

## `src/` Layout

```
web/src/
├── api/              # 核心后端接口封装（全 .js）
├── assets/
├── components/       # 通用组件
├── core/             # gin-vue-admin 核心配置
├── directive/
├── hooks/            # 组合式 hooks
├── pinia/
│   ├── index.js
│   └── modules/      # user/app/router/dictionary/params
├── plugin/           # 前端插件
├── router/           # index.js（动态路由为主）
├── style/
├── utils/            # request.js、dictionary、bus 等
├── view/             # 核心业务页面
├── App.vue
└── main.js
```

路径别名：`@/*` → `src/*`（`jsconfig.json`）。

---

## Plugin Layout

```
web/src/plugin/<name>/
├── api/           # *.js，走 @/utils/request
├── view/          # 业务页面 *.vue
├── components/    # 可选
├── form/          # 可选（如 announcement）
└── config.js      # 可选
```

现有插件示例：`ticket`、`camping`、`announcement`、`email`、`activityGuide`、`limitedActivity`。

Ticket 现状（精简）：

- `web/src/plugin/ticket/api/{order,product,scenic,user}.js`
- `web/src/plugin/ticket/view/{order,product,scenic,user,calendar}.vue`

LimitedActivity：

- `web/src/plugin/limitedActivity/api/{activity,order}.js`
- `web/src/plugin/limitedActivity/view/{activity,order}.vue`

插件一般**不**自建 Pinia；菜单路由多由后端动态下发。硬编码路由例外少见（如 camping H5 在 `web/src/router/index.js`）。

### H5 核销页约定

- 路由：`#/h5/verify?type=<biz>&code=<核销码或订单号>`，组件：`web/src/plugin/camping/view/h5Verify.vue`
- 已支持 `type`：`reservation` | `ticket` | `limitedActivity`
- 新增业务：在该页增加分支 + 对应插件 `*Public` API；**勿**破坏既有 type 行为；已鉴权后 query 变化应重新加载详情

---

## Where New Code Goes

| 变更类型 | 放置位置 |
|----------|----------|
| 系统级接口 | `web/src/api/` |
| 系统级页面 | `web/src/view/<module>/` |
| 可复用 UI | `web/src/components/` |
| 插件功能 | `web/src/plugin/<name>/` |
| 全局状态 | `web/src/pinia/modules/` |
| 可复用组合逻辑 | `web/src/hooks/` |

---

## Anti-patterns

- 页面里直接 `axios`，绕过 `@/utils/request` 与 `api/`
- 插件页面去改无关核心模块造成耦合
- 在 `components/` 塞仅某插件使用的页面级大文件（应放 plugin 内）
