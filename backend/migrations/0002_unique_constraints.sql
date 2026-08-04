-- ============================================================================
-- 0002_unique_constraints.sql (REL-03)
-- 同步幂等的数据库层硬保证（需先执行 0001 清理，否则创建索引会因重复行失败）。
-- MySQL 唯一索引允许多个 NULL，世界杯历史数据（external 为 NULL）不受影响。
-- ============================================================================

-- 比赛：同一外部来源 + 同一外部 ID 唯一（同步幂等、防止并发重复插入）
CREATE UNIQUE INDEX idx_matches_ext_unique
ON matches (external_provider, external_id);

-- 提醒：同一用户对同一场比赛只能有一条提醒
CREATE UNIQUE INDEX idx_reminders_user_match_unique
ON reminders (user_id, match_id);

-- 校验：查看索引是否创建成功
SHOW INDEX FROM matches WHERE Key_name = 'idx_matches_ext_unique';
SHOW INDEX FROM reminders WHERE Key_name = 'idx_reminders_user_match_unique';
