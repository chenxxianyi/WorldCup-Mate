# Database Migrations

Migrations are executed in numerical order. Each migration file is idempotent
(uses `IF NOT EXISTS` / `ALTER TABLE ... ADD COLUMN IF NOT EXISTS`) so it can
be run repeatedly without side effects.

## Running migrations

```bash
# Run all migrations against the configured MySQL
mysql -u root -p worldcup_mate < migrations/0001_duplicate_cleanup.sql
mysql -u root -p worldcup_mate < migrations/0002_unique_constraints.sql
mysql -u root -p worldcup_mate < migrations/0003_reminder_atomic_claim.sql
```

Or use the `migrate` tool:
```bash
migrate -path migrations -database "mysql://root:password@tcp(localhost:3306)/worldcup_mate" up
```

## Principles (DB-01)

- Migrations are versioned and always idempotent.
- No migration deletes or renames existing columns without a rollback path.
- Migration failures leave the database in a partially-migrated state that can
  be detected and resumed by re-running.
- Application startup does NOT run `AutoMigrate` in production; explicit
  migrations are the source of truth.
