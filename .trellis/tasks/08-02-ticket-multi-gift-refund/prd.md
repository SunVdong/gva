# 景点门票多次票赠送次数与整单比例退款

## Goal

为景点门票增加多次票赠送次数，统一「次数 × 数量」核销池，并按付费次数剩余比例做整单微信退款（用户端与后台一致）。核销先消耗付费次数 A 再消耗赠送次数 B；不拆码；多次票每单限购 1 张。

## Background

- 插件：`server/plugin/ticket`。SKU 有 `useTimes`，无赠送字段；订单 `totalUseTimes` 现状未乘 `quantity`；核销码为 `order_no`。
- 用户退款现状：仅 `verifiedTimes==0` 全额退；后台多次票退款只改状态、不调微信。
- 参考实现：`limitedActivity` 的 `CalcRefundFen` / `AdminRefund` / 按剩余回退占用。

## Requirements

### R1. SKU 赠送次数
- 多次票可配置单张赠送次数 `giftUseTimes`（p≥0）；单次票固定 0。

### R2. 多次票每单限购 1 张
- `ticketType==2` 时 `quantity` 必须为 1；与 `limitBuy` 并存。

### R3. 订单次数快照（乘数量）
- 单张：`m=useTimes`（≤0→1），`p=giftUseTimes`（单次→0）。
- 订单：`A=quantity×m`，`B=quantity×p`，`totalUseTimes=A+B`。
- 核销累加 `verifiedTimes`；先 A 后 B（`paidConsumed=min(verified,A)`）；用尽后 status=已核销。

### R4. 整单比例退款（单次+多次；用户端=后台）
- `remainingPaid=max(A-min(verified,A),0)`
- `refundFen=Round(payAmount×100×remainingPaid/A)`；`remainingPaid==0` 不可退。
- 用户端与后台均真实调用微信退款；落库 `refundAmount`。
- 替换后台「仅改状态」逻辑。

### R5. 日历 sold 回退
- 单次票（m=1）：`sold -= remainingPaid`
- 多次票（每单 1 张）：退款成功且 `remainingPaid>0` 时 `sold -= 1`

### R6. 不拆码
- 不新增凭证表；核销码继续 `order_no`；单次多张共享一码与次数池。

### R7. 兼容
- 新单写入 `paidUseTimes`/`giftUseTimes`；老单缺字段时按 `A=totalUseTimes,B=0` 兼容计算。

## Acceptance Criteria

- [ ] AC1: 多次票 SKU 可配置赠送次数；单次 gift=0。
- [ ] AC2: 多次票 quantity≠1 被拒。
- [ ] AC3: 新单 `totalUseTimes=quantity×(m+p)`，`paidUseTimes=quantity×m`。
- [ ] AC4: 核销先 A 后 B，用尽后已核销。
- [ ] AC5: 单次/多次在 remainingPaid>0 时可整单比例退，金额符合四舍五入公式。
- [ ] AC6: 用户端与后台退款规则及微信调用一致。
- [ ] AC7: sold 回退符合 R5。
- [ ] AC8: 不拆码；API `canRefund`/可退金额与规则一致。

## Out of Scope

- 按张拆码、用户自选按张退、多次票一单多张、转赠/改期/过期自动退、小程序仓外完整 UI 重构。
