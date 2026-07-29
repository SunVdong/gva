# 限时活动增加活动地点字段

## Goal

为限时活动增加选填的「活动地点」文本字段，后台可维护，活动列表/详情接口（含小程序）可返回，供调用方展示。

## Background

- 插件：`server/plugin/limitedActivity`；管理端：`web/src/plugin/limitedActivity/view/activity.vue`。
- `Activity`（表 `limited_activities`）现无地点字段；插件 `AutoMigrate` 负责加列。
- 小程序活动列表/详情直接返回 `Activity`；本仓库无限时活动小程序页面。
- 订单仅快照活动名称；本轮不扩展订单/核销。

## Requirements

- R1：`Activity` 新增地点字段：JSON / Go `address`，库列 `address`，字符串，建议 `size:256`。
- R2：管理端新增/编辑表单可填写地点（**选填**）；列表展示（空显示 `-`）；不做地点搜索。
- R3：创建/更新/详情/列表（管理端 + 小程序）响应含 `address`；历史数据为空字符串。
- R4：依赖现有 AutoMigrate，不写手工 SQL。

## Acceptance Criteria

- [x] AC1：管理端可保存并回显地点；留空可保存。（R1, R2）
- [x] AC2：管理端列表展示地点，空为 `-`。（R2）
- [x] AC3：管理端与小程序活动接口 JSON 含 `address`。（R3）
- [x] AC4：已有活动不填地点仍可正常读写。（R2, R4）

## Out of Scope

- 地图选点、经纬度、POI。
- 订单地点快照；订单详情 / 核销 H5 展示地点。
- 本仓库外的小程序 UI。

## Decisions

| 决策 | 结论 |
|------|------|
| 是否必填 | 选填 |
| 订单/核销 | 本轮不做 |
| 字段名 | `address`（非 `location`） |
| 字段形态 | 纯文本 |
