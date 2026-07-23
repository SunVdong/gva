# Quality Guidelines

> 后端代码审查与实现时必须守住的约定。

---

## Layering（不可违背）

1. `Router → API → Service → Model` 单向依赖
2. API 只做绑定、鉴权上下文取值、调 Service、写 `response`
3. Service 无 HTTP；无 `gin.Context`
4. 模块间经 `enter.go` 暴露的 Group / 包级 `Api`/`Service`/`Router` 通信，避免循环 import

权威叙述：`.claude/rules/project_rules.md` 后端规则章节。

---

## Swagger（API 强制）

每个对外 API 方法必须有完整注释块：`@Tags`、`@Summary`、`@Security`（如需）、`@accept`/`@Produce`、`@Param`、`@Success`、`@Router`。

示例形态见 `project_rules.md` 与现有 `server/api/v1/system/*.go`。注释是文档与前后端协作的来源，不可省略。

---

## Type Consistency

- 同一字段在 model / request / response / 前端用法上类型一致
- 指针与非指针转换在 Service 显式处理 nil
- JSON 主键为 `ID`（大写）；前端表格常用 `row-key="ID"`

---

## Plugin Quality Bar

新插件应对齐：

1. Model（含 request）→ Service → API → Router → initialize → `plugin.go`
2. `interfaces.Register` + `register.go` 匿名 import
3. AutoMigrate 覆盖本插件全部表模型
4. 菜单/API 权限初始化与路由分组（public/private）正确

推荐对照：`server/plugin/announcement/`（完整）、`server/plugin/ticket/`（精简但生产在用）。

---

## Cross-cutting Collaboration

- 响应：`{code, data, msg}`；分页：`PageResult` / 前端约定 `page, pageSize, total, list`
- 时间：与现有接口保持一致（勿擅自换格式）
- 接口变更先改 Swagger/注释，再改前端 `src/api` 或 `plugin/*/api`

---

## Forbidden Patterns

- 跨层调用（API→DB，Router→Service，Service→`gin.Context`）
- 新接口无 Swagger
- 自定义成功码（非 0）或抛弃统一 `response` 包
- 提交含密钥的 `config.yaml` 本地改动到公共分支（注意仓库已有配置文件惯例，勿把真实密钥写进 spec 示例）

---

## Verification Hints

- 编译：在 `server/` 下 `go build ./...`
- 确认新类型已加入对应 `enter.go`
- 插件：确认 `plugin/register.go` 已匿名 import
- 手动走一遍：绑定失败、业务失败、成功分页三条路径
