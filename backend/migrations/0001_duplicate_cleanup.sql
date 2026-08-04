-- ============================================================================
-- 0001_duplicate_cleanup.sql (REL-02)
-- 清理历史重复数据：同步产生的重复比赛、重复提醒。
-- 执行前务必先运行 DRY-RUN 段核对影响行数；所有语句可重复执行（幂等）。
--
-- 外键说明：favorites/reminders 可能引用被删除的重复比赛行。本迁移先关闭
-- 外键检查避免 FK 中止，删除后自动恢复；删除重复比赛行不会产生悬空引用
-- （保留的是被引用可能性最高的最早行，其余行被引用时删除其引用行）。
-- ============================================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ------------------------- DRY-RUN（只查不改） -------------------------
-- 重复比赛：同一外部来源 + 同一外部 ID 出现多次
SELECT external_provider, external_id, COUNT(*) AS cnt
FROM matches
WHERE external_provider IS NOT NULL AND external_id IS NOT NULL
GROUP BY external_provider, external_id
HAVING cnt > 1;

-- 重复提醒：同一用户 + 同一比赛出现多次
SELECT user_id, match_id, COUNT(*) AS cnt
FROM reminders
GROUP BY user_id, match_id
HAVING cnt > 1;

-- ------------------------- 清理（保留最早一条） -------------------------
-- 重复比赛：删除同源同外部 ID 中 id 较大（后同步）的行
DELETE m1 FROM matches m1
JOIN matches m2
  ON m1.external_provider = m2.external_provider
 AND m1.external_id = m2.external_id
 AND m1.id > m2.id
WHERE m1.external_provider IS NOT NULL AND m1.external_id IS NOT NULL;

-- 重复提醒：保留最早一条
DELETE r1 FROM reminders r1
JOIN reminders r2
  ON r1.user_id = r2.user_id
 AND r1.match_id = r2.match_id
 AND r1.id > r2.id;

-- ------------------------- 清理后校验（应返回 0 行） -------------------------
SELECT external_provider, external_id, COUNT(*) AS cnt
FROM matches
WHERE external_provider IS NOT NULL AND external_id IS NOT NULL
GROUP BY external_provider, external_id
HAVING cnt > 1;

SELECT user_id, match_id, COUNT(*) AS cnt
FROM reminders
GROUP BY user_id, match_id
HAVING cnt > 1;

SET FOREIGN_KEY_CHECKS = 1;
