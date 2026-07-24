# 限时活动 — 实施计划

## Checklist

1. **后端插件骨架**  
   - 新建 `server/plugin/limitedActivity/`（plugin.go、initialize/{api,gorm,menu,router}、model/service/api/router）。  
   - 在 `server/plugin/register.go` blank import 注册。

2. **数据模型与迁移**  
   - `Activity`、`ActivityOrder`、`ActivityOrderVerifyRecord` AutoMigrate。  
   - 订单状态与字段按 `design.md`。

3. **活动 Admin + Mini API**  
   - 后台 CRUD、显示状态、图片/二维码上传字段对接现有上传能力。  
   - 小程序：列表（仅显示+时间相关筛选）、详情。

4. **报名订单 + 名额**  
   - 下单事务占用 `sold`；超时关单释放；mini 订单列表/详情。

5. **支付接入**  
   - 扩展 `server/api/v1/mini/pay.go`：Create / Notify（及退款回调成功路径）支持 `limitedActivity` / 前缀 `A`。  
   - 回调幂等、金额校验对齐 ticket。

6. **核销**  
   - Service 多次核销 + 记录；公开查询/核销 API。  
   - 扩展 `web/src/plugin/camping/view/h5Verify.vue` + ticket/camping 同级的 api 封装。

7. **后台订单与按比例退款**  
   - 订单管理页（参考 `web/src/plugin/ticket/view/order.vue`）。  
   - 退款按钮：按未核销比例调微信退款；回退名额；禁止再核销。

8. **菜单与前端插件**  
   - `web/src/plugin/limitedActivity/`：活动管理、订单管理页面与 api。  
   - initialize Menu 注册。

9. **验证**  
   - 见下方 Validation；修问题后再请 `task.py start` 后的 check。

## Validation Commands

```bash
# 后端编译
cd server && go build -o nul .

# 有测试则跑插件相关（若暂无单测可跳过）
cd server && go test ./plugin/limitedActivity/...

# 前端（按仓库习惯）
cd web && npm run build
# 或至少保证相关 vue/js 无语法错误
```

手工验收对照 `prd.md` AC1–AC9（支付/退款需调试环境微信配置）。

## Risky Files / Rollback Points

| 风险点 | 说明 |
|--------|------|
| `server/api/v1/mini/pay.go` | 公共支付入口，改动需隔离分支、勿破坏 ticket |
| `web/src/plugin/camping/view/h5Verify.vue` | 共用核销页，扩展 type 时保持原 type 行为 |
| 名额 `sold` 并发 | 更新须事务 + 条件 `sold + qty <= quota` |

回滚：撤销 register import、隐藏菜单、回退 pay.go 与 h5Verify 改动。

## Follow-ups Before `task.py start`

- [x] `prd.md` 收敛完成  
- [x] `design.md` / `implement.md` 已写  
- [x] `implement.jsonl` / `check.jsonl` 已填真实 spec 条目  
- [ ] 用户审阅规划并同意开始实现  

## Notes

- 小程序页面不在本仓；实现以 API + Swagger 注释为准。  
- 不拆父子任务：同一插件交付链路紧耦合，单任务推进。
