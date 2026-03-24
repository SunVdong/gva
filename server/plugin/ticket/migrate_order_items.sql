-- 数据迁移脚本：将 order_items 数据合并到 orders 表
-- 前置条件：GORM AutoMigrate 已为 orders 表自动添加 sku_id/sku_name/price/quantity/visit_date 列
-- 执行完毕并确认数据正确后，可手动 DROP TABLE order_items;

-- 1. 回填数据（一单一SKU，取第一条 order_item）
UPDATE orders o
JOIN order_items oi ON oi.order_id = o.id
SET o.sku_id     = oi.sku_id,
    o.sku_name   = oi.sku_name,
    o.price      = oi.price,
    o.quantity   = oi.quantity,
    o.visit_date = oi.visit_date
WHERE o.sku_id IS NULL OR o.sku_id = 0;

-- 2. 验证：检查是否所有订单都已回填
SELECT COUNT(*) AS unfilled_orders
FROM orders
WHERE sku_id IS NULL OR sku_id = 0;

-- 3. 确认无误后删除 order_items 表
-- DROP TABLE IF EXISTS order_items;
