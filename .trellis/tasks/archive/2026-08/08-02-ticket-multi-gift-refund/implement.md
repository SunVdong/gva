# Implement: 门票赠送次数与整单比例退款

## Checklist（顺序）

1. **Model**
   - [ ] `ticket_sku.go`：`GiftUseTimes`
   - [ ] `order.go`：`PaidUseTimes`、`GiftUseTimes`、`RefundAmount`
   - [ ] request / 校验：SKU 保存时多次票 gift≥0；单次强制 0

2. **下单**
   - [ ] `service/order.go` `CreateOrder`：多次票 quantity==1；写入 A/B/total=A+B
   - [ ] mini 下单入口错误文案可读

3. **核销**
   - [ ] `VerifyOrder`：拒绝 status≠1（含 7）；上限用 totalUseTimes（已含赠送）
   - [ ] API/H5/mini 响应附带 `paidUseTimes`/`giftUseTimes`/`remainingTimes`（可选展示）

4. **退款核心**
   - [ ] `CalcRefundFen`（按 A/remainingPaid，`math.Round`）
   - [ ] `RequestRefund` / `ApplyRefundSuccess`：锁单→微信→refundAmount→sold 回退（R5）
   - [ ] 替换 `RefundPendingVerifyMultiTicket` 为真实退款
   - [ ] 改造 `server/api/v1/mini/pay.go` `Refund` 与回调落库（传入 refundFen / refundAmount）

5. **管理端 / 前端**
   - [ ] `product.vue`：多次票赠送次数表单项
   - [ ] `order.vue`：展示 A/B/已核销；退款确认展示预计退款额（对齐 limitedActivity）
   - [ ] `h5Verify.vue`：如需展示总次数含赠送，仅文案/字段，不改码体系

6. **契约**
   - [ ] mini `canRefund`：`status==1 && remainingPaid>0`（及有 wx_transaction_id）
   - [ ] 可选返回 `refundAmountFen` 预估

## Validation

```bash
# 编译
cd server && go build ./plugin/ticket/...
cd server && go build ./api/v1/mini/...

# 手工场景（支付沙箱/测试环境）
# 1) 多次票 gift=1,m=2 下单 quantity=2 → 应失败；quantity=1 → total=3
# 2) 核销 1 次后退款 → refund ≈ pay*2/2 = 全额? A=2, remainingPaid=1 → 一半
# 3) 核销 2 次（A 用尽）后再核赠送 → 不可退
# 4) 单次票 quantity=3，核 1 次 → refund 2/3；sold -= 2
# 5) 后台退款与用户端金额一致
```

## Risky files / Rollback points

| 文件 | 风险 |
|------|------|
| `server/api/v1/mini/pay.go` | 公共支付；改退款金额与回调时勿破坏 limitedActivity 分支 |
| `service/order.go` 退款事务 | 双花核销、sold 少减/多减 |
| 老单兼容分支 | 缺 A 字段时勿除零 |

回滚：保留 status=7 释放；sold 回退失败打 Warn 不阻断主状态（对齐 limitedActivity 可评估，优先事务内失败回滚）。

## Review gate before start

- [x] prd / design / implement 齐备
- [x] 用户确认可 `task.py start`
- [x] implement.jsonl / check.jsonl 已填真实 spec

## Spec update (Phase 3.3)

- [x] 新增 `.trellis/spec/backend/ticket-gift-proportional-refund.md`
- [x] 更新 `mini-pay-integration.md`（门票比例退 / 与活动分母差异）
- [x] 更新 `backend/index.md` 索引
- [x] `guides/cross-layer-thinking-guide.md` 增加退款检查清单
- [ ] 提交（用户要求先不 commit）
