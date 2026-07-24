# Implement: 门票多次票多场合核销

## Checklist

1. **Model**
   - [ ] `ticket_sku.go`：`SupportMultiVenue`
   - [ ] `order.go`：`SupportMultiVenue`（快照）
   - [ ] `order_verify_record.go`：`Venue`
   - [ ] 场合常量 + 校验（可放 `model/venue.go` 或 service 内）

2. **Service**
   - [ ] SKU 创建/更新：单次票强制 `SupportMultiVenue=false`
   - [ ] `CreateOrder`：写入快照
   - [ ] `VerifyOrder` / `VerifyOrderByOrderNoPublic`：按快照校验并写入 `Venue`
   - [ ] `GetList`：`ticketType` join 筛选；`venue` EXISTS 筛选

3. **API / Request**
   - [ ] `TicketOrderSearch` 增加 `ticketType`、`venue`
   - [ ] `VerifyOrderByCodePublic` 绑定 body `venue` 并下传 service
   - [ ] 确认 AutoMigrate 已覆盖上述 model（现有 gorm init）

4. **Admin 前端**
   - [ ] `product.vue`：多次票 checkbox「是否支持多场合」
   - [ ] `order.vue`：票种/场合筛选项；详情核销记录场合列
   - [ ] 场合 options 常量与后端 code 对齐

5. **H5**
   - [ ] `h5Verify.vue` ticket 分支：`supportMultiVenue` 时展示场合必选
   - [ ] `verifyTicketOrderByCodePublic` 传 `venue`
   - [ ] 核销记录列表展示场合 label（若有）

6. **自测**
   - [ ] 单次票：无开关、核销无场合
   - [ ] 多次票未勾选：核销无场合
   - [ ] 多次票勾选：无场合失败；有场合成功；二次可选不同场合
   - [ ] 改 SKU 开关不影响旧单
   - [ ] 后台票种/场合筛选；详情可见场合

## Validation Commands

```bash
# 后端编译
cd server && go build -o nul .
# 前端（按仓库惯例，可选）
cd web && npm run build
```

手动：管理端 SKU → 下单（或造单）→ H5 `#/h5/verify?type=ticket&code=` → 后台订单筛选/详情。

## Risky Files / Rollback

| 文件 | 风险 |
|------|------|
| `service/order.go` VerifyOrder | 核销主路径，必测多场合开/关 |
| `h5Verify.vue` | 共享页，勿影响 camping/limitedActivity 分支 |
| `api/order.go` VerifyOrderByCodePublic | 签名变更需兼容无 body 旧调用 |

回滚：还原上述改动；DB 新列可留。

## Status

- [x] 实现完成（trellis-implement）
- [x] 质量检查 PASS（trellis-check）
- [x] Spec 已更新：`.trellis/spec/backend/ticket-multi-venue-verify.md`
- [ ] Phase 3.4 用户确认后提交代码
