# Implement: 门票订单按核销场合统计

## Checklist

1. [x] Request：`TicketVenueVerifyStatsReq`（`month` string，可选）
2. [x] Service：`GetVenueVerifyStats(month string) (monthOut string, items []..., err)`
   - 解析 `YYYY-MM`；空则当前月
   - 半开区间查 `order_verify_records`
   - 白名单补齐 count=0
3. [x] API：`GetVenueVerifyStats` → `OkWithData`
4. [x] Router：`g.GET("getVenueVerifyStats", ...)`
5. [x] `initialize/api.go` 注册路径描述
6. [x] 前端 `order.js` 封装接口
7. [x] `order.vue` 汇总条 + 月份选择 + 独立拉取

## Validation

```bash
# 后端编译（在 server 目录）
go build ./plugin/ticket/...

# 手工
# 1. 打开订单管理：默认当前月，四场合有数字（含 0）
# 2. 切换到有多场合核销的月份：数字与库中该月 GROUP BY 一致
# 3. 空 venue 记录多的月份：总数不包含它们
# 4. 改列表筛选：汇总数字不变（直至改月份）
```

## Risky files / Rollback

- `service/order.go`：新增查询，勿改动现有 `GetList` / `VerifyOrder`
- `view/order.vue`：仅在搜索区与表格间插入块，勿改列表逻辑
- 回滚：移除 API 路由 + 前端汇总条；无 schema 回退

## Review gate before start

- [x] `prd.md` 已收敛
- [x] `design.md` / `implement.md` 齐全
- [ ] 用户审核规划通过后执行 `task.py start`
