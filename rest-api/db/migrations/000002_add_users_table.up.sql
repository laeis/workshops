CREATE TABLE IF NOT EXISTS users
(
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email      VARCHAR(255) UNIQUE NOT NULL,
    password   VARCHAR(100)        NOT NULL,
    timezone   VARCHAR(10)         NOT NULL,
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE tasks
    ADD CONSTRAINT fk_users_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;