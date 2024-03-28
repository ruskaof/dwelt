-- +goose Up
-- +goose StatementBegin
ALTER TABLE messages ADD COLUMN user_id BIGINT NOT NULL REFERENCES users(id) DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE messages DROP COLUMN user_id;
-- +goose StatementEnd
