# Error Handling

> Service 返回 `error`，API 转为统一 JSON 响应。

---

## Overview

统一响应包：`server/model/common/response/`（`response.go`、`common.go`）。

契约：

```json
{ "code": 0, "data": {}, "msg": "..." }
```

- 业务成功：`code = 0`（`SUCCESS`）
- 业务失败：`code = 7`（`ERROR`）
- 绝大多数接口 HTTP 状态仍为 **200**；未授权用 `NoAuth` → HTTP 401 + code 7

前端成功判断：`res.code === 0`（见 `web/src/utils/request.js`、Pinia user store）。

---

## Service Layer

- 签名返回 `error`（及业务结果）
- 业务失败：`errors.New("...")` / `fmt.Errorf("中文说明")` / 透传 GORM error
- **不**写 HTTP，**不**调 `response.*`

示例：`server/service/system/sys_user.go`、`server/plugin/ticket/service/order.go`。

---

## API Layer

1. `ShouldBindJSON` / `ShouldBindQuery` 失败 → `response.FailWithMessage(err.Error(), c)` 并 `return`
2. 调用 Service；若 `err != nil` → `FailWithMessage`（固定中文或 `err.Error()`）
3. 成功 → `Ok` / `OkWithMessage` / `OkWithData` / `OkWithDetailed`
4. 分页成功常用：

```go
response.OkWithDetailed(response.PageResult{
    List: list, Total: total, Page: page, PageSize: pageSize,
}, "获取成功", c)
```

`PageResult` 定义于 `server/model/common/response/common.go`。

参考 API：

- 核心：`server/api/v1/system/sys_user.go`、`sys_api.go`
- 插件：`server/plugin/ticket/api/order.go`
- 小程序：`server/plugin/ticket/api/mini/`、`server/plugin/camping/api/mini/`

常用 helper：`Ok`、`OkWithMessage`、`OkWithData`、`OkWithDetailed`、`Fail`、`FailWithMessage`、`FailWithDetailed`、`NoAuth`。

---

## Logging on Errors

核心 system API 常见模式：先 `global.GVA_LOG.Error("...", zap.Error(err))`，再 `FailWithMessage`。

Ticket 等业务插件 API 多数**只 Fail、不打日志**——这是现状。新代码：

- 核心/系统路径：对齐 system，失败记 Error
- 插件路径：至少对非预期错误（DB、panic 级）记日志；纯业务校验可只返回文案

---

## Anti-patterns

- Service 内 `c.JSON` / `response.Fail*`
- API 忽略 Service `error` 仍返回成功
- 改用非 `{code,data,msg}` 的私有响应结构
- 把可预期业务失败做成 HTTP 5xx（本项目惯例是 JSON code=7）
