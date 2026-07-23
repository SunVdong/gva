# Component Guidelines

> Vue SFC 与页面组件约定。

---

## Overview

- 全项目 `.vue` 使用 **`<script setup>`**（Composition API）
- **无** `lang="ts"`；脚本为 JavaScript
- UI：Element Plus；布局 class 常见 `gva-search-box`、`gva-table-box`
- 样式优先 UnoCSS 原子类 + Element Plus；避免内联 style

参考页面：`web/src/plugin/ticket/view/order.vue`。

---

## Page Pattern（列表页）

典型管理端列表：

1. 搜索区：`el-form` + 查询/重置
2. 表格区：`el-table`，`row-key="ID"`（对齐后端 JSON `ID`）
3. 分页：`page` / `pageSize` / `total`
4. 状态：`ref`/`reactive` 管理 loading、list、form
5. 数据：从 `plugin/.../api` 或 `@/api` 拉数；成功看 `res.code === 0`

---

## Reusable Components

- 目录：`web/src/components/`，按功能分子目录
- 使用 `defineProps` / `defineEmits`
- 可复用 UI 抽组件；页面只做业务编排

`project_rules.md` 中的 props/emits 示例仍适用（JS 版即可）。

---

## Composition Rules

- 页面：`view/` 或 `plugin/*/view/`
- 全局能力：优先 hooks / pinia，而非在多个页面复制
- 跨组件通知：已有 `@/utils/bus`（如 request 错误展示），勿再引入第二套事件总线

---

## Anti-patterns

- 新页面改回 Options API `export default { data() ... }` 作为主模式
- 忽略 `res.code`，把 Axios 原始错误当业务成功
- 表格 `row-key` 写成 `id` 小写（与后端 `ID` 不一致）
- 大段重复表格/搜索逻辑不抽组件或 composable
