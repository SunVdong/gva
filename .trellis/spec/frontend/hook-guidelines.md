# Hook Guidelines

> `web/src/hooks` 中的组合式函数。

---

## Overview

Hooks 封装可复用的 Composition 逻辑（字典、响应式布局、图表 option 等）。页面通过 import 使用，而不是复制粘贴。

---

## Existing Hooks（现状）

| 文件 | 导出 | 用途 |
|------|------|------|
| `web/src/hooks/useDict.js` | 命名导出 `useDict` / `useMultiDict` | 字典 |
| `web/src/hooks/use-windows-resize.js` | （见文件） | 窗口尺寸 |
| `web/src/hooks/responsive.js` | `export default function useResponsive` | 响应式 |
| `web/src/hooks/charts.js` | `export default function useChartOption` | ECharts option |

命名**不统一**（`useXxx.js` / `use-kebab.js` / 无 use 前缀文件名）。这是历史现状。

---

## Conventions for New Hooks

1. **优先**命名导出 + `use` 前缀：`export function useSomething(...)`
2. 文件名倾向：`useSomething.js`（与 `useDict.js` 对齐）
3. 只放可复用逻辑；单页一次性逻辑留在 SFC
4. 需要全局持久状态时用 Pinia，不要用 hook 假装全局 store
5. 字典类优先复用 `useDict` + `pinia/modules/dictionary.js` + `utils/dictionary.js`

---

## Anti-patterns

- 在 hook 里直接调 Element Plus 全局消息又隐式依赖路由，导致难以测试/复用（除非现有同类 hook 已如此）
- 新建与 `useDict` 重复的字典加载逻辑
- 把整页业务（含表格请求与权限）塞进一个巨型 hook
