-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS chats
(
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULl,
    name       TEXT      NOT NULL
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users_chats
(
    user_id BIGINT REFERENCES users (id) NOT NULL,
    chat_id BIGINT REFERENCES chats (id) NOT NULL,
    PRIMARY KEY (user_id, chat_id)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS messages
(
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP                    NOT NULL,
    text       TEXT,
    photo      BYTEA,
    chat_id    BIGINT REFERENCES chats (id) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE messages;
DROP TABLE users_chats;
DROP TABLE chats;
-- +goose StatementEnd
