ALTER TABLE `short_url_map`
    ADD COLUMN `expire_at` TIMESTAMP NULL DEFAULT NULL COMMENT '过期时间'
AFTER `surl`;