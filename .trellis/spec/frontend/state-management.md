# State Management

> 全局状态使用 Pinia；页面局部状态留在组件内。

---

## Overview

- 库：Pinia（**无 Vuex**）
- 入口：`web/src/pinia/index.js`（`createPinia()` + 再导出）
- 模式：**setup store** — `defineStore('id', () => { ... return { ... } })`

---

## Existing Stores

| 文件 | id |
|------|-----|
| `web/src/pinia/modules/user.js` | `user` |
| `web/src/pinia/modules/app.js` | `app` |
| `web/src/pinia/modules/router.js` | `router` |
| `web/src/pinia/modules/dictionary.js` | `dictionary` |
| `web/src/pinia/modules/params.js` | `params` |

`user.js` 特征：

- `ref` / `computed` 定义状态
- `useStorage('token')` + cookie `x-token`
- 调 `@/api/user`；成功条件 `res.code === 0`
- 通过函数（actions）修改状态，再在 return 中暴露

---

## What Goes Where

| 数据 | 放置 |
|------|------|
| 登录用户、token、权限 | `user` store |
| 布局/主题等应用配置 | `app` store |
| 动态路由 | `router` store |
| 字典缓存 | `dictionary` + `useDict` |
| 单页表格/表单/loading | 组件内 `ref`/`reactive` |
| 插件业务临时状态 | 默认组件内；仅跨页全局才新建 pinia module |

组件**不要**直接改写其他模块的 store 内部 ref；走 store 暴露的方法。

---

## Persistence & Side Channels

- Token：`useStorage` + cookies（与 `request.js` 的 `x-token` / `x-user-id` header 配合）
- 事件总线：`@/utils/bus`（非状态源，仅通知）

---

## Anti-patterns

- 引入 Vuex 或第二套全局状态库
- 在组件里 `userStore.token = xxx` 绕过 `setToken`
- 为每个插件页面建 Pinia module（通常不必要）
- Options 风格 `state: () => ({})` 的新 store（与现有 setup 风格不一致）
