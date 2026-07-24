# Quality Guidelines

> 前端实现与审查清单。

---

## Architecture

1. 模块职责单一：页面不直接拼 URL；一律走 `api` + `request`
2. 依赖单向：`view/components → api → request → backend`
3. 全局状态经 Pinia actions；插件自包含在 `src/plugin/<name>/`

---

## API Layer

统一模式（核心与插件相同）：

```javascript
import service from '@/utils/request'

export const getOrderList = (params) =>
  service({ url: '/ticket/order/getOrderList', method: 'get', params })
```

- 文件：`web/src/api/*.js` 或 `web/src/plugin/<name>/api/*.js`
- 封装：`web/src/utils/request.js`（`x-token`、`x-user-id`，`code === 0`）
- 注释：JSDoc 或与后端一致的 Summary/Router 注释

参考：`web/src/plugin/ticket/api/order.js`、`web/src/api/user.js`。

---

## Naming

| 种类 | 约定 |
|------|------|
| 文件 | kebab-case 或与现有模块一致的小写名（`order.vue`、`useDict.js`） |
| 组件名 | PascalCase |
| 变量/函数 | camelCase |
| 常量 | UPPER_SNAKE_CASE |

以邻接文件为准，避免同一目录混用多种风格的新名字。

---

## Style & UX

- 优先 UnoCSS；遵循 Element Plus
- 禁止无必要的内联 style
- 主题用 CSS 变量（与现有 style 体系一致）
- 列表处理 loading / 空态 / 错误提示（`ElMessage` 等）
- 路由懒加载与现有 router 模式一致

---

## Collaboration with Backend

- 以后端 Swagger/注释为准更新前端 api
- 响应 `{code, data, msg}`；错误走 request 拦截与页面提示
- 字段类型与命名（含 `ID`）前后端一致

### H5 核销（`h5Verify.vue`）

- 按 URL `type` 分支处理；改 ticket 时**禁止**改动 camping / limitedActivity 业务语义
- 门票多场合：以 `order.supportMultiVenue`（订单快照）决定是否必选场合；`venue` code 须与后端 `server/plugin/ticket/model/venue.go` 白名单一致
- 公开核销：`code` 走 query，`venue` 走 body（见 `web/src/plugin/ticket/api/order.js`）
- 契约细节：`.trellis/spec/backend/ticket-multi-venue-verify.md`

---

## Forbidden Patterns

- 页面内裸 `fetch`/`axios.create` 绕过 `request.js`
- 提交前不处理 `code !== 0`
- 引入 TypeScript/Vuex 作为默认新栈（与现状不符）
- 破坏 `gva-*` 布局惯例导致与系统页面风格割裂（除非独立 H5）

---

## Verification Hints

- 在 `web/`：`npm run build` 或项目惯用的 lint/dev 命令
- 手动验证：登录态 header、列表分页、失败提示
- 插件菜单是否由后端注册并可打开对应 `view`
