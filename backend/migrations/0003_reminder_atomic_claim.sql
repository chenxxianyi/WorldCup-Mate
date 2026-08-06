-- ============================================================================
-- 0003_reminder_atomic_claim.sql (REL-08 / ADM-13 / OBS-04)
-- 为提醒表添加原子认领、超时回收、退避重试和运营状态字段。
-- 适用于从旧版本升级到新版本的场景，幂等且可重复执行。
-- ============================================================================

ALTER TABLE reminders
    ADD COLUMN IF NOT EXISTS claim_token  VARCHAR(36) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS claimed_at   TIMESTAMP NULL,
    ADD COLUMN IF NOT EXISTS next_retry_at TIMESTAMP NULL,
    ADD COLUMN IF NOT EXISTS last_error   VARCHAR(500) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS notification_id INT NULL,
    ADD COLUMN IF NOT EXISTS worker_id    VARCHAR(36) NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_reminders_claim_token
    ON reminders (claim_token, status);
CREATE INDEX IF NOT EXISTS idx_reminders_next_retry_at
    ON reminders (next_retry_at, status, remind_at);
CREATE INDEX IF NOT EXISTS idx_reminders_last_error
    ON reminders (last_error);

-- 校验
SHOW COLUMNS FROM reminders WHERE Field IN (
    'claim_token', 'claimed_at', 'next_retry_at',
    'last_error', 'notification_id', 'worker_id'
);
