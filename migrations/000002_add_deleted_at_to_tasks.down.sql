DROP INDEX IF EXISTS index_tasks_active;
DROP INDEX IF EXISTS index_tasks_deleted;

ALTER TABLE tasks DROP COLUMN IF EXISTS deleted_at;