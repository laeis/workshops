CREATE TYPE category AS ENUM ('event', 'note');

CREATE TABLE IF NOT EXISTS tasks
(
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id     uuid NOT NULL,
    category    category DEFAULT 'event' NOT NULL,
    title       VARCHAR(50)  NOT NULL,
    description VARCHAR(355) NOT NULL,
    start_date  TIMESTAMP    NOT NULL,
    created_at  TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    DEFAULT CURRENT_TIMESTAMP
);