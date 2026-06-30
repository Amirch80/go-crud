ALTER TABLE tasks ADD COLUMN deleted_at TIMESTAMP DEFAULT NULL;

CREATE INDEX index_tasks_active ON tasks (created_at) WHERE deleted_at IS NULL;
CREATE INDEX index_tasks_deleted ON tasks (deleted_at) WHERE deleted_at IS NOT NULL;