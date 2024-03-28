-- +goose Up
-- +goose StatementBegin
CREATE INDEX messages_created_at_idx ON messages (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX messages_created_at_idx;
-- +goose StatementEnd
