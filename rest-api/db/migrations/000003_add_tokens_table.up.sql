CREATE TABLE IF NOT EXISTS tokens
(
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    token      VARCHAR(255) NOT NULL,
    user_id    uuid,
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_id
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE cascade
);


