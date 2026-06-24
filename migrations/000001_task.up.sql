CREATE TABLE IF NOT EXISTS tasks
(
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    status      SMALLINT     NOT NULL DEFAULT 1,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP    NOT NULL DEFAULT NOW()

    CONSTRAINT task_status_valid CHECK (status IN (1, 2, 3))
    );

CREATE INDEX index_task_status ON tasks (status);
CREATE INDEX ON tasks(created_at);
CREATE INDEX ON tasks(updated_at);