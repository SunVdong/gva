# 门票订单按核销场合统计

## Goal

在景区门票订单管理页增加按核销场合的核销次数汇总，支持按核销月份查看，默认当前月，方便运营对比各点位核销量。

## Background

- 多场合核销已落地（`.trellis/spec/backend/ticket-multi-venue-verify.md`）。
- 核销记录表 `order_verify_records.venue` 存场合 code；非多场合订单核销时 `venue` 为空。
- 固定场合：`zhongshanling` / `zhaozhao` / `lululand` / `hongshan`（`server/plugin/ticket/model/venue.go`）。
- 订单管理页已有场合筛选，无场合汇总统计；ticket 插件无现成聚合 API。

## Requirements

- R1. 按核销场合汇总「核销次数」（对 `order_verify_records` 按 `venue` 计数）；不做订单数、金额统计。
- R2. 独立「核销月份」筛选：按 `verified_at` 落在该自然月内；默认当前月；不与订单列表其它筛选项联动。
- R3. UI 放在订单管理页搜索区下方、表格上方：月份选择器 + 各场合核销次数。
- R4. 只统计四个固定白名单场合；忽略 `venue` 为空的记录；无数据的场合显示 0。

## Acceptance Criteria

- [ ] AC1. 订单管理页可看到四个场合各自的核销次数（中山陵 / 爪爪 / lululand / 红山）。
- [ ] AC2. 进入页面时统计默认按当前自然月；可切换月份并刷新数字。
- [ ] AC3. 切换月份后数字仅反映该月 `verified_at` 落在月内的核销记录。
- [ ] AC4. 空 `venue` 记录不计入任何场合，也不出现「未指定场合」项。
- [ ] AC5. 某场合当月无核销时显示 0，而非缺项。
- [ ] AC6. 订单列表原有筛选 / 分页行为不受影响（统计与列表解耦）。

## Out of Scope

- H5 核销流程、SKU 多场合开关、camping / limitedActivity。
- 订单数 / 金额 / 导出 / 图表大屏。
- 任意起止日期（仅按自然月）。
- 单独统计页或弹窗。
- Schema 迁移（复用现有 `order_verify_records`）。

## Decisions

| 决策 | 结论 |
|------|------|
| 指标 | 仅核销次数 |
| 时间 | 独立核销月份，默认当前月 |
| UI | 列表上方汇总条 |
| 空 venue | 忽略，不单独计入 |
