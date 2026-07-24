# Frontend Development Guidelines

> gin-vue-admin 前端约定。仓库现状为 **Vue 3 + JavaScript（无 TypeScript）**。

---

## Overview

前端根目录：`web/`。技术栈：Vue 3.5 + Vite + Pinia + Element Plus + UnoCSS + Vue Router + Axios + VueUse。

依赖链：`页面/组件 → api 封装 → @/utils/request → 后端`。

---

## Guidelines Index

| Guide | Description | Status |
|-------|-------------|--------|
| [Directory Structure](./directory-structure.md) | `src/` 与插件目录 | Filled |
| [Component Guidelines](./component-guidelines.md) | SFC、`script setup`、Element Plus | Filled |
| [Hook Guidelines](./hook-guidelines.md) | `src/hooks` 组合式函数 | Filled |
| [State Management](./state-management.md) | Pinia setup store | Filled |
| [Type Safety](./type-safety.md) | JS + JSDoc / 字段一致性（非 TS） | Filled |
| [Quality Guidelines](./quality-guidelines.md) | 命名、样式、API 注释 | Filled |

---

## Golden Paths

| 主题 | 锚定文件 |
|------|----------|
| Axios 封装 | `web/src/utils/request.js` |
| 插件 API | `web/src/plugin/ticket/api/order.js` |
| 插件页面 | `web/src/plugin/ticket/view/order.vue`、`web/src/plugin/limitedActivity/view/` |
| H5 核销 | `web/src/plugin/camping/view/h5Verify.vue`（多 type） |
| Pinia | `web/src/pinia/modules/user.js` |
| Hooks | `web/src/hooks/useDict.js` |

---

## Source of Conventions

- `.claude/rules/project_rules.md` 前端章节
- `web/src` 真实代码（全量 `.vue` 为 `<script setup>`，无 `.ts`）
