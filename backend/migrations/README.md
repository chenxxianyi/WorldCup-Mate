# 数据库迁移（REL-01）

## 执行规则

1. 迁移文件按**编号顺序**执行，每个文件只执行一次。
2. 已在生产库执行的编号**不要修改文件内容**——需要变更时新增下一个编号的文件。
3. 迁移前先备份（`mysqldump`）或至少执行文件内的 DRY-RUN 段核对影响行数。
4. 执行方式（宿主机 MySQL，非 Docker）：

```bash
mysql -u root -p worldcup < backend/migrations/0001_duplicate_cleanup.sql
mysql -u root -p worldcup < backend/migrations/0002_unique_constraints.sql
```

Docker 环境：

```bash
docker compose exec -T mysql mysql -u root -p worldcup < backend/migrations/0001_duplicate_cleanup.sql
```

5. 迁移失败时：先修复数据（如重复行），再重跑该文件。0001 语句为幂等设计可重跑；0002 的 `CREATE UNIQUE INDEX` 重复执行会报 `Duplicate key name`——属"已应用"的正常信号，无需处理。

## 文件清单

| 编号 | 内容 | 状态 |
|------|------|------|
| 0001_duplicate_cleanup.sql | REL-02：清理历史重复的比赛与提醒数据（含 DRY-RUN 段） | 待执行 |
| 0002_unique_constraints.sql | REL-03：比赛与提醒唯一约束（依赖 0001 先执行） | 待执行 |

## 说明

- `matches` 的唯一约束基于 `(external_provider, external_id)`；世界杯历史比赛这两列为 NULL，**不受唯一索引影响**（MySQL 唯一索引允许多个 NULL）。
- `reminders` 唯一约束基于 `(user_id, match_id)`，防止同一用户对同一场比赛重复设置提醒。
- 表结构层面的后续变更（如新增列）由 GORM AutoMigrate 在服务启动时自动处理（只增不删），与编号迁移互不冲突。
