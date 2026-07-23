# Type Safety

> 本仓库前端**没有 TypeScript**。类型安全靠约定、JSDoc 与前后端字段对齐。

---

## Reality Check

| 指标 | 现状 |
|------|------|
| `web/src` 下 `.ts` / `.tsx` | 0 |
| Vue `lang="ts"` | 无 |
| `tsconfig` / `vue-tsc` | 无 |
| 脚本语言 | JavaScript + `jsconfig.json` 路径别名 |

**不要**在未立项迁移的情况下把新文件写成 `.ts` 或 `lang="ts"`，以免与工具链脱节。

---

## Practical Rules

1. **API 封装**：参数与返回用 JSDoc 说明（`@param` / `@returns`）；可参考 `project_rules.md` 示例与带 Swagger 风格注释的 `web/src/api/*.js`
2. **与后端对齐**：
   - 成功码：`code === 0`
   - 分页字段：`page`、`pageSize`、`total`、`list`
   - 主键：`ID`（大写）
   - 同一业务字段前后端语义一致（number/string/boolean；后端指针 JSON 后可能是 `null`）
3. **Props**：`defineProps` 写清 `type` / `required` / `default`
4. **禁止**用 TypeScript 类型断言思维硬改运行时结构；以接口真实返回为准

---

## When You Need Stronger Typing

若未来引入 TS，应先统一工具链（`tsconfig`、构建、lint），再渐进迁移；在那之前本文件描述的 JS 约定仍是规范。

---

## Anti-patterns

- 假设响应一定有 `data.list` 却不判断 `code`
- 把后端 `ID` 当 `id` 用导致 row-key/详情查询失败
- 新增 `.ts` 文件但仓库无法类型检查
- model 与页面魔法字符串状态码不一致且无注释
