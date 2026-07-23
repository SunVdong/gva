# Logging Guidelines

> 结构化日志以 Zap 为主，入口优先 `global.GVA_LOG`。

---

## Overview

- 类型：`server/global/global.go` → `GVA_LOG *zap.Logger`
- 字段：`zap.Error(err)`、`zap.String(...)`、`zap.Any(...)` 等
- GORM 日志桥接：`server/initialize/internal/gorm_logger_writer.go`

---

## Preferred Usage

```go
global.GVA_LOG.Info("message", zap.String("key", val))
global.GVA_LOG.Warn("message", zap.Error(err))
global.GVA_LOG.Error("message", zap.Error(err))
global.GVA_LOG.Debug("message", ...)
```

真实引用：

- `server/initialize/reload.go`
- `server/api/v1/system/sys_api.go`、`sys_auto_code.go`
- `server/service/system/sys_user.go`、`jwt_black_list.go`
- `server/utils/upload/aliyun_oss.go`
- `server/plugin/ticket/plugin.go`（Info）
- `server/plugin/ticket/service/order.go`（Warn）

---

## What to Log

| 场景 | 级别 | 说明 |
|------|------|------|
| 启动/插件注册/配置加载 | Info | 可运维可见 |
| 可恢复异常、降级、重试 | Warn | 如订单服务告警 |
| Service/API 非预期失败、外部依赖失败 | Error | 带 `zap.Error` |
| 详细诊断 | Debug | 默认生产勿刷屏 |

避免：成功路径逐条 Debug 刷屏；日志中打印密码、token、完整身份证等敏感信息。

---

## Inconsistency（现状）

部分代码使用 `zap.L()` 而非 `global.GVA_LOG`（例：`server/plugin/ticket/initialize/gorm.go`）。

**新代码一律用 `global.GVA_LOG`**，除非在 `GVA_LOG` 尚未初始化的极早启动阶段。

---

## Anti-patterns

- `fmt.Println` / `log.Println` 作为业务日志
- Error 级别不带 `err` 字段导致无法排查
- 与 `response.FailWithMessage` 重复拼一大段用户文案进日志却无结构化字段
